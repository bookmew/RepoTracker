package service

import (
	"RepoTracker/src/entity"
	"RepoTracker/src/repository"
)

type TokenService struct {
	repo *repository.TokenRepository
}

func NewTokenService(repo *repository.TokenRepository) *TokenService {
	return &TokenService{repo: repo}
}

func (s *TokenService) GetAll() ([]entity.Token, error) {
	return s.repo.GetAll()
}
