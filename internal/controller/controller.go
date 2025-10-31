package controller

import (
	"github.com/MMII0220/MiniBank/internal/domain/contracts"
	// "github.com/MMII0220/MiniBank/internal/service"
)

type Controller struct {
	service contracts.ServiceI
}

func NewController(service contracts.ServiceI) *Controller {
	return &Controller{
		service: service,
	}
}
