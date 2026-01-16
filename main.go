package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"json_formatter/handlers"
)

// PageData holds values passed into the template.
type PageData struct {
	RawInput   string
	Key        string
	Value      string
	Formatted  string
	Matches    []string
	KeyMatches []string
	Error      string
}

// App encapsulates dependencies for HTTP handlers.
type App struct {
	tmpl *template.Template
}

//go:embed templates/*.html
var tmplFS embed.FS

//go:embed static/*
var staticFS embed.FS

func main() {
	tmpl := template.Must(template.ParseFS(tmplFS, "templates/layout.html", "templates/index.html"))

	app := &App{tmpl: tmpl}

	r := gin.Default()
	staticContent, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatalf("static FS error: %v", err)
	}
	r.GET("/static/*filepath", func(c *gin.Context) {
		file := strings.TrimPrefix(c.Param("filepath"), "/")
		if file == "" {
			c.Status(http.StatusNotFound)
			return
		}
		c.FileFromFS(file, http.FS(staticContent))
	})

	r.GET("/", app.handleIndex)
	r.NoRoute(app.handleIndex)

	r.POST("/format", app.handleFormat)
	r.POST("/find/key", app.handleFindKey)
	r.POST("/find/value", app.handleFindValue)
	r.POST("/minify", app.handleMinify)
	r.POST("/toyaml", app.handleToYAML)
	r.POST("/extract/key", app.handleExtractKey)

	r.GET("/healthz", app.handleHealth)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}
	if port[0] != ':' {
		port = ":" + port
	}
	log.Println(fmt.Sprintf("listening on http://localhost%s", port))
	if err := r.Run(port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func (a *App) renderPage(c *gin.Context, data PageData, status int) {
	c.Status(status)
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := a.tmpl.ExecuteTemplate(c.Writer, "base", data); err != nil {
		c.String(http.StatusInternalServerError, "template error")
	}
}

func (a *App) handleIndex(c *gin.Context) {
	a.renderPage(c, PageData{}, http.StatusOK)
}

func (a *App) handleFormat(c *gin.Context) {
	raw := c.PostForm("jsonInput")
	key := c.PostForm("key")
	value := c.PostForm("value")

	formatted, matches, keyMatches, err := handlers.Process(raw, key, value)
	data := PageData{
		RawInput:   raw,
		Key:        key,
		Value:      value,
		Formatted:  formatted,
		Matches:    matches,
		KeyMatches: keyMatches,
	}
	if err != nil {
		data.Error = err.Error()
		a.renderPage(c, data, http.StatusBadRequest)
		return
	}
	a.renderPage(c, data, http.StatusOK)
}

func (a *App) handleFindKey(c *gin.Context) {
	raw := c.PostForm("jsonInput")
	key := c.PostForm("key")
	formatted, matches, _, err := handlers.Process(raw, key, "")
	data := PageData{RawInput: raw, Key: key, Formatted: formatted, Matches: matches}
	if err != nil {
		data.Error = err.Error()
		a.renderPage(c, data, http.StatusBadRequest)
		return
	}
	a.renderPage(c, data, http.StatusOK)
}

func (a *App) handleFindValue(c *gin.Context) {
	raw := c.PostForm("jsonInput")
	value := c.PostForm("value")
	formatted, _, keyMatches, err := handlers.Process(raw, "", value)
	data := PageData{RawInput: raw, Value: value, Formatted: formatted, KeyMatches: keyMatches}
	if err != nil {
		data.Error = err.Error()
		a.renderPage(c, data, http.StatusBadRequest)
		return
	}
	a.renderPage(c, data, http.StatusOK)
}

func (a *App) handleMinify(c *gin.Context) {
	raw := c.PostForm("jsonInput")
	minified, err := handlers.Minify(raw)
	data := PageData{RawInput: raw, Formatted: minified}
	if err != nil {
		data.Error = err.Error()
		a.renderPage(c, data, http.StatusBadRequest)
		return
	}
	a.renderPage(c, data, http.StatusOK)
}

func (a *App) handleToYAML(c *gin.Context) {
	raw := c.PostForm("jsonInput")
	out, err := handlers.ToYAML(raw)
	data := PageData{RawInput: raw, Formatted: out}
	if err != nil {
		data.Error = err.Error()
		a.renderPage(c, data, http.StatusBadRequest)
		return
	}
	a.renderPage(c, data, http.StatusOK)
}

func (a *App) handleExtractKey(c *gin.Context) {
	raw := c.PostForm("jsonInput")
	key := c.PostForm("key")
	out, err := handlers.ExtractKeyJSON(raw, key)
	data := PageData{RawInput: raw, Key: key, Formatted: out}
	if err != nil {
		data.Error = err.Error()
		a.renderPage(c, data, http.StatusBadRequest)
		return
	}
	a.renderPage(c, data, http.StatusOK)
}

func (a *App) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
