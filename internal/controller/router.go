package controller

import (
	"os"

	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/gin-gonic/gin"
)

func (ctr *Controller) SetupRoutes() {
	r := gin.Default()

	r.GET("/ping", ctr.healthCheck)
	r.GET("/health/redis", ctr.redisHealth)

	auth := r.Group("/auth")
	{
		auth.POST("/register", ctr.registerHandler)
		auth.POST("/login", ctr.loginHandler)
		auth.POST("/refresh", ctr.refreshTokenHandler)
	}

	admin := r.Group("/admin")
	admin.Use(ctr.AuthMiddleware(domain.RoleAdmin))
	{
		admin.POST("/blockUnblock/:id", ctr.blockUnblockAccountHandler)
		admin.GET("/getAuditLogs", ctr.getAuditLogsHandler)
	}

	api := r.Group("/api")
	api.Use(ctr.AuthMiddleware(domain.RoleUser))
	{
		api.POST("/deposit", ctr.depositHandler)
		api.POST("/withdraw", ctr.withdrawHandler)
		api.POST("/transfer", ctr.transferHandler)
		api.GET("/history", ctr.historyLogs)
		api.GET("/accounts", ctr.getAllAccountsHandler)
	}

	r.Run(":" + os.Getenv("ROUTER_RUN"))
}
