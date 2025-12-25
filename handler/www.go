package handler

import (
	"fmt"
	"net/http"
	"net/url"
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
	return &HandlerContent{}
}

func (h *HandlerContent) IsAble(url *url.URL) bool {
	if strings.HasPrefix(url.Path, "/www") {
		return true
	}

	return false
}

/*
generate index
*/
func (h *HandlerContent) Handle(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "wip")
	return nil
}

//www?q=contentpath/id

//  fsys fs.FS
/*
	hasFirst := r.URL.Query().Has("first")
	first := r.URL.Query().Get("first")
	hasSecond := r.URL.Query().Has("second")
	second := r.URL.Query().Get("second")
*/

/*
import "html/template"

- Markdown:




*/

/*

func genIndex(files []string) {
	var (
		indexFilePath = *contentDirPath + "/index.md"
		index         = strings.Builder{}
	)

	index.WriteString("# Index.md - Overview\n")

	for _, file := range files {
		index.WriteString(
			fmt.Sprintf("* [%s](%s)\n", file, file),
		)

		logger.Debug(file) // prefix?
	}

	//nolint:gosec // file should be public
	if err := os.WriteFile(indexFilePath, []byte(index.String()), 0444); err != nil { // read for everyone, write for none
		panic(err)
	}

	logger.Debug("generated index file", *indexHTMLSitePath)
}

// func(w http.ResponseWriter, r *http.Request)
func indexHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info(r.RemoteAddr, "=>", r.Method, r.URL.Path)

	// valid path?
	switch r.URL.Path {
	case "/":
		// redir to index
		logger.Info("redirecting ip from '/' to '/index.md'") // PLUS PREFIX
		http.Redirect(w, r, "/index.md", http.StatusPermanentRedirect)

	default:
		// Immer zuerst Status setzen, dann schreiben!
		w.WriteHeader(http.StatusNotFound)
		http.Redirect(w, r, *notFoundHTMLURI, http.StatusPermanentRedirect)
	}
}

// contentHandler()
func contentHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info(r.RemoteAddr, "=>", r.Method, r.URL.Path)

	// need okay from uris
	_, oke := uris[r.URL.Path]
	if !oke {
		http.Redirect(w, r, *notFoundHTMLURI, http.StatusPermanentRedirect)

		return
	}

	// Read File
	//nolint:gosec // file inclusion is the Goal
	var path = *contentDirPath + r.URL.Path

	content, err := os.ReadFile(path)
	if err != nil {
		logger.Error("ReadFile err from", *contentDirPath+r.URL.Path, err.Error())
		http.Redirect(w, r, *notFoundHTMLURI, http.StatusPermanentRedirect)

		return
	}

	htmlContent := renderMarkdownToHTML(content)

	// Template laden
	tmpl, err := template.New("base.html").
		Funcs(template.FuncMap{
			"safeHTML": func(s template.HTML) template.HTML { return s },
		}).
		ParseFiles(filepath.Join(*templateDirPath, "base.html"))
	if err != nil {
		logger.Error("Template New", err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		http.Redirect(w, r, *internalServerErrorURI, http.StatusPermanentRedirect)

		return
	}

	// rendern
	//nolint:gosec // markdown to HTML, no big deal
	var data = struct {
		Title   string
		Content template.HTML
	}{
		Title:   "hbsdsrv – " + r.URL.Path,
		Content: template.HTML(htmlContent),
	}

	if err := tmpl.Execute(w, data); err != nil {
		logger.Error("Template Execution", err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		http.Redirect(w, r, *internalServerErrorURI, http.StatusPermanentRedirect)
	}
}

func setupMux() (*http.ServeMux, error) {
	var (
		mux   = http.NewServeMux()
		files []string
	)

	// Scan ContentDir
	if err := filepath.Walk(
		*contentDirPath,

		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("got err from filepath.Walk: %s", err.Error())
			}

			if !info.IsDir() {
				relPath, err := filepath.Rel(*contentDirPath, path)
				if err != nil {
					return fmt.Errorf("filepath.Rel path err: %s", err.Error())
				}

				relPath = filepath.ToSlash(relPath) // \ → /
				uriPath := "/" + relPath
				files = append(files, uriPath)

				uris[uriPath] = contentHandler
			}

			return nil
		},
	); err != nil {
		return nil, err
	}

	genIndex(files)

	// register
	for p, f := range uris {
		mux.HandleFunc(p, f)
	}

	return mux, nil
}
*/
