package service

import (
	"github.com/lozhkindm/banking-auth/domain"
	"github.com/lozhkindm/banking-auth/dto"
	"github.com/lozhkindm/banking-auth/errs"
)

type AuthService interface {
	Login(request dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
}

type DefaultAuthService struct {
	repo domain.AuthRepository
}

func (s DefaultAuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, *errs.AppError) {
	login, err := s.repo.FindByCredentials(req.Username, req.Password)

	if err != nil {
		return nil, err
	}

	token, err := login.GenerateToken()

	if err != nil {
		return nil, err
	}

	res := dto.LoginResponse{
		Token: *token,
	}

	return &res, nil
}

func NewAuthService(repo domain.AuthRepository) AuthService {
	return DefaultAuthService{repo: repo}
}
