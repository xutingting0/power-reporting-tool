// file: main.go
package main

import (
	"backend/api"
	"backend/db"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库连接
	db.Init()

	// 创建数据库表
	db.CreateTables()

	// 初始化数据
	db.InitData()

	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 公开路由
	r.POST("/api/login", api.Login)
	r.POST("/api/logout", api.Logout)

	// 需要认证的路由组
	authGroup := r.Group("/api")
	authGroup.Use(api.AuthMiddleware())
	{
		// 用户相关
		authGroup.GET("/users/current", api.GetCurrentUser)
		authGroup.PUT("/users/profile", api.UpdateUserProfile)

		// 报备单相关
		authGroup.POST("/reports", api.CreateReport)
		authGroup.GET("/reports/user", api.GetUserReports)
		authGroup.GET("/reports/pending", api.GetPendingReports)
		authGroup.GET("/reports/:id", api.GetReportByID)
		authGroup.PUT("/reports/:id", api.UpdateReport)
		authGroup.POST("/reports/:id/submit", api.SubmitReport)
		authGroup.DELETE("/reports/:id", api.DeleteReport)
		authGroup.POST("/reports/batch-delete", api.BatchDeleteReports)
		authGroup.GET("/reports/:id/logs", api.GetReportLogs)
		authGroup.GET("/reports/export", api.ExportReports)
		authGroup.GET("/reports/statistics", api.GetReportStatistics)
		authGroup.GET("/dashboard/stats", api.GetDashboardStats)
		authGroup.GET("/reports/recent", api.GetRecentReports)

		// 模板相关
		authGroup.GET("/templates", api.GetAllTemplates)
		authGroup.GET("/templates/:id", api.GetTemplateByID)
		authGroup.POST("/templates/:id/use", api.UseTemplate)

		// 配置相关
		authGroup.GET("/config/units", api.GetConfigUnits)
		authGroup.GET("/config/substations", api.GetConfigSubstations)
		authGroup.GET("/config/substations/grouped", api.GetSubstationsGrouped)
		authGroup.GET("/config/safety-templates", api.GetSafetyTemplates)
		authGroup.GET("/config/required-safety-measures", api.GetRequiredSafetyMeasures)
		authGroup.GET("/config/device-types", api.GetDeviceTypes)

		// 管理员路由组
		adminGroup := authGroup.Group("/admin")
		adminGroup.Use(api.AdminMiddleware())
		{
			// 用户管理
			adminGroup.GET("/users", api.GetAllUsers)
			adminGroup.POST("/users", api.CreateUser)
			adminGroup.PUT("/users/:username", api.UpdateUser)
			adminGroup.DELETE("/users/:username", api.DeleteUser)

			// 模板管理
			adminGroup.POST("/templates", api.CreateTemplate)
			adminGroup.PUT("/templates/:id", api.UpdateTemplate)
			adminGroup.DELETE("/templates/:id", api.DeleteTemplate)

			// 报备单管理
			adminGroup.GET("/reports/all", api.GetAllReports)
			adminGroup.POST("/reports/:id/approve", api.ApproveReport)
			adminGroup.POST("/reports/:id/reject", api.RejectReport)

			// 配置管理
			adminGroup.POST("/config/units", api.AddConfigUnit)
			adminGroup.PUT("/config/units/:id", api.UpdateConfigUnit)
			adminGroup.DELETE("/config/units/:id", api.DeleteConfigUnit)
			adminGroup.POST("/config/substations", api.AddConfigSubstation)
			adminGroup.PUT("/config/substations/:id", api.UpdateConfigSubstation)
			adminGroup.DELETE("/config/substations/:id", api.DeleteConfigSubstation)

			// 系统管理
			adminGroup.GET("/system/status", api.GetSystemStatus)
			adminGroup.POST("/system/backup", api.BackupDatabase)
		}
	}

	log.Println("服务器启动在 :8000")
	if err := r.Run(":8000"); err != nil {
		log.Fatal("启动服务器失败:", err)
	}
}
