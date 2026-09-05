package service

import (
	"errors"
	"sort"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/model"
	"github.com/sheet-platform/backend/internal/repository"
)

type SheetService struct {
	sheetRepo *repository.SheetMusicRepo
	tagRepo   *repository.TagRepo
	fileRepo  *repository.SheetFileRepo
	logger    *zap.Logger
}

type SheetServiceParams struct {
	fx.In
	SheetRepo *repository.SheetMusicRepo
	TagRepo   *repository.TagRepo
	FileRepo  *repository.SheetFileRepo
	Logger    *zap.Logger
}

func NewSheetService(p SheetServiceParams) *SheetService {
	return &SheetService{
		sheetRepo: p.SheetRepo,
		tagRepo:   p.TagRepo,
		fileRepo:  p.FileRepo,
		logger:    p.Logger,
	}
}

type SheetDetail struct {
	model.SheetMusic
	Tags        []model.Tag       `json:"tags"`
	Files       []model.SheetFile `json:"files,omitempty"`
	Audio       []model.AudioFile `json:"audio,omitempty"`
	LikeCount   int64             `json:"like_count"`
	IsLiked     bool              `json:"is_liked"`
	IsPurchased bool              `json:"is_purchased"`
	HasFreeFile bool              `json:"has_free_file"`
	HasPaidFile bool              `json:"has_paid_file"`
}

type SheetListReq struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Keyword  string `json:"keyword"`
	TagIDs   []uint `json:"tag_ids,omitempty"`
}

type SheetListItem struct {
	model.SheetMusic
	Tags []model.Tag `json:"tags"`
}

type SheetListResp struct {
	List     []SheetListItem `json:"list"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"size"`
}

func (s *SheetService) List(req *SheetListReq) (*SheetListResp, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	sheets, total, err := s.sheetRepo.List(req.Page, req.PageSize, req.Keyword, req.TagIDs)
	if err != nil {
		s.logger.Error("list sheets failed", zap.Error(err))
		return nil, errors.New("查询乐谱列表失败")
	}

	sheetIDs := make([]uint, len(sheets))
	for i, sheet := range sheets {
		sheetIDs[i] = sheet.ID
	}
	tagMap, err := s.tagRepo.MapBySheetIDs(sheetIDs)
	if err != nil {
		s.logger.Error("load list tags failed", zap.Error(err))
		tagMap = map[uint][]model.Tag{}
	}

	list := make([]SheetListItem, len(sheets))
	for i, sheet := range sheets {
		tags := tagMap[sheet.ID]
		if tags == nil {
			tags = []model.Tag{}
		}
		list[i] = SheetListItem{SheetMusic: sheet, Tags: tags}
	}

	return &SheetListResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (s *SheetService) Detail(id uint) (*SheetDetail, error) {
	sheet, err := s.sheetRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("乐谱不存在")
		}
		s.logger.Error("find sheet failed", zap.Error(err))
		return nil, errors.New("查询乐谱失败")
	}

	// Increment view count (best-effort)
	_ = s.sheetRepo.IncrementView(id)

	tags, err := s.tagRepo.FindBySheetID(id)
	if err != nil {
		s.logger.Error("load sheet tags failed", zap.Error(err))
		tags = []model.Tag{}
	}
	files, err := s.fileRepo.ListBySheetID(id)
	if err != nil {
		s.logger.Error("load sheet files failed", zap.Error(err))
	}
	audio, err := s.fileRepo.ListAudioBySheetID(id)
	if err != nil {
		s.logger.Error("load sheet audio failed", zap.Error(err))
	}

	return &SheetDetail{
		SheetMusic: *sheet,
		Tags:       tags,
		Files:      files,
		Audio:      audio,
	}, nil
}

func (s *SheetService) FindByID(id uint) (*model.SheetMusic, error) {
	return s.sheetRepo.FindByID(id)
}

func (s *SheetService) ListTags() ([]model.Tag, error) {
	return s.tagRepo.ListAll()
}

// normalizePageSync 校验并按时间升序排列曲谱同步点
func normalizePageSync(points model.PageSyncPoints) (model.PageSyncPoints, error) {
	if len(points) == 0 {
		return nil, nil
	}
	for _, p := range points {
		if p.Time < 0 {
			return nil, errors.New("同步点时间不能为负数")
		}
		if p.Page < 1 {
			return nil, errors.New("同步点页码必须从 1 开始")
		}
	}
	sort.SliceStable(points, func(i, j int) bool { return points[i].Time < points[j].Time })
	for i := 1; i < len(points); i++ {
		if points[i].Time == points[i-1].Time {
			return nil, errors.New("同步点时间不能重复")
		}
	}
	return points, nil
}

func (s *SheetService) Create(sheet *model.SheetMusic, tagNames []string) error {
	if sheet.Title == "" {
		return errors.New("标题不能为空")
	}
	pageSync, err := normalizePageSync(sheet.PageSync)
	if err != nil {
		return err
	}
	sheet.PageSync = pageSync

	if err := s.sheetRepo.Create(sheet); err != nil {
		return err
	}
	return s.syncTags(sheet.ID, tagNames)
}

func (s *SheetService) Update(sheet *model.SheetMusic, tagNames []string) error {
	pageSync, err := normalizePageSync(sheet.PageSync)
	if err != nil {
		return err
	}
	sheet.PageSync = pageSync

	existing, err := s.sheetRepo.FindByID(sheet.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("乐谱不存在")
		}
		return err
	}

	// 保留不可由表单修改的字段，避免 Save 把它们清零
	sheet.CreatedAt = existing.CreatedAt
	sheet.Status = existing.Status
	sheet.ViewCount = existing.ViewCount
	sheet.UserID = existing.UserID

	if err := s.sheetRepo.Update(sheet); err != nil {
		return err
	}
	return s.syncTags(sheet.ID, tagNames)
}

func (s *SheetService) syncTags(sheetMusicID uint, tagNames []string) error {
	tagIDs := make([]uint, 0, len(tagNames))
	for _, name := range tagNames {
		if name == "" {
			continue
		}
		tag, err := s.tagRepo.FindOrCreate(name, "genre")
		if err != nil {
			s.logger.Error("find or create tag failed", zap.String("name", name), zap.Error(err))
			return errors.New("保存标签失败")
		}
		tagIDs = append(tagIDs, tag.ID)
	}
	if err := s.tagRepo.SetSheetTags(sheetMusicID, tagIDs); err != nil {
		s.logger.Error("set sheet tags failed", zap.Error(err))
		return errors.New("保存标签失败")
	}
	return nil
}

func (s *SheetService) Delete(id uint) error {
	_, err := s.sheetRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("乐谱不存在")
		}
		return err
	}
	return s.sheetRepo.Delete(id)
}
