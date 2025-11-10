package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/MMII0220/MiniBank/internal/domain/contracts"
	"github.com/MMII0220/MiniBank/internal/errs"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service contracts.ServiceI
}

func NewController(service contracts.ServiceI) *Controller {
	return &Controller{
		service: service,
	}
}

func (ctr *Controller) translateError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, errs.ErrDatabaseError):
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
	case errors.Is(err, errs.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	case errors.Is(err, errs.ErrUserAlreadyRegistered):
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
	case errors.Is(err, errs.ErrAccessDenied):
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
	case errors.Is(err, errs.ErrAccountBlocked):
		c.JSON(http.StatusForbidden, gin.H{"error": "Account is blocked"})
	case errors.Is(err, errs.ErrInsufficientFunds):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
	case errors.Is(err, errs.ErrInvalidAmount):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
	case errors.Is(err, errs.ErrDailyLimitExceeded):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Daily limit exceeded"})
	case errors.Is(err, errs.ErrInvalidToken):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
	case errors.Is(err, errs.ErrTokenExpired):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
	case errors.Is(err, errs.ErrRefreshTokenExpired):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
	case errors.Is(err, errs.ErrOperationNotAllowed):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Operation not allowed"})
	default:
		// Неизвестная ошибка - возвращаем 500 и логируем
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}

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

		// Проверяем роль если указана
		if requiredRole != "" && user.Role != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		// Сохраняем пользователя в контексте
		c.Set("currentUser", user)
		c.Next()
	}
}
