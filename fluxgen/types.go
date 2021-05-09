package fluxgen

import (
	"html/template"
	"time"
)

type Page struct {
	title        string
	date         time.Time
	template     string
	href         string
	oldExtention string
	newExtension string
	filename     string
	content      template.HTML
	meta         map[string]interface{}
	postsList    *Pages
	resources    *Resources
	fluxConfig   *FluxConfig
}

type Pages []Page

type FluxConfig map[string]interface{}

type Resources map[string]string
