package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

//login => auth
//https://hbsdsrv.1337.cx/login.html

type HandlerLogin struct {
	loginSite []byte
	//creds
}

func NewHandlerLogin() *HandlerLogin {
	loginSite, err := os.ReadFile(loginHTML)
	if err != nil {
		panic(err)
	}

	return &HandlerLogin{
		loginSite: loginSite,
	}
}

func (h *HandlerLogin) IsAble(url *url.URL) bool {
	if url.Path == "/login" {
		return true
	}

	return false
}

func (h *HandlerLogin) Handle(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		w.Write(h.loginSite)
		return nil

	case http.MethodPost: //login incoming
		fmt.Fprintf(w, `{ message:"wip" }`)
		return nil

	default:
		Error(w, http.StatusMethodNotAllowed)
		return nil
	}
}
