package service

import (
	"errors"
	"strconv"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/config"
	"github.com/sheet-platform/backend/internal/dto/request"
	"github.com/sheet-platform/backend/internal/model"
	"github.com/sheet-platform/backend/internal/pkg"
	"github.com/sheet-platform/backend/internal/repository"
)

type AuthService struct {
	userRepo *repository.UserRepo
	cfg      *config.Config
	logger   *zap.Logger
	db       *gorm.DB
}

type AuthServiceParams struct {
	fx.In
	UserRepo *repository.UserRepo
	Cfg      *config.Config
	Logger   *zap.Logger
	DB       *gorm.DB
}

func NewAuthService(p AuthServiceParams) *AuthService {
	return &AuthService{userRepo: p.UserRepo, cfg: p.Cfg, logger: p.Logger, db: p.DB}
}

type LoginResponse struct {
	Token    string `json:"token"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	IsAdmin  bool   `json:"is_admin"`
}

func (s *AuthService) Login(req *request.LoginReq) (*LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("邮箱或密码错误")
		}
		s.logger.Error("find user failed", zap.Error(err))
		return nil, errors.New("登录失败，请稍后重试")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("邮箱或密码错误")
	}

	if user.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}

	ttl, err := time.ParseDuration(s.cfg.JWT.TTL)
	if err != nil {
		ttl = 24 * time.Hour
	}

	token, err := pkg.GenerateJWT(s.cfg.JWT.Secret, ttl, user.ID, user.Username, 0)
	if err != nil {
		s.logger.Error("generate jwt failed", zap.Error(err))
		return nil, errors.New("登录失败，请稍后重试")
	}

	return &LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		IsAdmin:  user.IsAdmin,
	}, nil
}

func (s *AuthService) SendVerificationCode(req *request.SendCodeReq) error {
	exists, err := s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		s.logger.Error("check email exists failed", zap.Error(err))
		return errors.New("发送验证码失败，请稍后重试")
	}
	if exists {
		return errors.New("该邮箱已被注册")
	}

	code := pkg.GenerateAndStoreCode(req.Email)

	if err := pkg.SendVerificationEmail(&s.cfg.Email, req.Email, code); err != nil {
		s.logger.Error("send email failed", zap.Error(err))
		return errors.New("邮件发送失败，请稍后重试")
	}

	s.logger.Info("verification code sent",
		zap.String("email", req.Email),
		zap.String("code", code),
	)
	return nil
}

type ProfileResponse struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	IsAdmin   bool   `json:"is_admin"`
	Points    int    `json:"points"`
	CreatedAt string `json:"created_at"`
}

func (s *AuthService) GetProfile(userID uint) (*ProfileResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	email := ""
	if user.Email != nil {
		email = *user.Email
	}
	return &ProfileResponse{
		UserID:    user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Email:     email,
		Avatar:    user.Avatar,
		IsAdmin:   user.IsAdmin,
		Points:    user.Points,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *AuthService) Register(req *request.RegisterReq) (*LoginResponse, error) {
	if !pkg.VerifyCode(req.Email, req.Code) {
		return nil, errors.New("验证码错误或已过期")
	}

	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		s.logger.Error("check user exists failed", zap.Error(err))
		return nil, errors.New("注册失败，请稍后重试")
	}
	if exists {
		return nil, errors.New("用户名已存在")
	}

	emailExists, err := s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		s.logger.Error("check email exists failed", zap.Error(err))
		return nil, errors.New("注册失败，请稍后重试")
	}
	if emailExists {
		return nil, errors.New("该邮箱已被注册")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("hash password failed", zap.Error(err))
		return nil, errors.New("注册失败，请稍后重试")
	}

	nickname := req.Nickname
	if nickname == "" {
		nickname = req.Username
	}

	registerPoints := s.getRegisterPoints()

	email := req.Email
	user := &model.User{
		Username: req.Username,
		Password: string(hashed),
		Email:    &email,
		Nickname: nickname,
		Points:   registerPoints,
		Status:   1,
	}

	if err := s.userRepo.Create(user); err != nil {
		s.logger.Error("create user failed", zap.Error(err))
		return nil, errors.New("注册失败，请稍后重试")
	}

	ttl, err := time.ParseDuration(s.cfg.JWT.TTL)
	if err != nil {
		ttl = 24 * time.Hour
	}
	token, err := pkg.GenerateJWT(s.cfg.JWT.Secret, ttl, user.ID, user.Username, 0)
	if err != nil {
		s.logger.Error("generate jwt failed", zap.Error(err))
		return nil, errors.New("注册成功，但登录失败，请尝试登录")
	}

	return &LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		IsAdmin:  user.IsAdmin,
	}, nil
}

func (s *AuthService) getRegisterPoints() int {
	var cfg model.SystemConfig
	if err := s.db.Where("config_key = ?", "register_points").First(&cfg).Error; err != nil {
		return 100
	}
	pts, err := strconv.Atoi(cfg.ConfigVal)
	if err != nil {
		return 100
	}
	return pts
}
