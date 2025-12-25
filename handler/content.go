/*
 * Serve Markdown Content
 */
package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Dr-Deep/hbsdsrv-www/config"
	"github.com/Dr-Deep/logging-go"
)

type HandlerContent struct {
	contentPaths map[string]string

	cfg    *config.Configuration
	logger *logging.Logger
}

func NewHandlerContent(logger *logging.Logger, cfg *config.Configuration) *HandlerContent {
	contentPaths, err := genFsMap(cfg.Application.WWWDirectory, contentURL)
	if err != nil {
		logger.Error("gen Error", err.Error())
	}

	return &HandlerContent{
		contentPaths: contentPaths,
		cfg:          cfg,
		logger:       logger,
	}
}

func (h *HandlerContent) IsAble(url *url.URL) bool {
	if strings.HasPrefix(url.Path, contentURL) {
		return true
	}

	return false
}

func (h *HandlerContent) Handle(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		Error(w, http.StatusMethodNotAllowed)
		return nil
	}

	// valid path?
	fsPath, oke := h.contentPaths[r.URL.Path]
	if !oke {
		Error(w, http.StatusNotFound)
		return nil
	}

	// read file
	//nolint:gosec
	markdownContent, err := os.ReadFile(fsPath) // file inclusion is the Goal
	if err != nil {
		Error(w, http.StatusInternalServerError)
		h.logger.Error("file reading error", err.Error())
		return nil
	}

	// markdown to html
	htmlContent := template.HTML(
		renderMarkdownToHTML(markdownContent),
	)

	site, err := renderHTMLTemplate(
		baseHTML,
		struct {
			Title   string
			Content template.HTML
		}{
			Title:   fmt.Sprintf("hbsdsrv - %s", r.URL.Path),
			Content: htmlContent,
		},
	)
	if err != nil {
		h.logger.Error("template error", err.Error())
		Error(w, http.StatusInternalServerError)
		return nil
	}

	// respond
	w.Write([]byte(site))

	return nil
}
