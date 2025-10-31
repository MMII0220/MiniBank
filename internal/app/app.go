package app

import (
	"log"

	"github.com/MMII0220/MiniBank/config"
	"github.com/MMII0220/MiniBank/internal/controller"
	"github.com/MMII0220/MiniBank/internal/repository"
	"github.com/MMII0220/MiniBank/internal/service"
)

// dbConn, err := config.InitDB()
// 	if err != nil {

// AppRun starts the application in main.go
func AppRun() {
	dbConn, err := config.InitDB()
	if err != nil {
		log.Fatal("failed to initialize database: ", err)
	}
	defer config.CloseDB()

	rep := repository.NewRepository(dbConn)
	svc := service.NewService(rep)
	ctr := controller.NewController(svc)

	ctr.SetupRoutes()

	//init cron(interface service)
}
