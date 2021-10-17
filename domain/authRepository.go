package domain

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lozhkindm/banking-lib/errs"
	"github.com/lozhkindm/banking-lib/logger"
)

type AuthRepository interface {
	FindByCredentials(username, password string) (*Login, *errs.AppError)
	SaveRefreshToken(refreshToken string) *errs.AppError
	RefreshTokenExists(refreshToken string) *errs.AppError
}

type AuthRepositoryDB struct {
	client *sqlx.DB
}

func (d AuthRepositoryDB) RefreshTokenExists(refreshToken string) *errs.AppError {
	var token string

	sqlSelect := "SELECT refresh_token FROM refresh_token_store WHERE refresh_token = ?"

	if err := d.client.Get(&token, sqlSelect, refreshToken); err == nil {
		return nil
	} else if err == sql.ErrNoRows {
		return errs.NewUnauthorizedError("Refresh token not found")
	} else {
		logger.Error("Error while getting refresh token: " + err.Error())
		return errs.NewDatabaseError()
	}
}

func (d AuthRepositoryDB) SaveRefreshToken(refreshToken string) *errs.AppError {
	sqlInsert := "INSERT INTO refresh_token_store (refresh_token) VALUES (?)"

	_, err := d.client.Exec(sqlInsert, refreshToken)

	if err != nil {
		logger.Error("Error while creating a new refresh token: " + err.Error())
		return errs.NewDatabaseError()
	}

	return nil
}

func (d AuthRepositoryDB) FindByCredentials(username, password string) (*Login, *errs.AppError) {
	var login Login

	sqlVerify := `SELECT username, u.customer_id, role, group_concat(a.account_id) as account_numbers FROM users u
                  LEFT JOIN accounts a ON a.customer_id = u.customer_id
                WHERE username = ? and password = ?
                GROUP BY a.customer_id`

	err := d.client.Get(&login, sqlVerify, username, password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewUnauthorizedError("invalid credentials")
		} else {
			logger.Error("Error while finding user by credentials: " + err.Error())
			return nil, errs.NewDatabaseError()
		}
	}

	return &login, nil
}

func NewAuthRepositoryDB(client *sqlx.DB) AuthRepositoryDB {
	return AuthRepositoryDB{client: client}
}
