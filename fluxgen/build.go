package fluxgen

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/sprig/v3"
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
	Resources    *Resources
	FluxConfig   *FluxConfig
}

type Pages []Page

type FluxConfig map[string]interface{}

type Resources map[string]string

func FluxBuild() {
	FluxClean()
	fluxConfig := parseFluxConfig(ConfigFile)
	resources := loadResources(CSSDir)
	pageList, postList := parsePages(&fluxConfig, &resources)
	By(descendingOrderByDate).Sort(pageList)
	By(descendingOrderByDate).Sort(postList)
	parseHTMLTemplates(TemplatesDir, pageList, postList)
	processStaticFolders(CSSDir)
	processStaticFolders(AssetsDir)
}

func loadResources(path ...string) Resources {
	resources := Resources{}
	for _, i := range path {
		err := filepath.WalkDir(i, func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() {
				resources[strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))] = path
			}
			return nil
		})
		if err != nil {
			log.Fatalf("Error walking (%v) - %v", i, err)
		}
	}
	return resources
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

func parsePages(config *FluxConfig, resources *Resources) (Pages, Pages) {
	var pageList Pages
	var postList Pages
	err := filepath.Walk(PagesDir, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".md" {
			mdPage := parseMarkdown(path, config, resources)
			pageList = append(pageList, mdPage)
			postList = append(postList, mdPage)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("[Error Walking (%v)] - %v", PagesDir, err)
	}

	dirContent, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatalf("[Error Reading (%v)] - %v", ".", err)
	}

	for _, f := range dirContent {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".html" {
			htmlPage := parseHTML(f.Name(), config, resources)
			pageList = append(pageList, htmlPage)
		}
	}

	return pageList, postList
}

func parseHTMLTemplates(path string, pages Pages, posts Pages) {
	templateList := getHTMLTemplatesList()
	tmpl, err := template.New("default").Funcs(sprig.FuncMap()).ParseFiles(templateList...)
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
				err := os.MkdirAll(filepath.Join(SiteDir, path), 0744)
				if err != nil {
					return err
				}
				fmt.Printf("Creating Folder: %v\n", path)
			} else {
				err := copyFile(path, filepath.Join(SiteDir, path))
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
