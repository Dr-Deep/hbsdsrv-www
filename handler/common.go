package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

const (
	indexHTML = "./html/index.html"
	loginHTML = "./html/login.html"
	errorHTML = "./html/error.html"
	baseHTML  = "./html/base.html"

	contentURL     = "/www"
	contentMdIndex = "./www/index.md"
	assetsURL      = "/assets"
)

/*
HTMLDirectory   string   `yaml:"html-dir"`
WWWDirectory    string   `yaml:"www-dir"`
AssetsDirectory string   `yaml:"assets-dir"`
*/

// returns rendered html, error
func renderHTMLTemplate(templateFilePath string, templateData any) (string, error) {
	t, err := template.New(filepath.Base(templateFilePath)).
		Funcs(template.FuncMap{
			"safeHTML": func(s template.HTML) template.HTML { return s },
		}).
		ParseFiles(templateFilePath)
	if err != nil {
		return "", err
	}

	var buf = &bytes.Buffer{}
	if err := t.ExecuteTemplate(buf, filepath.Base(templateFilePath), templateData); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func renderMarkdownToHTML(_markdown []byte) string {
	var (
		p   = parser.New()
		doc = p.Parse(_markdown)
		r   = html.NewRenderer(
			html.RendererOptions{
				Flags: html.CommonFlags,
			},
		)
	)

	return string(markdown.Render(doc, r))
}

// returns map[urlPath]fsPath
func genFsMap(directory, urlPrefix string) (map[string]string, error) {
	var pathMap = map[string]string{}

	// walk & store
	var walkFunc = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walkFunc lastErr: %s: %s", path, err.Error())
		}

		if !info.IsDir() {
			// relative path
			relativePath, err := filepath.Rel(directory, path)
			if err != nil {
				return fmt.Errorf("walkFunc curErr: %s: %s", path, err.Error())
			}

			uriPath := urlPrefix + "/" + relativePath
			pathMap[uriPath] = path
		}

		return nil
	}

	if err := filepath.Walk(
		directory,
		walkFunc,
	); err != nil {
		return nil, err
	}

	return pathMap, nil
}
