package service

import (
	"fox/internal/modules/game/domain"
	"fox/internal/modules/game/repo"
)

type RNGFactory func() domain.RNG

type Service struct {
	repo repo.GameRepo
	rng  RNGFactory
}

func New(r repo.GameRepo, rng RNGFactory) *Service {
	return &Service{repo: r, rng: rng}
}
