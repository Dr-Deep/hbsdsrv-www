package handler

import (
	"net/http"
	"net/url"
	"os"

	"github.com/Dr-Deep/hbsdsrv-www/config"
	"github.com/Dr-Deep/logging-go"
)

/*
 * Generate Index (/www)
 */

type HandlerIndex struct {
	indexSite []byte

	cfg    *config.Configuration
	logger *logging.Logger
}

func NewHandlerIndex(logger *logging.Logger, cfg *config.Configuration) *HandlerIndex {
	indexSite, err := os.ReadFile(indexHTML)
	if err != nil {
		panic(err)
	}

	return &HandlerIndex{
		indexSite: indexSite,
		cfg:       cfg,
		logger:    logger,
	}
}

func (h *HandlerIndex) IsAble(url *url.URL) bool {
	if url.Path == "/" || url.Path == "/index" || url.Path == "/index.html" {
		return true
	}

	return false
}

func (h *HandlerIndex) Handle(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		Error(w, http.StatusMethodNotAllowed)
		return nil
	}

	w.Write(h.indexSite)

	return nil
}
