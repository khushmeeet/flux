package fluxgen

import (
	"bytes"
	"encoding/json"
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
	"path/filepath"
	"strings"
	"time"
)

type Page struct {
	Title      string
	Date       time.Time
	Template   string
	Href       string
	Extension  string
	Content    template.HTML
	MetaData   map[string]interface{}
	AllPages   *Pages
	FluxConfig *FluxConfig
}

type Pages []Page

type FluxConfig map[string]interface{}

func FluxBuild() {
	fluxConfig := parseFluxConfig(ConfigFile)
	pageList := parsePages(PagesFolder, &fluxConfig)

	By(descendingOrderByDate).Sort(pageList)

	parseHTMLTemplates(TemplatesFolder, pageList)
}

func parseFluxConfig(path string) FluxConfig {
	fluxConfig := make(FluxConfig)
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("[Error Reading (%v)] - %v", path, err)
	}
	err = json.Unmarshal(configFile, &fluxConfig)
	if err != nil {
		log.Fatalf("[Error Unmarshalling (%v)] - %v", path, err)
	}
	return fluxConfig
}

func parseHTMLTemplates(path string, pages Pages) {
	funcMap := template.FuncMap{}
	tmpl, err := template.New("index").Funcs(funcMap).ParseGlob(path + "/*.html")
	if err != nil {
		log.Fatalf("[Error Parsing Template Dir (%v)] - %v", path, err)
	}

	for _, p := range pages {
		p.AllPages = &pages
		buffer, err := p.applyTemplate(tmpl)
		if err != nil {
			log.Fatalf("[Error Applying Template to Page] - %v", err)
		}

		err = ioutil.WriteFile(filepath.Join(SiteFolder, p.Href+p.Extension), buffer.Bytes(), 07444)
		if err != nil {
			log.Fatalf("[Error Writing File (%v)] - %v", p.Href, err)
		}
		fmt.Printf("Writing File: %v\n", p.Href)
	}
}

func (p *Page) applyTemplate(t *template.Template) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	templateFile := p.Template + ".html"
	err := t.ExecuteTemplate(buffer, templateFile, p)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func parsePages(filePath string, config *FluxConfig) Pages {
	var pageList Pages
	err := filepath.Walk(filePath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".md" {
			page := parseMarkdown(path, config)
			pageList = append(pageList, page)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("[Error Walking (%v)] - %v", filePath, err)
	}
	return pageList
}

func parseMarkdown(path string, config *FluxConfig) Page {
	var buff bytes.Buffer
	context := parser.NewContext()
	md := goldmark.New(
		goldmark.WithExtensions(meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"))),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)

	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("[Error Reading (%v)] - %v", path, err)
	}

	err = md.Convert(file, &buff, parser.WithContext(context))
	if err != nil {
		log.Fatalf("[Error Parsing Markdown File (%v)] - %v", path, err)
	}
	frontMatter := meta.Get(context)

	date, err := time.Parse("2006-01-02", frontMatter["date"].(string))
	if err != nil {
		log.Fatalf("[Error Parsing Time (%v)] - %v", frontMatter["date"].(string), err)
	}

	page := Page{
		Title:      frontMatter["title"].(string),
		Date:       date,
		Template:   frontMatter["template"].(string),
		Href:       strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		Extension:  ".html",
		Content:    template.HTML(buff.Bytes()),
		MetaData:   make(map[string]interface{}),
		FluxConfig: config,
	}

	for k, v := range frontMatter {
		if k != "title" && k != "date" && k != "template" {
			page.MetaData[k] = v
		}
	}

	return page
}
