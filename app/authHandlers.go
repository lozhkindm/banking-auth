package app

import (
	"github.com/lozhkindm/banking-auth/service"
	"net/http"
)

type AuthHandlers struct {
	service service.AuthService
}

func (h AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {

}

func (h AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {

}

func (h AuthHandlers) Verify(w http.ResponseWriter, r *http.Request) {

}
