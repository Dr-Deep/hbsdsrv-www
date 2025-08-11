// publish markdown files via HTTP3
package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Dr-Deep/logging-go"

	http3 "github.com/quic-go/quic-go/http3"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

/*
socket?
markdown dir?
template dir?

dieser teil bekommt von nginx /www



// uri prefix support/var
*/

var (
	logger *logging.Logger

	// Global Vars
	uris = map[string]func(http.ResponseWriter, *http.Request){}

	// Generell Flags
	socket = flag.String(
		"socket",
		"localhost:8082",
		"socket to listen on",
	)

	certFilePath = flag.String(
		"cert",
		"./fullchain.pem", // /usr/local/etc/letsencrypt/live/vmd171781.contaboserver.net/fullchain.pem
		"ssl certificate",
	)
	keyFilePath = flag.String(
		"key",
		"./privkey.pem", // /usr/local/etc/letsencrypt/live/vmd171781.contaboserver.net/privkey.pem
		"ssl key",
	)

	// Dirs
	templateDirPath = flag.String(
		"templates",
		"./templates", // /usr/local/www
		"HTML templates directory",
	)
	contentDirPath = flag.String(
		"content",
		"./content", // /usr/local/www/www
		"Markdown directory",
	)

	// Default Sites
	indexHTMLSitePath = flag.String(
		"index",
		*contentDirPath+"/index.md",
		"where to put the generated index.md",
	)
	notFoundHTMLURI = flag.String(
		"404",
		"/404.html",
		"404.html URI",
	)
	internalServerErrorURI = flag.String(
		"50x",
		"/50x.html",
		"50x.html URI",
	)
)

func init() {
	flag.Parse()
	setupLogger()
}

func setupLogger() {
	logger = logging.NewLogger(os.Stdout)
}

func main() {
	mux, err := setupMux()
	if err != nil {
		panic(err)
	}

	logger.Info("listening on",
		fmt.Sprintf("https://%s/", *socket),
	)
	mux.HandleFunc("/", indexHandler)

	if err := http3.ListenAndServeTLS(
		*socket,
		*certFilePath,
		*keyFilePath,
		mux,
	); err != nil {
		panic(err)
	}

	logger.Info("Quitting...")
}

func genIndex(files []string) {
	var (
		indexFilePath = *contentDirPath + "/index.md"
		index         = strings.Builder{}
	)
	logger.Debug("generating index file", *indexHTMLSitePath)

	index.WriteString("# Index.md - Overview\n")
	for _, file := range files {
		index.WriteString(
			fmt.Sprintf("* [%s](%s)\n", file, file),
		)
		logger.Debug(file) // prefix?
	}

	if err := os.WriteFile(indexFilePath, []byte(index.String()), 0444); err != nil { // read for everyone, write for none
		panic(err)
	}
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
	var path = *contentDirPath + r.URL.Path
	content, err := os.ReadFile(path)
	if err != nil {
		logger.Error("OS ReadFile from", *contentDirPath+r.URL.Path, err.Error())
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
	var mux = http.NewServeMux()

	// Scan ContentDir
	var files []string
	if err := filepath.Walk(
		*contentDirPath,

		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				relPath, err := filepath.Rel(*contentDirPath, path)
				if err != nil {
					return err
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

func renderMarkdownToHTML(md []byte) string {
	var (
		p   = parser.New()
		doc = p.Parse(md)
		r   = html.NewRenderer(
			html.RendererOptions{
				Flags: html.CommonFlags,
			},
		)
	)

	return string(markdown.Render(doc, r))
}
