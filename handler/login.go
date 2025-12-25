package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/Dr-Deep/hbsdsrv-www/config"
	"github.com/Dr-Deep/logging-go"
)

//login => auth
//https://hbsdsrv.1337.cx/login.html

type HandlerLogin struct {
	loginSite []byte
	//creds

	cfg    *config.Configuration
	logger *logging.Logger
}

func NewHandlerLogin(logger *logging.Logger, cfg *config.Configuration) *HandlerLogin {
	loginSite, err := os.ReadFile(loginHTML)
	if err != nil {
		panic(err)
	}

	return &HandlerLogin{
		loginSite: loginSite,
		cfg:       cfg,
		logger:    logger,
	}
}

func (h *HandlerLogin) IsAble(url *url.URL) bool {
	return (url.Path == "/login")
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
