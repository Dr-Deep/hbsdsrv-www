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

// generates www/index.md file
func genIndex(contentPaths map[string]string) error {
	var (
		index = strings.Builder{}
	)
	{ // Generate
		index.WriteString("# Index.md - Overview\n")

		for urlPath, fsPath := range contentPaths {
			index.WriteString(
				fmt.Sprintf("* [%s](%s)\n", fsPath, urlPath),
			)
		}
	}

	{ // store (should be public)
		//nolint:gosec
		if err := os.WriteFile(contentMdIndex, []byte(index.String()), 0444); err != nil { // read for everyone, write for none
			return err
		}
	}

	return nil
}

type HandlerContent struct {
	contentPaths map[string]string

	cfg    *config.Configuration
	logger *logging.Logger
}

func NewHandlerContent(logger *logging.Logger, cfg *config.Configuration) *HandlerContent {
	var (
		contentPaths map[string]string

		getContentPaths = func() {
			_contentPaths, err := genFsMap(cfg.Application.WWWDirectory, contentURL)
			if err != nil {
				logger.Error("gen Error", err.Error())
			}

			contentPaths = _contentPaths
		}
	)

	getContentPaths()

	if err := genIndex(contentPaths); err != nil {
		logger.Error("gen Error", err.Error())
	}

	getContentPaths()

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
