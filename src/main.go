package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// RenderTemplate tries to render a template consisting of one or more files.
// Data is accessible from the template. The result is written to w.
func RenderTemplate(w http.ResponseWriter, data any, files ...string) {

	af := []string{}
	for _, f := range files {
		af = append(af, "templates/"+f)
	}

	tmpl, err := template.ParseFiles(af...)
	if err != nil {
		log.Printf("Error parsing template : %s", err)
		return
	}
	w.Header().Add("content-type", "text/html")
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template : %s", err)
	}
}

func FrontView(w http.ResponseWriter, r *http.Request) {

	fmt.Println("home_view", r.Method, r.URL.String())

	if r.Method == "GET" {

		data := struct { // google "golang anonymous struct"
			Title string
			Pages []Page
		}{
			Title: "Pages",
			Pages: Store.TopPages(),
		}
		RenderTemplate(w, data, "layout.html", "home.html")

	} else {

		fmt.Println("SEARCH")

		kw := r.FormValue("search")
		fmt.Println(kw)

		data := struct { // google "golang anonymous struct"
			Title   string
			Results []Page
			Keyword string
		}{
			Title:   "Search",
			Results: []Page{},
			Keyword: kw,
		}

		for _, n := range Store.Pages {
			if strings.Contains(n.Title, kw) {

				fmt.Println("FOUND", n)
				data.Results = append(data.Results, n)

			}
		}
		RenderTemplate(w, data, "layout.html", "home.html")
	}
}

func PageView(w http.ResponseWriter, r *http.Request) {

	fmt.Println("generic_view", r.Method, r.URL.String())

	if r.Method == "GET" {

		path := NewPath(r.URL.Path)

		fmt.Println(r.URL.Query().Has("edit"))

		page, found := Store.FindPage(path)

		if !found {
			page = &Page{}
			page.URL = path
			page.Title = r.URL.Path
		}

		data := struct {
			Title    string
			Page     Page
			Children []Page
			Parent   *Page
			New      bool
			Edit     bool
		}{
			Title:    r.URL.Path,
			Page:     *page,
			Children: Store.FindChildren(*page),
			Parent:   Store.FindParent(*page),
			New:      !found,
			Edit:     r.URL.Query().Has("edit"),
		}
		RenderTemplate(w, data, "layout.html", "page.html")

	} else {

		// save data
		n := Page{}
		n.URL = NewPath(r.URL.String())
		n.Title = r.FormValue("title")
		n.Content = r.FormValue("content")
		n.Save()

		http.Redirect(w, r, "/", http.StatusFound)

	}

}

type Config struct {
	Port int
}

func LoadConfig(filename string) Config {
	c := Config{}
	if _, err := toml.DecodeFile(filename, &c); err != nil {
		log.Fatalln(err)
	}
	return c
}

func main() {

	config := LoadConfig("./config.toml")

	static_path := "./assets/"

	Store.LoadAll()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	//	r.Use(middleware.Logger)

	fs := http.FileServer(http.Dir(static_path))
	r.Handle("/assets/*", http.StripPrefix("/assets/", fs))

	r.Get("/", FrontView)
	r.Post("/", FrontView)

	r.Get("/*", PageView)
	r.Post("/*", PageView)

	addr := fmt.Sprintf("0.0.0.0:%d", config.Port)
	log.Printf("Starting server http://" + addr + "\n")
	http.ListenAndServe(addr, r)
}
