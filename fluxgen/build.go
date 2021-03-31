package fluxgen

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type Page struct {
	Href       string
	Name       string
	Content    template.HTML
	MetaData   map[string]interface{}
	FluxConfig FluxConfig
}

type FluxConfig map[string]interface{}

func FluxBuild() {
	fc := readConfigFile()
	tmplMap := parseTemplatesWithPartials()

	pageSlice, err := parsePages(fc)
	if err != nil {
		log.Fatal("Unable to parse Markdown files!")
	}

	for _, i := range pageSlice {
		executeTemplates(tmplMap, i)
	}
}

func parsePages(fc FluxConfig) ([]Page, error) {
	var buff bytes.Buffer
	var pageSlice []Page
	context := parser.NewContext()
	md := goldmark.New(
		goldmark.WithExtensions(meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"))),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)

	err := filepath.Walk(PagesFolder, func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".md" {
			fileName := strings.Split(filepath.Base(path), ".")[0]
			formattedFileName := strings.Join(strings.Split(fileName, " "), "-")

			markdownFile, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			if err := md.Convert(markdownFile, &buff, parser.WithContext(context)); err != nil {
				return err
			}
			metaData := meta.Get(context)
			pageData := Page{
				Href:       "/" + formattedFileName,
				Name:       formattedFileName,
				Content:    template.HTML(buff.Bytes()),
				MetaData:   metaData,
				FluxConfig: fc,
			}
			pageSlice = append(pageSlice, pageData)

		}
		return nil
	})
	if err != nil {
		log.Fatal("Not able to scan Pages folder!")
	}

	return pageSlice, nil
}

func executeTemplates(tmplMap map[string]*template.Template, page Page) {
	file, err := os.Create(path.Join(SiteFolder, page.Name+".html"))
	if err != nil {
		log.Fatal("Unable to create HTML file!")
	}

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

func readConfigFile() FluxConfig {
	var configMap FluxConfig
	configFile, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		log.Fatal("Unable to read Config File!")
	}

	err = json.Unmarshal(configFile, &configMap)
	if err != nil {
		log.Fatal("Unable to parse Config File!")
	}

	return configMap
}
