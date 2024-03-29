package fluxgen

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Joker/hpp"
	_ "github.com/flosch/pongo2-addons"
	"github.com/flosch/pongo2/v4"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/muesli/termenv"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

var green = termenv.ColorProfile().Color("#29bc89")

func (p *page) applyTemplate() (string, error) {
	tmpl, err := pongo2.FromFile(p.template + ".html")
	if err != nil {
		return "", err
	}

	ctx := pongo2.Context{
		"Title":   p.Title,
		"Date":    p.Date,
		"Content": p.Content,
		"Meta":    p.Meta,
		"Posts":   p.PostsList,
		"Flux":    p.fluxConfig,
		"GetPath": getResource(p.resources),
	}

	out, err := tmpl.Execute(ctx)
	if err != nil {
		return "", err
	}
	return hpp.PrPrint(out), nil
}

func parseMarkdown(path string, config *fluxConfig, r *resources) (page, error) {
	var buff bytes.Buffer
	context := parser.NewContext()
	md := goldmark.New(
		goldmark.WithExtensions(meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
			),
			extension.Footnote,
			extension.GFM,
			emoji.Emoji,
			mathjax.MathJax,
		),
		goldmark.WithParserOptions(
			parser.WithBlockParsers(),
			parser.WithAttribute(),
		),
	)

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return page{}, nil
	}

	err = md.Convert(file, &buff, parser.WithContext(context))
	if err != nil {
		return page{}, err
	}
	frontMatter := meta.Get(context)

	var parsedDate time.Time
	if val, ok := frontMatter["data"]; ok {
		parsedDate, err = time.Parse("2006-01-02", val.(string))
		if err != nil {
			return page{}, err
		}
	}

	metaData := make(map[string]interface{})
	for k, v := range frontMatter {
		if k != "title" && k != "template" {
			metaData[k] = v
		}
	}

	page := page{
		Title:      frontMatter["title"].(string),
		Date:       parsedDate,
		template:   filepath.Join(TemplatesDir, frontMatter["template"].(string)),
		filename:   filepath.Base(path),
		Content:    template.HTML(buff.Bytes()),
		Meta:       make(map[string]interface{}),
		fluxConfig: config,
		resources:  r,
	}
	page.setHref(path)
	return page, nil
}

func parseHTML(path string, config *fluxConfig, resources *resources) (page, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return page{}, err
	}

	page := page{
		Title:      "",
		Date:       time.Time{},
		template:   strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		filename:   filepath.Base(path),
		Content:    template.HTML(file),
		Meta:       make(map[string]interface{}),
		fluxConfig: config,
		resources:  resources,
	}
	page.setHref(path)
	return page, nil
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

func (p *page) setHref(path string) {
	if filepath.Base(path) == "index.html" {
		p.Href = "/"
	} else if filepath.Ext(path) == ".html" {
		p.Href = filepath.Join("/", strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)), "/")
	} else {
		p.Href = filepath.Join("/", strings.TrimSuffix(path, filepath.Base(path)))
	}
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

func createFileWriteDir(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.MkdirAll(filePath, 0744)
		if err != nil {
			return err
		}
	}
	return nil
}

func getResource(r *resources) func(v string) string {
	return func(v string) string {
		return (*r)[v]
	}
}

func printMsg(msg, status string) {
	var info string
	if status == "tick" {
		info = termenv.String("✔").Foreground(green).String()
	} else if status == "party" {
		info = termenv.String("️🎉").String()
	}
	fmt.Printf("%s %s\n", info, msg)
}
