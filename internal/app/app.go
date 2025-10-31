package app

import (
	"log"

	"github.com/MMII0220/MiniBank/config"
	"github.com/MMII0220/MiniBank/internal/controller"
)

// AppRun starts the application in main.go
func AppRun() {
	if err := config.InitDB(); err != nil {
		log.Fatal("failed to initialize database: ", err)
	}
	defer config.CloseDB()

	controller.SetupRoutes()

	//init cron(interface service)
}
