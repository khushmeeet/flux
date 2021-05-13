package fluxgen

import (
	"html/template"
	"net/http"
	"time"
)

type Page struct {
	Title        string
	Date         time.Time
	template     string
	Href         string
	oldExtention string
	newExtension string
	filename     string
	Content      template.HTML
	Meta         map[string]interface{}
	PostsList    *Pages
	resources    *Resources
	fluxConfig   *FluxConfig
}

type Pages []Page

type FluxConfig map[string]interface{}

type Resources map[string]string

type responseObserver struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}
