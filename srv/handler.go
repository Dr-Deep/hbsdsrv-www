package srv

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/Dr-Deep/hbsdsrv-www/handler"
)

/*
! ratelimit by nginx
*/
type Handler interface {
	IsAble(url *url.URL) bool
	Handle(w http.ResponseWriter, r *http.Request) error
}

func (www *WWWServer) Handle(w http.ResponseWriter, r *http.Request) {
	www.logger.Info(r.RemoteAddr, trunc(r.UserAgent()), "=>", r.Method, r.RequestURI)

	// Allowed Host filter
	if r.Host != www.cfg.Application.AllowedHost {
		http.Redirect(w, r, fmt.Sprintf("http://%s/", www.cfg.Application.AllowedHost), http.StatusMovedPermanently)
	}

	//
	var ourHandler Handler
	for _, h := range www.handlers {
		if h.IsAble(r.URL) {
			ourHandler = h
		}
	}

	// 404
	if ourHandler == nil {
		handler.Error(w, http.StatusNotFound)
		return
	}

	if err := ourHandler.Handle(w, r); err != nil {
		handler.Error(w, http.StatusInternalServerError)
		www.logger.Error(fmt.Sprintf("%v", ourHandler), err.Error())
	}
}
