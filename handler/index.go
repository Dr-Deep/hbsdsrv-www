package handler

import (
	"net/http"
	"net/url"
	"os"
)

/*
 * Generate Index (/www)
 */

type HandlerIndex struct {
	indexSite []byte
}

func NewHandlerIndex() *HandlerIndex {
	indexSite, err := os.ReadFile(indexHTML)
	if err != nil {
		panic(err)
	}

	return &HandlerIndex{
		indexSite: indexSite,
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
