package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/config"
	"github.com/sheet-platform/backend/internal/model"
	"github.com/sheet-platform/backend/internal/pkg"
	"github.com/sheet-platform/backend/internal/router"
)

func main() {
	configPath := flag.String("config", "config.yaml", "config file path")
	flag.Parse()

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("init logger failed: %v", err)
	}
	defer logger.Sync()

	cfg := config.Init(*configPath)
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)


	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("connect database failed", zap.Error(err))
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.UserRole{},
		&model.Permission{},
		&model.PermissionModule{},
		&model.PermissionModuleRelation{},
		&model.RoleModule{},
		&model.SheetMusic{},
		&model.SheetFile{},
		&model.AudioFile{},
		&model.SheetComment{},
		&model.SheetLike{},
		&model.Tag{},
		&model.SheetTag{},
		&model.SheetOrder{},
		&model.SystemConfig{},
		&model.SystemLog{},
		&model.BannedIp{},
		&model.UserLoginLog{},
		&model.ContactMessage{},
		&model.FriendLink{},
		&model.DownloadRecord{},
		&model.SitePvHistory{},
	); err != nil {
		logger.Fatal("auto migrate failed", zap.Error(err))
	}

	if err := seedSystemConfigs(db); err != nil {
		logger.Fatal("seed system config failed", zap.Error(err))
	}

	// 每小时归档内存 PV 到 DB（启动时先做一次，捞回之前未归档的数据）
	pkg.ArchivePV(db)
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			pkg.ArchivePV(db)
		}
	}()

	r := gin.New()
	router.Setup(r, db, cfg, logger)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("server starting", zap.String("addr", addr))

	if err := r.Run(addr); err != nil {
		logger.Fatal("server start failed", zap.Error(err))
	}
}
