package app

import (
	"encoding/json"
	"github.com/lozhkindm/banking-auth/dto"
	"github.com/lozhkindm/banking-auth/service"
	"net/http"
)

type AuthHandlers struct {
	service service.AuthService
}

func (h AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		res, err := h.service.Login(req)

		if err != nil {
			writeResponse(w, err.Code, err.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, res)
		}
	}
}

//func (h AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
//
//}
//
//func (h AuthHandlers) Verify(w http.ResponseWriter, r *http.Request) {
//
//}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
