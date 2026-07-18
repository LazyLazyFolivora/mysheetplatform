package router

import (
	"net/http"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	"github.com/sheet-platform/backend/internal/config"
	"github.com/sheet-platform/backend/internal/handler"
	"github.com/sheet-platform/backend/internal/middleware"
	"github.com/sheet-platform/backend/internal/pkg"
	"github.com/sheet-platform/backend/internal/repository"
	"github.com/sheet-platform/backend/internal/service"
)

func Setup(r *gin.Engine, db *gorm.DB, cfg *config.Config, logger *zap.Logger) {
	r.Use(middleware.CORS())
	r.Use(func(c *gin.Context) {
		pkg.IncrementPV()
		c.Next()
	})
	r.Use(gin.Recovery())
	if cfg.Server.Mode == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 阻止直接访问 PDF 文件，只能通过鉴权下载接口获取
	r.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/uploads/sheets_pdf/") {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	})

	// 上传文件静态服务（封面、图片、音频、HLS 分段）
	r.Static("/uploads", cfg.Upload.Dir)

	api := r.Group("/api")
	api.Use(middleware.IPBan(db))

	// Repos
	userRepo := repository.NewUserRepo(db)
	sheetRepo := repository.NewSheetMusicRepo(db)
	tagRepo := repository.NewTagRepo(db)
	fileRepo := repository.NewSheetFileRepo(db)
	permRepo := repository.NewPermissionRepo(db)
	roleRepo := repository.NewRoleRepo(db)
	orderRepo := repository.NewOrderRepo(db)
	dlRecordRepo := repository.NewDownloadRecordRepo(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg, logger, db)
	sheetService := service.NewSheetService(sheetRepo, tagRepo, fileRepo, logger)
	fileService := service.NewFileService(fileRepo, cfg, logger, db)
	permService := service.NewPermissionService(permRepo, roleRepo, logger)
	orderService := service.NewOrderService(db, sheetRepo, orderRepo, cfg, logger)
	adminService := service.NewAdminService(db, logger)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	sheetHandler := handler.NewSheetHandler(sheetService, db, userRepo, fileRepo, dlRecordRepo, cfg)
	fileHandler := handler.NewFileHandler(fileService)
	permHandler := handler.NewPermissionHandler(permService)
	orderHandler := handler.NewOrderHandler(orderService, cfg)
	adminHandler := handler.NewAdminHandler(adminService)

	// ===== Public routes =====
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/send-code", authHandler.SendCode)
	api.GET("/sheet-music", sheetHandler.List)
	api.GET("/sheet-music/:id", sheetHandler.Detail)
	api.GET("/sheet-music/:id/download", sheetHandler.Download)
	api.GET("/tags", sheetHandler.ListTags)
	api.GET("/friend-links", adminHandler.ListPublicFriendLinks)
	api.POST("/contact-messages", adminHandler.CreateContactMessage)
	api.POST("/alipay/notify", orderHandler.Notify)
	api.GET("/alipay/order/status", orderHandler.OrderStatus)
	api.GET("/orders/:order_no/pay", orderHandler.Pay)

	// ===== Authenticated routes =====
	auth := api.Group("", middleware.Auth(cfg, db))
	{
		auth.GET("/profile", authHandler.Profile)
		auth.POST("/orders", orderHandler.Create)
		auth.GET("/orders", orderHandler.List)
	}

	// ===== Admin routes =====
	admin := api.Group("/admin", middleware.AdminAuth(cfg, db))
	{
		admin.GET("/dashboard", adminHandler.Dashboard)
		admin.GET("/users", adminHandler.ListUsers)
		admin.POST("/users/:id/ban", adminHandler.BanUser)
		admin.POST("/users/:id/unban", adminHandler.UnbanUser)

		admin.POST("/sheet-music", sheetHandler.Create)
		admin.PUT("/sheet-music/:id", sheetHandler.Update)
		admin.DELETE("/sheet-music/:id", sheetHandler.Delete)
		admin.POST("/sheet-music/upload-pdf", fileHandler.UploadSheetPDF)
		admin.POST("/sheet-music/upload-audio", fileHandler.UploadAudio)
		admin.POST("/upload/image", fileHandler.UploadImage)

		admin.GET("/permissions/tree", permHandler.Tree)
		admin.POST("/permissions/modules", permHandler.CreateModule)
		admin.PUT("/permissions/modules/:id", permHandler.UpdateModule)
		admin.DELETE("/permissions/modules/:id", permHandler.DeleteModule)
		admin.GET("/permissions/modules/:id/permissions", permHandler.GetModulePermissions)
		admin.POST("/permissions/modules/:id/permissions", permHandler.AssignModulePermissions)

		admin.GET("/permissions", permHandler.ListPermissions)
		admin.POST("/permissions", permHandler.CreatePermission)
		admin.PUT("/permissions/:id", permHandler.UpdatePermission)
		admin.DELETE("/permissions/:id", permHandler.DeletePermission)

		admin.GET("/roles", permHandler.ListRoles)
		admin.POST("/roles", permHandler.CreateRole)
		admin.PUT("/roles/:id", permHandler.UpdateRole)
		admin.DELETE("/roles/:id", permHandler.DeleteRole)
		admin.GET("/roles/:id/modules", permHandler.GetRoleModules)
		admin.POST("/roles/:id/modules", permHandler.AssignRoleModules)

		admin.GET("/contact-messages", adminHandler.ListMessages)
		admin.POST("/contact-messages/:id/read", adminHandler.MarkMessageRead)
		admin.DELETE("/contact-messages/:id", adminHandler.DeleteMessage)

		admin.GET("/friend-links", adminHandler.ListFriendLinks)
		admin.POST("/friend-links", adminHandler.CreateFriendLink)
		admin.PUT("/friend-links/:id", adminHandler.UpdateFriendLink)
		admin.DELETE("/friend-links/:id", adminHandler.DeleteFriendLink)

		admin.GET("/banned-ip", adminHandler.ListBannedIPs)
		admin.POST("/banned-ip", adminHandler.CreateBannedIP)
		admin.DELETE("/banned-ip/:id", adminHandler.DeleteBannedIP)

		admin.GET("/system-config", adminHandler.ListSystemConfigs)
		admin.POST("/system-config", adminHandler.CreateSystemConfig)
		admin.PUT("/system-config", adminHandler.UpdateSystemConfig)
		admin.DELETE("/system-config/:key", adminHandler.DeleteSystemConfig)
	}
}
