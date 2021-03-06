package app

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/lozhkindm/banking-auth/config"
	"github.com/lozhkindm/banking-auth/domain"
	"github.com/lozhkindm/banking-auth/service"
	"log"
	"net/http"
	"time"
)

func Start() {
	router := mux.NewRouter()
	dbClient := getDbClient()

	ah := AuthHandlers{
		service: service.NewAuthService(
			domain.NewAuthRepositoryDB(dbClient),
			domain.GetRolePermissions(),
		),
	}

	router.HandleFunc("/auth/login", ah.Login).Methods(http.MethodPost)
	//router.HandleFunc("/auth/register", ah.Register).Methods(http.MethodPost)
	router.HandleFunc("/auth/verify", ah.Verify).Methods(http.MethodGet)
	router.HandleFunc("/auth/refresh", ah.Refresh).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(config.NewServerConfig().AsString(), router))
}

func getDbClient() *sqlx.DB {
	client, err := sqlx.Open("mysql", config.NewDbConfig().AsDataSource())

	if err != nil {
		panic(err)
	}

	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)

	return client
}
