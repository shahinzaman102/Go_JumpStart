package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/go-chi/chi/v5"
)

// Page represents a wiki page with a title and body
type Page struct {
	Title string
	Body  []byte // Page content as raw bytes
}

// save writes the Page's body to a file in data/<Title>.txt
func (p *Page) save() error {
	if err := os.MkdirAll("data", 0700); err != nil { // ensure directory exists
		return err
	}
	return os.WriteFile("data/"+p.Title+".txt", p.Body, 0600) // file accessible only by owner
}

// loadPage reads a page from disk and returns a Page struct
func loadPage(title string) (*Page, error) {
	body, err := os.ReadFile("data/" + title + ".txt")
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// wikiTemplates is loaded lazily to avoid panics during init/tests
var wikiTemplates *template.Template

// LoadWikiTemplates parses the edit and view templates
func LoadWikiTemplates() error {
	var err error
	wikiTemplates, err = template.ParseFiles("templates/edit.html", "templates/view.html")
	return err
}

// validLink matches [[PageName]] style links in wiki text
var validLink = regexp.MustCompile(`\[(.+?)\]`)

// renderTemplate processes links and renders a wiki template
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	bodyStr := string(p.Body)

	// Convert [[PageName]] into <a href="/view/PageName">PageName</a>
	processed := validLink.ReplaceAllStringFunc(bodyStr, func(s string) string {
		m := validLink.FindStringSubmatch(s)
		link := m[1]
		escaped := url.PathEscape(link) // this encodes "page 1" to "page%201"
		return fmt.Sprintf(`<a href="/view/%s">%s</a>`, escaped, link)
	})

	err := wikiTemplates.ExecuteTemplate(w, tmpl+".html", struct {
		Title string
		Body  template.HTML
	}{
		Title: p.Title,
		Body:  template.HTML(processed), // trusted HTML
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ViewWiki handles GET /view/{title}
// If page exists, it renders it; else redirects to /edit/{title}
func ViewWiki(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "title")
	decodedTitle, _ := url.PathUnescape(title)

	p, err := loadPage(decodedTitle)
	if err != nil {
		http.Redirect(w, r, "/edit/"+decodedTitle, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// EditWiki handles GET /edit/{title}
// Loads existing page or prepares an empty one for editing
func EditWiki(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "title")
	decodedTitle, _ := url.PathUnescape(title)

	p, err := loadPage(decodedTitle)
	if err != nil {
		p = &Page{Title: decodedTitle} // new empty page
	}

	tmpl := template.Must(template.ParseFiles("templates/edit.html"))
	err = tmpl.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// SaveWiki handles POST /save/{title}
// Saves submitted form content and redirects to /view/{title}
func SaveWiki(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "title")
	decodedTitle, _ := url.PathUnescape(title)

	body := r.FormValue("body")
	p := &Page{Title: decodedTitle, Body: []byte(body)}

	if err := p.save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+decodedTitle, http.StatusFound)
}
