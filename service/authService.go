package service

import "github.com/lozhkindm/banking-auth/domain"

type AuthService interface {
}

type DefaultAuthService struct {
	repo domain.AuthRepository
}

func NewAuthService(repo domain.AuthRepository) AuthService {
	return DefaultAuthService{repo: repo}
}
