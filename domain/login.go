package domain

import (
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"github.com/lozhkindm/banking-lib/errs"
	"strings"
	"time"
)

type Login struct {
	Username   string         `db:"username"`
	CustomerId sql.NullString `db:"customer_id"`
	Accounts   sql.NullString `db:"account_numbers"`
	Role       string         `db:"role"`
}

func (l Login) GenerateAccessToken() (string, *errs.AppError) {
	authToken := NewAuthToken(l.ClaimsForAccessToken())
	return authToken.NewAccessToken()
}

func (l Login) GenerateRefreshToken() (string, *errs.AppError) {
	authToken := NewAuthToken(l.ClaimsForAccessToken())
	return authToken.newRefreshToken()
}

func (l Login) ClaimsForAccessToken() AccessClaims {
	if l.Accounts.Valid && l.CustomerId.Valid {
		return l.claimsForUser()
	} else {
		return l.claimsForAdmin()
	}
}

func (l Login) claimsForUser() AccessClaims {
	accounts := strings.Split(l.Accounts.String, ",")

	return AccessClaims{
		CustomerId: l.CustomerId.String,
		Accounts:   accounts,
		Username:   l.Username,
		Role:       l.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AccessTokenDuration).Unix(),
		},
	}
}

func (l Login) claimsForAdmin() AccessClaims {
	return AccessClaims{
		Role:     l.Role,
		Username: l.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AccessTokenDuration).Unix(),
		},
	}
}
