package controller

import (
	"net/http"

	"github.com/MMII0220/MiniBank/internal/controller/dto"
	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/MMII0220/MiniBank/internal/logger"
	"github.com/gin-gonic/gin"
)

func (ctr *Controller) registerHandler(c *gin.Context) {
	log := logger.GetLogger()
	log.Info().Str("endpoint", "register").Msg("Registration request received")

	var req dto.ReqRegisterHTTP
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid registration request format")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var role domain.Role
	if req.Role != "" && (req.Role == "admin" || req.Role == "user") {
		role = domain.Role(req.Role)
	} else {
		role = domain.RoleUser
	}

	domainReq := req.ToDomain()
	user, err := ctr.service.Register(domainReq, role)
	if err != nil {
		ctr.translateError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": user.ID})
}

func (ctr *Controller) loginHandler(c *gin.Context) {
	var req dto.ReqLoginHTTP
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domainReq := req.ToDomain()
	tokenResponse, err := ctr.service.Login(domainReq)
	if err != nil {
		ctr.translateError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}

func (ctr *Controller) refreshTokenHandler(c *gin.Context) {
	var req dto.ReqRefreshTokenHTTP
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domainReq := req.ToDomain()
	tokenResponse, err := ctr.service.RefreshToken(domainReq)
	if err != nil {
		ctr.translateError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}
