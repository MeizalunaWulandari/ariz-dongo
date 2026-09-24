package service

import (
	"github.com/MeizalunaWulandari/ariz-dongo/internal/config"
	"github.com/MeizalunaWulandari/ariz-dongo/internal/repository"
)

type HealthService struct {
	config     config.Config
	repository *repository.Repository
}

type HealthResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Status  string `json:"status"`
}

func NewHealthService(
	cfg config.Config,
	repo *repository.Repository,
) *HealthService {
	return &HealthService{
		config:     cfg,
		repository: repo,
	}
}

func (s *HealthService) Check() HealthResponse {
	return HealthResponse{
		Name:    s.config.AppName,
		Version: s.config.AppVersion,
		Status:  "ok",
	}
}
