package controller

import (
	"net/http"
	"strings"

	"github.com/MMII0220/MiniBank/internal/controller/dto"
	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/gin-gonic/gin"
)

func (ctr *Controller) registerHandler(c *gin.Context) {
	var req dto.ReqRegisterHTTP
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Определяем роль: если указана в запросе, используем её, иначе - user по умолчанию
	var role domain.Role
	if req.Role != "" && (req.Role == "admin" || req.Role == "user") {
		role = domain.Role(req.Role)
	} else {
		role = domain.RoleUser // по умолчанию
	}

	domainReq := req.ToDomain()
	user, err := ctr.service.Register(domainReq, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	token, err := ctr.service.Login(domainReq)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Middleware для проверки JWT и роли
func (ctr *Controller) AuthMiddleware(requiredRole domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}

		tokenStr := parts[1]
		user, err := ctr.service.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		if requiredRole != "" && user.Role != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		// Сохраняем пользователя в контекст
		c.Set("currentUser", user)
		c.Next()
	}
}
