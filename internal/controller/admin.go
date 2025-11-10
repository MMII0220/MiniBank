package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MMII0220/MiniBank/internal/controller/dto"
	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/gin-gonic/gin"
)

func (ctr *Controller) blockUnblockAccountHandler(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(domain.User)
	if !currentUser.IsAdmin() {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	var req dto.ReqAdminAccountActionHTTP
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
	if err := ctr.service.BlockUnblockAccount(accountID, block, currentUser.ID, req.Reason); err != nil {
		ctr.translateError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("account %d %s", accountID, map[bool]string{true: "blocked", false: "unblocked"}[block])})
}

func (ctr *Controller) getAuditLogsHandler(c *gin.Context) {
	logs, err := ctr.service.AuditLogs()
	if err != nil {
		ctr.translateError(c, err)
		return
	}
	c.JSON(200, logs)
}
