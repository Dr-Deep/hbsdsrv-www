package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
*/

var (
	// Global Vars
	uris = map[string]func(http.ResponseWriter, *http.Request){}

	// Flags
	socket = *flag.String(
		"socket",
		"localhost:8082",
		"socket to listen on",
	)
	certFilePath = *flag.String(
		"cert",
		"./fullchain.pem", //usr/local/etc/letsencrypt/live/vmd171781.contaboserver.net/fullchain.pem
		"ssl certificate",
	)
	keyFilePath = *flag.String(
		"key",
		"./privkey.pem", //usr/local/etc/letsencrypt/live/vmd171781.contaboserver.net/privkey.pem
		"ssl key",
	)
	templateDirPath = *flag.String(
		"templates",
		"./templates", //usr/local/www
		"HTML templates directory",
	)
	contentDirPath = *flag.String(
		"content",
		"./content", //usr/local/www/www
		"Markdown directory",
	)
)

func init() {
	flag.Parse()

}

func main() {
	var mux = setupMux()

	fmt.Printf("listening on https://%s/", socket)
	mux.HandleFunc("/", indexHandler)

	if err := http3.ListenAndServeTLS(
		socket,
		certFilePath,
		keyFilePath,
		mux,
	); err != nil {
		panic(err)
	}

	fmt.Println("bye")
}

func genIndex(files []string) {
	var index = strings.Builder{}
	index.WriteString("# Index.md - Overview\n")
	for _, file := range files {
		index.WriteString(
			fmt.Sprintf("* [%s](%s)\n", file, file),
		)
	}

	if err := os.WriteFile(contentDirPath+"/index.md", []byte(index.String()), os.ModePerm); err != nil {
		panic(err)
	}
}

// func(w http.ResponseWriter, r *http.Request)
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// valid Path?

	switch r.URL.Path {
	case "/":
		// redir to index
		http.Redirect(w, r, "/index.md", http.StatusPermanentRedirect)

	default:
		// Immer zuerst Status setzen, dann schreiben!
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 – Seite nicht gefunden"))
	}
}

// contentHandler()
func contentHandler(w http.ResponseWriter, r *http.Request) {

	// need okay from uris
	_, oke := uris[r.URL.Path]
	if !oke {
		http.Redirect(w, r, "/404.html", http.StatusPermanentRedirect)
	}

	// Read File
	var path = contentDirPath + r.URL.Path
	content, err := os.ReadFile(path)
	if err != nil {
		http.Redirect(w, r, "/404.html", http.StatusPermanentRedirect)
	}

	htmlContent := renderMarkdownToHTML(content)

	// Template laden
	tmpl, err := template.New("base.html").
		Funcs(template.FuncMap{
			"safeHTML": func(s template.HTML) template.HTML { return s },
		}).
		ParseFiles(filepath.Join(templateDirPath, "base.html"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 – Fehler beim Template"))
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 – Template rendering error"))
	}
}

// Todo: uri regexp?
func setupMux() (mux *http.ServeMux) {
	mux = http.NewServeMux()

	// Scan ContentDir
	var files []string
	if err := filepath.Walk(
		contentDirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				relPath, err := filepath.Rel(contentDirPath, path)
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
		fmt.Println("content dir Error:", err)
	}

	genIndex(files)

	// register
	for p, f := range uris {
		mux.HandleFunc(p, f)
	}

	return
}

func renderMarkdownToHTML(md []byte) string {

	p := parser.New()
	doc := p.Parse(md)

	r := html.NewRenderer(
		html.RendererOptions{
			Flags: html.CommonFlags,
		},
	)
	return string(markdown.Render(doc, r))
}
