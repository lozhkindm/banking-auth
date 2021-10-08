package config

import (
	"fmt"
	"os"
)

type Db struct {
	User     string
	Password string
	Address  string
	Port     string
	Name     string
}

func NewDbConfig() Db {
	conf := Db{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWD"),
		Address:  os.Getenv("DB_ADDR"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
	}

	if conf.User == "" {
		conf.User = "root"
	}

	if conf.Password == "" {
		conf.Password = "codecamp"
	}

	if conf.Address == "" {
		conf.Address = "localhost"
	}

	if conf.Port == "" {
		conf.Port = "13306"
	}

	if conf.Name == "" {
		conf.Name = "banking"
	}

	return conf
}

func (db Db) AsDataSource() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db.User, db.Password, db.Address, db.Port, db.Name)
}
