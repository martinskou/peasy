package main

import (
	"fmt"
	"net/http"
	"strings"
)

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
