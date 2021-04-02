package fluxgen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
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

type Pages struct {
	Pages []Page
}

type FluxConfig map[string]interface{}

func FluxBuild() {
	err := parseStaticFiles()
	if err != nil {
		log.Fatal("Unable to parse static files!", err)
	}

	fc := readConfigFile()
	tmplMap := parseTemplatesWithPartials()

	pageSlice, err := parsePages(fc)
	if err != nil {
		log.Fatal("Unable to parse Markdown files!")
	}

	for _, i := range pageSlice {
		executeTemplates(tmplMap, i)
	}

	err = parseGeneralTemplates(pageSlice)
	if err != nil {
		errF := fmt.Errorf("error from parseGeneralTemplates() : %v", err)
		log.Fatal(errF.Error())
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

func parseGeneralTemplates(pages []Page) error {
	var partials []string
	partialsFiles, err := ioutil.ReadDir(path.Join(TemplatesFolder, "partials"))
	for _, f := range partialsFiles {
		partials = append(partials, filepath.Join(TemplatesFolder, "partials", f.Name()))
	}

	dir, err := os.ReadDir(".")
	if err != nil {
		return err
	}

	for _, file := range dir {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".html" {
			partials = append(partials, file.Name())
			tmpl, err := template.New(file.Name()).ParseFiles(partials...)
			if err != nil {
				return err
			}
			f, err := os.Create(path.Join(SiteFolder, file.Name()))
			if err != nil {
				return err
			}
			err = tmpl.Execute(f, Pages{Pages: pages})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func parseStaticFiles() error {
	err := filepath.Walk(StaticFolder, func(path string, info fs.FileInfo, err error) error {
		fmt.Println(filepath.Join(SiteFolder, path))
		if info.IsDir() {
			fmt.Println(filepath.Join(StaticFolder, path))
			err := os.Mkdir(filepath.Join(SiteFolder, path), 0777)
			if err != nil {
				return err
			}
		}
		if !info.IsDir() {
			err = copy(path, filepath.Join(SiteFolder, path))
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func copy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}
