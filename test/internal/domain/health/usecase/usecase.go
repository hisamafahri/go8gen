package usecase

import "abc/internal/domain/health"

type Health struct {
	healthRepo health.Repository
}

func NewHealthUseCase(health health.Repository) *Health {
	return &Health {
		healthRepo: health,
	}
}

func (u *Health) Readiness() error {
	return u.healthRepo.Readiness()
}