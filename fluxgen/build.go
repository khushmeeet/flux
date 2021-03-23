package fluxgen

import (
	"bytes"
	"fmt"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

type Page struct {
	Title    string
	Short    string
	Template string
	Date     time.Time
	Html     []byte
	Tags     []string
	Images   []string
}

func markdownToHtml(src []byte) ([]byte, error) {
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
		return []byte{}, err
	}

	return buff.Bytes(), nil
}

func Generate() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to retrieve current working directory")
	}

	err = filepath.Walk(path.Join(currentDir, PagesFolder), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if filepath.Ext(path) == ".md" {
			markdownFile, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatal("Unable to read Markdown file!")
			}
			Html, err := markdownToHtml(markdownFile)
			if err != nil {
				log.Fatal("Unable to convert Markdown to HTML!")
			}
			fmt.Println(string(Html))

		}
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
}
