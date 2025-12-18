package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

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

func render(tmpl *template.Template, c *gin.Context, rawJSON, key, value string) {
	formatted, matches, keyMatches, err := handlers.Process(rawJSON, key, value)
	data := PageData{
		RawInput:   rawJSON,
		Key:        key,
		Value:      value,
		Formatted:  formatted,
		Matches:    matches,
		KeyMatches: keyMatches,
	}
	if err != nil {
		data.Error = err.Error()
	}
	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, "template error")
	}
}

func main() {
	tmplPath := filepath.Join("templates", "index.html")
	tmpl := template.Must(template.ParseFiles(tmplPath))

	r := gin.Default()
	r.Static("/static", "static")

	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
		c.Header("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(c.Writer, PageData{}); err != nil {
			c.String(http.StatusInternalServerError, "template error")
		}
	})

	r.NoRoute(func(c *gin.Context) {
		c.Status(http.StatusOK)
		c.Header("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(c.Writer, PageData{}); err != nil {
			c.String(http.StatusInternalServerError, "template error")
		}
	})

	r.POST("/format", func(c *gin.Context) {
		render(tmpl, c, c.PostForm("jsonInput"), c.PostForm("key"), c.PostForm("value"))
	})

	r.POST("/find/key", func(c *gin.Context) {
		render(tmpl, c, c.PostForm("jsonInput"), c.PostForm("key"), "")
	})

	r.POST("/find/value", func(c *gin.Context) {
		render(tmpl, c, c.PostForm("jsonInput"), "", c.PostForm("value"))
	})

	r.POST("/minify", func(c *gin.Context) {
		rawJSON := c.PostForm("jsonInput")
		minified, err := handlers.Minify(rawJSON)
		data := PageData{RawInput: rawJSON, Formatted: minified}
		if err != nil {
			data.Error = err.Error()
		}
		c.Status(http.StatusOK)
		c.Header("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(c.Writer, data); err != nil {
			c.String(http.StatusInternalServerError, "template error")
		}
	})

	r.POST("/toyaml", func(c *gin.Context) {
		rawJSON := c.PostForm("jsonInput")
		yamlOut, err := handlers.ToYAML(rawJSON)
		data := PageData{RawInput: rawJSON, Formatted: yamlOut}
		if err != nil {
			data.Error = err.Error()
		}
		c.Status(http.StatusOK)
		c.Header("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(c.Writer, data); err != nil {
			c.String(http.StatusInternalServerError, "template error")
		}
	})

	r.POST("/extract/key", func(c *gin.Context) {
		rawJSON := c.PostForm("jsonInput")
		key := c.PostForm("key")
		out, err := handlers.ExtractKeyJSON(rawJSON, key)
		data := PageData{RawInput: rawJSON, Key: key, Formatted: out}
		if err != nil {
			data.Error = err.Error()
		}
		c.Status(http.StatusOK)
		c.Header("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(c.Writer, data); err != nil {
			c.String(http.StatusInternalServerError, "template error")
		}
	})

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Println("listening on http://localhost:8888")
	if err := r.Run(":8888"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
