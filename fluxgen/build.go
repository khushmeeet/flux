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
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Page struct {
	Title        string
	Date         time.Time
	Template     string
	Href         string
	OldExtension string
	NewExtension string
	Content      template.HTML
	MetaData     map[string]interface{}
	PageList     *Pages
	PostList     *Pages
	FluxConfig   *FluxConfig
}

type Pages []Page

type FluxConfig map[string]interface{}

func FluxBuild() {
	fluxConfig := parseFluxConfig(ConfigFile)
	pageList, postList := parsePages(PagesFolder, &fluxConfig)
	By(descendingOrderByDate).Sort(pageList)
	By(descendingOrderByDate).Sort(postList)
	parseHTMLTemplates(TemplatesFolder, pageList, postList)
	processAssets(PagesFolder)
	processStatic(StaticFolder)
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

func parsePages(filePath string, config *FluxConfig) (Pages, Pages) {
	var pageList Pages
	var postList Pages
	err := filepath.Walk(filePath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && (filepath.Ext(path) == ".md" || filepath.Ext(path) == ".html") {
			page := parseMarkdown(path, config)
			pageList = append(pageList, page)
			if filepath.Ext(path) == ".md" {
				postList = append(postList, page)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("[Error Walking (%v)] - %v", filePath, err)
	}
	return pageList, postList
}

func parseHTMLTemplates(path string, pages Pages, posts Pages) {
	pagesList, _ := filepath.Glob(PagesFolder + "/*.html")
	templatesList, _ := filepath.Glob(TemplatesFolder + "/*.html")
	allTemplatesList := append(pagesList, templatesList...)

	funcMap := template.FuncMap{}
	tmpl, err := template.New("index").Funcs(funcMap).ParseFiles(allTemplatesList...)
	if err != nil {
		log.Fatalf("[Error Parsing Template Dir (%v)] - %v", path, err)
	}

	for _, p := range pages {
		p.PageList = &pages
		p.PostList = &posts
		buffer, err := p.applyTemplate(tmpl)
		if err != nil {
			log.Fatalf("[Error Applying Template to Page] - %v", err)
		}

		err = ioutil.WriteFile(filepath.Join(SiteFolder, p.Href+p.NewExtension), buffer.Bytes(), 07444)
		if err != nil {
			log.Fatalf("[Error Writing File (%v)] - %v", p.Href, err)
		}
		fmt.Printf("Writing File: %v\n", p.Href+p.OldExtension)
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

	var parsedDate time.Time
	date, ok := frontMatter["date"]
	if ok {
		parsedDate, err = time.Parse("2006-01-02", date.(string))
		if err != nil {
			log.Fatalf("[Error Parsing Time (%v)] - %v", date.(string), err)
		}
	}

	var templateFile string
	if getMapValue(frontMatter, "template") == "" {
		templateFile = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	} else {
		templateFile = getMapValue(frontMatter, "template")
	}

	page := Page{
		Title:        getMapValue(frontMatter, "title"),
		Date:         parsedDate,
		Template:     templateFile,
		Href:         strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		OldExtension: filepath.Ext(path),
		NewExtension: ".html",
		Content:      template.HTML(buff.Bytes()),
		MetaData:     make(map[string]interface{}),
		FluxConfig:   config,
	}

	for k, v := range frontMatter {
		if k != "title" && k != "date" && k != "template" {
			page.MetaData[k] = v
		}
	}
	return page
}

func processAssets(filePath string) {
	err := filepath.Walk(filePath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) != ".md" && filepath.Ext(path) != ".html" {
			err := copyFile(path, filepath.Join(SiteFolder, filepath.Base(path)))
			fmt.Printf("Copying File: %v\n", path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("[Error Walking (%v)] - %v", filePath, err)
	}
}

func processStatic(filePath string) {
	err := filepath.Walk(filePath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			err := os.MkdirAll(filepath.Join(SiteFolder, path), 0744)
			if err != nil {
				return err
			}
			fmt.Printf("Creating Folder: %v\n", path)
		} else {
			err := copyFile(path, filepath.Join(SiteFolder, path))
			if err != nil {
				return err
			}
			fmt.Printf("Copying File: %v\n", path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("[Error Walking (%v)] - %v", filePath, err)
	}
}
