package handler

import (
	"fmt"
	"net/http"
	"net/url"
)

/*
 * Generate Index (/www)
 */

type HandlerIndex struct{}

func (h *HandlerIndex) IsAble(url *url.URL) bool {
	if url.Path == "/" || url.Path == "/index" || url.Path == "/index.html" {
		return true
	}

	return false
}

func (h *HandlerIndex) Handle(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		http.Error(w, "405 MethodNotAllowed", http.StatusMethodNotAllowed)
		return nil
	}

	fmt.Fprintf(w, "index.html?")

	return nil
}
