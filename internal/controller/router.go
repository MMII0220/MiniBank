package controller

import (
	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/gin-gonic/gin"
	"os"
)

func (ctr *Controller) SetupRoutes() {
	r := gin.Default()

	r.GET("/ping", ctr.healthCheck)

	auth := r.Group("/auth")
	{
		auth.POST("/register", ctr.registerHandler)
		auth.POST("/login", ctr.loginHandler)
	}

	admin := r.Group("/admin")
	admin.Use(ctr.AuthMiddleware(domain.RoleAdmin))
	{
		admin.POST("/blockUnblock/:id", ctr.blockUnblockAccountHandler)
		admin.GET("/getAuditLogs", ctr.getAuditLogsHandler)
	}

	api := r.Group("/api")
	api.Use(ctr.AuthMiddleware(domain.RoleUser)) // или можно создать специальную роль "any"
	{
		api.POST("/deposit", ctr.depositHandler)
		api.POST("/withdraw", ctr.withdrawHandler)
		api.POST("/transfer", ctr.transferHandler)
		api.GET("/history", ctr.historyLogs)
	}

	r.Run(":" + os.Getenv("ROUTER_RUN")) // listen and serve on
}
