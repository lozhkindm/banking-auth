package dto

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/lozhkindm/banking-auth/domain"
)

type RefreshRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r RefreshRequest) IsAccessTokenValid() *jwt.ValidationError {
	_, err := jwt.Parse(r.AccessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HmacSampleSecret), nil
	})

	if err != nil {
		var ve *jwt.ValidationError

		if errors.As(err, &ve) {
			return ve
		}
	}

	return nil
}
