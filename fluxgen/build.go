package fluxgen

import (
	"encoding/json"
	"fmt"
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
	FileName     string
	Content      template.HTML
	MetaData     map[string]interface{}
	PageList     *Pages
	PostList     *Pages
	FluxConfig   *FluxConfig
}

type Pages []Page

type FluxConfig map[string]interface{}

func FluxBuild() {
	FluxClean()
	fluxConfig := parseFluxConfig(ConfigFile)
	pageList, postList := parsePages(PagesFolder, &fluxConfig)
	By(descendingOrderByDate).Sort(pageList)
	By(descendingOrderByDate).Sort(postList)
	parseHTMLTemplates(TemplatesFolder, pageList, postList)
	//processAssets(PagesFolder)
	processStaticFolders(CSSFolder)
	processStaticFolders(AssetsFolder)
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

	tmpl, err := template.New("index").Funcs(sprig.FuncMap()).ParseFiles(allTemplatesList...)
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

		fileWritePath := createFileWritePath(p.FileName, p.Href)
		createFileWriteDir(fileWritePath)

		err = ioutil.WriteFile(filepath.Join(fileWritePath, "index.html"), buffer.Bytes(), 0744)
		if err != nil {
			log.Fatalf("[Error Writing File (%v)] - %v", p.Href, err)
		}
		fmt.Printf("Writing File: %v\n", p.Href+p.OldExtension)
	}
}

func processStaticFolders(filePath string) {
	if _, err := os.Stat(filePath); err == nil {
		err = filepath.Walk(filePath, func(path string, info fs.FileInfo, err error) error {
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
	} else {
		fmt.Printf("\"%v\" does not exists... Skipping", filePath)
	}
}
