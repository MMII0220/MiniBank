package controller

import (
	// "github.com/MMII0220/MiniBank/internal/domain"
	// "github.com/MMII0220/MiniBank/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// type Limit struct {
// 	ID            int64         `json:"id" db:"id"`
// 	UserID        int64         `json:"user_id" db:"user_id"`
// 	OperationType OperationType `json:"operation_type" db:"operation_type"`
// 	DailyAmount   float64       `json:"daily_amount" db:"daily_amount"`
// 	DailyCount    int           `json:"daily_count" db:"daily_count"`
// 	LastReset     time.Time     `json:"last_reset" db:"last_reset"`
// }

func getLimit(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Limit controller is working!"})
}

func updateLimit(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Update limit controller is working!"})
}
