package fluxgen

import (
	"html/template"
	"net/http"
	"time"
)

type page struct {
	Title      string
	Date       time.Time
	template   string
	Href       string
	filename   string
	Content    template.HTML
	Meta       map[string]interface{}
	PostsList  *pages
	resources  *resources
	fluxConfig *fluxConfig
}

type pages []page

type fluxConfig map[string]interface{}

type resources map[string]string

type responseObserver struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}
