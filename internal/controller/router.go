package controller

import (
	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/gin-gonic/gin"
	"os"
)

func SetupRoutes() {
	r := gin.Default()

	r.GET("/ping", healthCheck)

	auth := r.Group("/auth")
	{
		auth.POST("/register", registerHandler)
		auth.POST("/login", loginHandler)
	}

	admin := r.Group("/admin")
	admin.Use(AuthMiddleware(domain.RoleAdmin))
	{
		admin.POST("/blockUnblock/:id", blockUnblockAccountHandler)
		admin.GET("/getAuditLogs", getAuditLogsHandler)
	}

	api := r.Group("/api")
	api.Use(AuthMiddleware(domain.RoleUser)) // или можно создать специальную роль "any"
	{
		api.POST("/deposit", depositHandler)
		api.POST("/withdraw", withdrawHandler)
		api.POST("/transfer", transferHandler)
		api.GET("/history", historyLogs)
	}

	r.Run(":" + os.Getenv("ROUTER_RUN")) // listen and serve on
}
