package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/config"
	"github.com/sheet-platform/backend/internal/handler"
	"github.com/sheet-platform/backend/internal/middleware"
	"github.com/sheet-platform/backend/internal/model"
	"github.com/sheet-platform/backend/internal/pkg"
	"github.com/sheet-platform/backend/internal/repository"
	"github.com/sheet-platform/backend/internal/scheduler"
	"github.com/sheet-platform/backend/internal/seed"
	"github.com/sheet-platform/backend/internal/service"
)

func New(configPath string) *fx.App {
	return fx.New(
		// ---- infrastructure ----
		fx.Provide(newLogger),
		fx.Provide(func(logger *zap.Logger) *config.Config {
			return config.Init(configPath)
		}),
		fx.Provide(newDB),
		fx.Provide(newEngine),

		// ---- repositories ----
		fx.Provide(
			repository.NewUserRepo,
			repository.NewSheetMusicRepo,
			repository.NewTagRepo,
			repository.NewSheetFileRepo,
			repository.NewPermissionRepo,
			repository.NewRoleRepo,
			repository.NewOrderRepo,
			repository.NewDownloadRecordRepo,
			repository.NewContactMessageRepo,
			repository.NewFriendLinkRepo,
			repository.NewBannedIPRepo,
			repository.NewSystemConfigRepo,
		),

		// ---- services ----
		fx.Provide(
			service.NewAuthService,
			service.NewAdminService,
			service.NewPermissionService,
			service.NewFileService,
			service.NewOrderService,
			service.NewSheetService,
		),

		// ---- handlers ----
		fx.Provide(
			handler.NewAuthHandler,
			handler.NewAdminHandler,
			handler.NewFileHandler,
			handler.NewPermissionHandler,
			handler.NewOrderHandler,
			handler.NewSheetHandler,
		),

		// ---- side effects ----
		fx.Invoke(migrateAndSeed),
		fx.Invoke(setupRoutes),
		fx.Invoke(invokeScheduler),
		fx.Invoke(invokeServer),
	)
}

// ---- provider functions ----

func newLogger() (*zap.Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}
	zap.ReplaceGlobals(logger)
	return logger, nil
}

func newDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}
	return db, nil
}

func newEngine(cfg *config.Config) *gin.Engine {
	if cfg.Server.Mode == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	return gin.New()
}

// ---- invoke functions ----

func migrateAndSeed(db *gorm.DB) error {
	if err := db.AutoMigrate(model.AllModels()...); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	return seed.SystemConfigs(db)
}

func setupRoutes(
	engine *gin.Engine,
	db *gorm.DB,
	cfg *config.Config,
	authHandler *handler.AuthHandler,
	sheetHandler *handler.SheetHandler,
	fileHandler *handler.FileHandler,
	orderHandler *handler.OrderHandler,
	permHandler *handler.PermissionHandler,
	adminHandler *handler.AdminHandler,
) {
	engine.Use(middleware.CORS())
	engine.Use(func(c *gin.Context) {
		pkg.IncrementPV()
		c.Next()
	})
	engine.Use(gin.Recovery())
	engine.Use(middleware.BlockPDF())

	engine.Static("/uploads", cfg.Upload.Dir)

	api := engine.Group("/api")
	api.Use(middleware.IPBan(db))

	public := api.Group("")
	auth := api.Group("", middleware.Auth(cfg, db))
	admin := api.Group("/admin", middleware.AdminAuth(cfg, db))

	authHandler.RegisterRoutes(public, auth, admin)
	sheetHandler.RegisterRoutes(public, auth, admin)
	fileHandler.RegisterRoutes(public, auth, admin)
	orderHandler.RegisterRoutes(public, auth, admin)
	permHandler.RegisterRoutes(public, auth, admin)
	adminHandler.RegisterRoutes(public, auth, admin)
}

func invokeScheduler(lc fx.Lifecycle, db *gorm.DB, logger *zap.Logger) {
	stop := scheduler.Start(db, logger)
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			stop()
			return nil
		},
	})
}

func invokeServer(lc fx.Lifecycle, engine *gin.Engine, cfg *config.Config, logger *zap.Logger) {
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{Addr: addr, Handler: engine}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("server starting", zap.String("addr", addr))
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("server shutting down")
			return srv.Shutdown(ctx)
		},
	})
}
