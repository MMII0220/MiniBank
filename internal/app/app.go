package app

import (
	"log"

	"github.com/MMII0220/MiniBank/config"
	"github.com/MMII0220/MiniBank/internal/controller"
	"github.com/MMII0220/MiniBank/internal/redis"
	"github.com/MMII0220/MiniBank/internal/repository"
	"github.com/MMII0220/MiniBank/internal/service"
)

// AppRun starts the application in main.go
func AppRun() {
	dbConn, err := config.InitDB()
	if err != nil {
		log.Fatal("failed to initialize database: ", err)
	}
	defer config.CloseDB()

	if err := redis.InitRedisConnection(); err != nil {
		log.Printf("WARNING: Cannot connect to Redis: %v", err)
		log.Printf("Application will continue without Redis caching")
	} else {
		log.Println("Redis connected successfully")
	}

	// redisClient := redis.GetRedisClient()

	rep := repository.NewRepository(dbConn)
	svc := service.NewService(rep)
	ctr := controller.NewController(svc)

	ctr.SetupRoutes()
}
