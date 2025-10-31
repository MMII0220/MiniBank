package service

import (
	"github.com/MMII0220/MiniBank/internal/domain/contracts"
)

type Service struct {
	repo contracts.RepositoryI
}

func NewService(repo contracts.RepositoryI) *Service {
	return &Service{
		repo: repo,
	}
}
