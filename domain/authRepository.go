package domain

import "github.com/jmoiron/sqlx"

type AuthRepository interface {
}

type AuthRepositoryDB struct {
	client *sqlx.DB
}

func NewAuthRepositoryDB(client *sqlx.DB) AuthRepositoryDB {
	return AuthRepositoryDB{client: client}
}
