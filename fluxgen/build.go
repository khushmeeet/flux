package fluxgen

import (
	"bytes"
	"fmt"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

func Generate() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to retrieve current working directory")
	}

	err = filepath.Walk(path.Join(currentDir, PagesFolder), func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".md" {
			markdownFile, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatal("Unable to read Markdown file!")
			}
			convertedHtml, metaData, err := markdownToHtml(markdownFile)
			if err != nil {
				log.Fatal("Unable to convert Markdown to HTML!")
			}
			parseHTMLTemplate(convertedHtml, metaData)
		}
		return err
	})
	if err != nil {
		log.Fatal("Not able to scan Pages folder!")
	}
}

func markdownToHtml(src []byte) ([]byte, map[string]interface{}, error) {
	var buff bytes.Buffer

	md := goldmark.New(
		goldmark.WithExtensions(meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"))),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)

	context := parser.NewContext()
	if err := md.Convert(src, &buff, parser.WithContext(context)); err != nil {
		return []byte{}, nil, err
	}
	metaData := meta.Get(context)
	return buff.Bytes(), metaData, nil
}

func parseHTMLTemplate(convertedHtml []byte, metaData map[string]interface{}) {
	pageData := make(map[string]interface{})
	pageData["content"] = template.HTML(convertedHtml)
	fmt.Println(string(convertedHtml))
	for k, v := range metaData {
		pageData[k] = v
	}

	tmpl := parsePartials()
	file, err := os.Create(path.Join(SiteFolder, pageData["template"].(string)))
	err = tmpl.Execute(file, pageData)
	if err != nil {
		log.Fatal("Could not parse HTML Template!", err)
	}
}

func parsePartials() *template.Template {
	tmpl := template.New("index.html")
	var t *template.Template
	err := filepath.Walk(TemplatesFolder, func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".html" {
			t, err = tmpl.ParseFiles(path)
			if err != nil {
				return err
			}
			return nil
		}
		return nil
	})
	if err != nil {
		log.Fatal("Not able to scan Templates folder!")
	}
	return template.Must(t, err)
}
