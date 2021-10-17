package service

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/lozhkindm/banking-auth/domain"
	"github.com/lozhkindm/banking-auth/dto"
	"github.com/lozhkindm/banking-lib/errs"
	"github.com/lozhkindm/banking-lib/logger"
)

type AuthService interface {
	Login(request dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
	Verify(urlParams map[string]string) *errs.AppError
	Refresh(request dto.RefreshRequest) (*dto.LoginResponse, *errs.AppError)
}

type DefaultAuthService struct {
	repo            domain.AuthRepository
	rolePermissions domain.RolePermissions
}

func (s DefaultAuthService) Refresh(req dto.RefreshRequest) (*dto.LoginResponse, *errs.AppError) {
	if ve := req.IsAccessTokenValid(); ve != nil {
		if ve.Errors != jwt.ValidationErrorExpired {
			return nil, errs.NewUnauthorizedError("invalid access token")
		}

		if err := s.repo.RefreshTokenExists(req.RefreshToken); err != nil {
			return nil, err
		}

		if accessToken, err := domain.NewAccessTokenFromRefreshToken(req.RefreshToken); err != nil {
			return nil, err
		} else {
			return &dto.LoginResponse{AccessToken: accessToken}, nil
		}
	}

	return nil, errs.NewUnauthorizedError("current access token is not expired")
}

func (s DefaultAuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, *errs.AppError) {
	var login *domain.Login
	var err *errs.AppError
	var accessToken, refreshToken string

	if login, err = s.repo.FindByCredentials(req.Username, req.Password); err != nil {
		return nil, err
	}

	if accessToken, err = login.GenerateAccessToken(); err != nil {
		return nil, err
	}

	if refreshToken, err = login.GenerateRefreshToken(); err != nil {
		return nil, err
	}

	if err = s.repo.SaveRefreshToken(refreshToken); err != nil {
		return nil, err
	}

	res := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return &res, nil
}

func (s DefaultAuthService) Verify(urlParams map[string]string) *errs.AppError {
	if jwtToken, err := jwtTokenFromString(urlParams["token"]); err != nil {
		return err
	} else {
		if jwtToken.Valid {
			claims := jwtToken.Claims.(*domain.AccessClaims)

			if claims.IsUserRole() {
				if !claims.IsRequestVerifiedWithTokenClaims(urlParams) {
					return errs.NewUnauthorizedError("request not verified with the token claims")
				}
			}

			isAuthorized := s.rolePermissions.IsAuthorizedFor(claims.Role, urlParams["route_name"])

			if !isAuthorized {
				return errs.NewForbiddenError(fmt.Sprintf("%s role is not authorized", claims.Role))
			}

			return nil
		} else {
			return errs.NewUnauthorizedError("invalid token")
		}
	}
}

func jwtTokenFromString(tokenString string) (*jwt.Token, *errs.AppError) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HmacSampleSecret), nil
	})

	if err != nil {
		logger.Error("Error while parsing token: " + err.Error())
		return nil, errs.NewUnauthorizedError("Cannot parse the token")
	}

	return token, nil
}

func NewAuthService(repo domain.AuthRepository, perms domain.RolePermissions) AuthService {
	return DefaultAuthService{
		repo:            repo,
		rolePermissions: perms,
	}
}
