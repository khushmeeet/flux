package fluxgen

import (
	"html/template"
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
	PostList     *Pages
	Resources    *Resources
	FluxConfig   *FluxConfig
}

type Pages []Page

type FluxConfig map[string]interface{}

type Resources map[string]string
