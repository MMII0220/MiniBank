package controller

import (
	"fmt"
	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/MMII0220/MiniBank/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	// "time"
)

func blockUnblockAccountHandler(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(domain.User)
	if !currentUser.IsAdmin() {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	var req domain.ReqAdminAccountAction
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	accountID, err := strconv.Atoi(c.Param("id"))
	if err != nil || accountID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	block := req.Block

	// Controller передает только HTTP параметры в Service
	if err := service.BlockUnblockAccount(accountID, block, currentUser.ID, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("account %d %s", accountID, map[bool]string{true: "blocked", false: "unblocked"}[block])})
}

func getAuditLogsHandler(c *gin.Context) {
	logs, err := service.AuditLogs()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, logs)
}
