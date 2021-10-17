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

// Verify
// http://localhost:8181/auth/verify?token=somevalidtokenstring&routeName=GetCustomer&customer_id=2000&account_id=95470
func (h AuthHandlers) Verify(w http.ResponseWriter, r *http.Request) {
	urlParams := make(map[string]string)

	for k := range r.URL.Query() {
		urlParams[k] = r.URL.Query().Get(k)
	}

	if urlParams["token"] != "" {
		err := h.service.Verify(urlParams)

		if err != nil {
			writeResponse(w, err.Code, notAuthorizedResponse(err.Message))
		} else {
			writeResponse(w, http.StatusOK, authorizedResponse())
		}
	} else {
		writeResponse(w, http.StatusForbidden, notAuthorizedResponse("missing token"))
	}
}

func (h AuthHandlers) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		res, err := h.service.Refresh(req)

		if err != nil {
			writeResponse(w, err.Code, err.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, res)
		}
	}
}

func notAuthorizedResponse(msg string) map[string]interface{} {
	return map[string]interface{}{
		"is_authorized": false,
		"message":       msg,
	}
}

func authorizedResponse() map[string]bool {
	return map[string]bool{"is_authorized": true}
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
