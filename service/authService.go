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
}

type DefaultAuthService struct {
	repo            domain.AuthRepository
	rolePermissions domain.RolePermissions
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

func (s DefaultAuthService) Verify(urlParams map[string]string) *errs.AppError {
	if jwtToken, err := jwtTokenFromString(urlParams["token"]); err != nil {
		return err
	} else {
		if jwtToken.Valid {
			claims := jwtToken.Claims.(*domain.Claims)

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
	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
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
