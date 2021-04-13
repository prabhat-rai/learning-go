package main

import (
	"html/template"
	"path/filepath"
	"prabhat-rai.in/snippetbox/pkg/forms"
	"prabhat-rai.in/snippetbox/pkg/models"
	"time"
)

type templateData struct {
	CSRFToken string
	CurrentYear int
	Flash string
	Form *forms.Form
	IsAuthenticated bool
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

func noescape(str string) template.HTML {
	return template.HTML(str)
}

var functions = template.FuncMap{
	"humanDate": humanDate,
	"noEscape": noescape,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// Loop through the pages one-by-one.
	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
