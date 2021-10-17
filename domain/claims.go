package domain

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const HmacSampleSecret = "hmacSampleSecret"
const AccessTokenDuration = time.Hour
const RefreshTokenDuration = time.Hour * 24 * 30

type AccessClaims struct {
	CustomerId string   `json:"customer_id"`
	Accounts   []string `json:"accounts"`
	Username   string   `json:"username"`
	Role       string   `json:"role"`
	jwt.StandardClaims
}

type RefreshClaims struct {
	TokenType  string   `json:"token_type"`
	CustomerId string   `json:"cid"`
	Accounts   []string `json:"accounts"`
	Username   string   `json:"un"`
	Role       string   `json:"role"`
	jwt.StandardClaims
}

func (a AccessClaims) IsUserRole() bool {
	return a.Role == "user"
}

func (a AccessClaims) IsValidCustomerId(customerId string) bool {
	return a.CustomerId == customerId
}

func (a AccessClaims) IsValidAccountId(accountId string) bool {
	if accountId != "" {
		accountFound := false

		for _, a := range a.Accounts {
			if a == accountId {
				accountFound = true
				break
			}
		}

		return accountFound
	}

	return true
}

func (a AccessClaims) IsRequestVerifiedWithTokenClaims(urlParams map[string]string) bool {
	if a.CustomerId != urlParams["customer_id"] {
		return false
	}

	if !a.IsValidAccountId(urlParams["account_id"]) {
		return false
	}

	return true
}

func (a AccessClaims) RefreshTokenClaims() RefreshClaims {
	return RefreshClaims{
		TokenType:  "refresh_token",
		CustomerId: a.CustomerId,
		Accounts:   a.Accounts,
		Username:   a.Username,
		Role:       a.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(RefreshTokenDuration).Unix(),
		},
	}
}

func (r RefreshClaims) AccessTokenClaims() AccessClaims {
	return AccessClaims{
		CustomerId: r.CustomerId,
		Accounts:   r.Accounts,
		Username:   r.Username,
		Role:       r.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AccessTokenDuration).Unix(),
		},
	}
}
