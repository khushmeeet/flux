package fluxgen

import (
	"bytes"
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
	"strings"
)

type Page struct {
	Href     string
	Name     string
	Content  template.HTML
	MetaData map[string]interface{}
}

func Generate() {
	tmplMap := parseTemplatesWithPartials()
	err := filepath.Walk(PagesFolder, func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".md" {
			markdownFile, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatal("Unable to read Markdown file!")
			}
			page, err := parsePages(path, markdownFile)
			if err != nil {
				log.Fatal("Unable to convert Markdown to HTML!")
			}
			executeTemplates(tmplMap, page)
		}
		return err
	})
	if err != nil {
		log.Fatal("Not able to scan Pages folder!")
	}
}

func parsePages(path string, src []byte) (Page, error) {
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
		return Page{}, err
	}
	metaData := meta.Get(context)
	pageData := Page{
		Href:     "/" + filepath.Base(path),
		Name:     filepath.Base(path),
		Content:  template.HTML(buff.Bytes()),
		MetaData: metaData,
	}
	return pageData, nil
}

func executeTemplates(tmplMap map[string]*template.Template, page Page) {
	fileName := strings.Split(page.Name, ".")[0]
	formattedFileName := strings.Join(strings.Split(fileName, " "), "-")
	file, err := os.Create(path.Join(SiteFolder, formattedFileName+".html"))
	err = tmplMap[page.MetaData["template"].(string)].Execute(file, page)
	if err != nil {
		log.Fatal("Unable to execute templates!", err)
	}
}

func parseTemplatesWithPartials() map[string]*template.Template {
	parsedTmplMap := make(map[string]*template.Template)
	templatesInfo, err := ioutil.ReadDir(TemplatesFolder)
	if err != nil {
		log.Fatal("Error reading Templates Folder!")
	}
	for _, i := range templatesInfo {
		if !i.IsDir() {
			tmpl := template.New(i.Name())
			var t *template.Template
			err = filepath.Walk(TemplatesFolder, func(path string, info fs.FileInfo, err error) error {
				if filepath.Ext(path) == ".html" {
					t, err = tmpl.ParseFiles(path)
					if err != nil {
						return err
					}
				}
				return nil
			})
			if err != nil {
				log.Fatal("Not able to walk Templates folder!")
			}
			parsedTmplMap[i.Name()] = t
		}
	}
	return parsedTmplMap
}
