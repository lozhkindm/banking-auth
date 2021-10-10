package domain

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lozhkindm/banking-auth/errs"
	"github.com/lozhkindm/banking-auth/logger"
)

type AuthRepository interface {
	FindByCredentials(username, password string) (*Login, *errs.AppError)
}

type AuthRepositoryDB struct {
	client *sqlx.DB
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
			return nil, errs.NewAuthorizationError("invalid credentials")
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
