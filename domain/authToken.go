package domain

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/lozhkindm/banking-lib/errs"
	"github.com/lozhkindm/banking-lib/logger"
)

type AuthToken struct {
	token *jwt.Token
}

func (at AuthToken) NewAccessToken() (string, *errs.AppError) {
	signedToken, err := at.token.SignedString([]byte(HmacSampleSecret))

	if err != nil {
		logger.Error("Error while signing the access token: " + err.Error())
		return "", errs.NewUnauthorizedError("Cannot sign the access token")
	}

	return signedToken, nil
}

func (at AuthToken) newRefreshToken() (string, *errs.AppError) {
	accessClaims := at.token.Claims.(AccessClaims)
	refreshClaims := accessClaims.RefreshTokenClaims()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	signedToken, err := token.SignedString([]byte(HmacSampleSecret))

	if err != nil {
		logger.Error("Error while signing the refresh token: " + err.Error())
		return "", errs.NewUnauthorizedError("Cannot sign the refresh token")
	}

	return signedToken, nil
}

func NewAuthToken(claims AccessClaims) AuthToken {
	return AuthToken{
		token: jwt.NewWithClaims(jwt.SigningMethodHS256, claims),
	}
}

func NewAccessTokenFromRefreshToken(refreshToken string) (string, *errs.AppError) {
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(HmacSampleSecret), nil
	})

	if err != nil {
		return "", errs.NewUnauthorizedError("invalid or expired refresh token")
	}

	refreshClaims := token.Claims.(*RefreshClaims)
	accessClaims := refreshClaims.AccessTokenClaims()

	return NewAuthToken(accessClaims).NewAccessToken()
}
