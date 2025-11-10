package controller

import (
	"net/http"

	"github.com/MMII0220/MiniBank/internal/controller/dto"
	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/MMII0220/MiniBank/internal/logger"
	appredis "github.com/MMII0220/MiniBank/internal/redis"
	"github.com/gin-gonic/gin"
)

// Checking if there is a connection to the server
func (ctr *Controller) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ping": "pong",
	})
}

// Redis health check
func (ctr *Controller) redisHealth(c *gin.Context) {
	if err := appredis.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"redis": "down",
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"redis": "ok",
	})
}

// Adding money to bank-account
func (ctr *Controller) depositHandler(c *gin.Context) {
	log := logger.GetLogger()
	currentUser := c.MustGet("currentUser").(domain.User)

	log.Info().
		Int("user_id", currentUser.ID).
		Str("endpoint", "deposit").
		Msg("Deposit request received")

	var req dto.ReqTransactionHTTP
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Warn().Err(err).Int("user_id", currentUser.ID).Msg("Invalid deposit request format")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.CardNumber == "" && req.PhoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "either card_number or phone_number must be provided"})
		return
	}

	domainReq := req.ToDomain()
	err = ctr.service.Deposit(int(currentUser.ID), domainReq)
	if err != nil {
		ctr.translateError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deposit successful"})
}

// Withdrawing money from bank-account
func (ctr *Controller) withdrawHandler(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(domain.User)

	var req dto.ReqTransactionHTTP
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.CardNumber == "" && req.PhoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "either card_number or phone_number must be provided"})
		return
	}

	domainReq := req.ToDomain()
	err = ctr.service.Withdraw(int(currentUser.ID), domainReq)
	if err != nil {
		ctr.translateError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Withdraw successful"})
}

// Transferring money between bank-accounts
func (ctr *Controller) transferHandler(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(domain.User)

	var req dto.ReqTransferHTTP
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.FromCardNumber == "" && req.FromPhoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from_card_number or from_phone_number must be provided"})
		return
	}

	if req.ToCardNumber == "" && req.ToPhoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to_card_number or to_phone_number must be provided"})
		return
	}

	domainReq := req.ToDomain()
	err := ctr.service.Transfer(int(currentUser.ID), domainReq)
	if err != nil {
		ctr.translateError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer successful"})
}

func (ctr *Controller) historyLogs(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(domain.User)

	transactions, err := ctr.service.HistoryLogs(currentUser.ID)
	if err != nil {
		ctr.translateError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"history_logs": transactions})
}

func (ctr *Controller) getAllAccountsHandler(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(domain.User)

	accounts, err := ctr.service.GetAllAccounts(currentUser.ID)
	if err != nil {
		ctr.translateError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"accounts":    accounts,
		"total_count": len(accounts),
	})
}
