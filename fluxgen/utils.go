package fluxgen

import (
	"bytes"
	"fmt"
	"github.com/flosch/pongo2/v4"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (p *Page) setHref(path string) {
	if filepath.Base(path) == "index.html" {
		p.Href = "/"
	} else if filepath.Ext(path) == ".html" {
		p.Href = filepath.Join(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)), "/")
	} else {
		p.Href = strings.TrimSuffix(path, filepath.Base(path))
	}
}

func (p *Page) applyTemplate() (string, error) {
	tmpl, err := pongo2.FromFile(p.Template + ".html")
	if err != nil {
		panic(err)
	}

	fmt.Println(p.PostList)

	ctx := pongo2.Context{
		"title":     p.Title,
		"date":      p.Date,
		"content":   p.Content,
		"meta":      p.MetaData,
		"posts":     *p.PostList,
		"flux":      p.FluxConfig,
		"resources": p.Resources,
	}

	out, err := tmpl.Execute(ctx)
	if err != nil {
		return "", err
	}
	return out, nil
}

func copyFile(src, dst string) error {
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

func createFileWritePath(fileName string, filePath string) string {
	fileWritePath := ""
	if fileName == "index.html" {
		fileWritePath = SiteDir
	} else {
		fileWritePath = filepath.Join(SiteDir, filePath)
	}
	return fileWritePath
}

func createFileWriteDir(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		_ = os.MkdirAll(filePath, 0744)
	} else {
		fmt.Println("Directory already exists")
	}
}

func parseMarkdown(path string, config *FluxConfig, r *Resources) Page {
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

	parsedDate, err := time.Parse("2006-01-02", frontMatter["date"].(string))
	if err != nil {
		log.Fatalf("[Error Parsing Time (%v)] - %v", frontMatter["date"].(string), err)
	}

	metaData := make(map[string]interface{})
	for k, v := range frontMatter {
		if k != "title" && k != "date" && k != "template" {
			metaData[k] = v
		}
	}

	page := Page{
		Title:        frontMatter["title"].(string),
		Date:         parsedDate,
		Template:     filepath.Join(TemplatesDir, frontMatter["template"].(string)),
		OldExtension: filepath.Ext(path),
		NewExtension: ".html",
		FileName:     filepath.Base(path),
		Content:      template.HTML(buff.Bytes()),
		MetaData:     make(map[string]interface{}),
		FluxConfig:   config,
		Resources:    r,
	}
	page.setHref(path)
	return page
}

func parseHTML(path string, config *FluxConfig, resources *Resources) Page {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("[Error Reading (%v)] - %v", path, err)
	}

	page := Page{
		Title:        "",
		Date:         time.Time{},
		Template:     strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		OldExtension: filepath.Ext(path),
		NewExtension: ".html",
		FileName:     filepath.Base(path),
		Content:      template.HTML(file),
		MetaData:     make(map[string]interface{}),
		FluxConfig:   config,
		Resources:    resources,
	}
	page.setHref(path)
	return page
}
