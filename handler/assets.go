package handler

import (
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Dr-Deep/hbsdsrv-www/config"
	"github.com/Dr-Deep/logging-go"
)

/*
 * Static Assets
 */

type HandlerAssets struct {
	assetPaths map[string]string
	logger     *logging.Logger
}

func NewHandlerAssets(logger *logging.Logger, cfg *config.Configuration) *HandlerAssets {
	assetsPaths, err := gen(cfg.Application.AssetsDirectory)
	if err != nil {
		logger.Error("gen Error", err.Error())
	}

	return &HandlerAssets{
		assetPaths: assetsPaths,
	}
}

func (h *HandlerAssets) IsAble(url *url.URL) bool {
	if strings.HasPrefix(url.Path, assetsURL) {
		return true
	}

	return false
}

func (h *HandlerAssets) Handle(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		Error(w, http.StatusMethodNotAllowed)
		return nil
	}

	// valid path?
	fsPath, oke := h.assetPaths[r.URL.Path]
	if !oke {
		Error(w, http.StatusNotFound)
		return nil
	}

	// read file
	//nolint:gosec
	b, err := os.ReadFile(fsPath) // file inclusion is the Goal
	if err != nil {
		Error(w, http.StatusInternalServerError)
		//? log
		return nil
	}

	w.Write(b)

	return nil
}
