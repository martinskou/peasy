package main

import (
	"fmt"
	"net/http"
)

/*


Top:
add new page, name _______ , parent _________ , fullurl ________

Next:
list of pages
page  edit, delete, move

*/

func NewPageForm() Form {
	f := NewForm("New page", "Create page", "/admin/page/new/")
	f.AddField(NewField("Name", "", "name", Text, ""))
	return f
}

func AdminDashboardView(w http.ResponseWriter, r *http.Request) {

	npf := NewPageForm()

	if r.Method == "GET" {

		data := struct { // google "golang anonymous struct"
			Title       string
			Page        Page
			Pages       []Page
			NewPageForm Form
		}{
			Title:       "Pages",
			Page:        Page{},
			Pages:       Store.TopPages(),
			NewPageForm: npf,
		}
		RenderTemplate(w, data, "layout.html", "admin_dashboard.html")

	} else {

	}
}

func AdminPageEditView(w http.ResponseWriter, r *http.Request) {

	id := r.FormValue("id")
	name := r.FormValue("name")

	p := Page{}
	if id == "" {
		p.Id = Uuid()
		p.Title = name
	} else {
		// get page from storage.
	}

	data := struct {
		Title string
		Page  Page
	}{
		Title: fmt.Sprintf("Page: %s", name),
		Page:  p,
	}

	RenderTemplate(w, data, "layout.html", "admin_page_edit.html")

}

func AdminPageSave(w http.ResponseWriter, r *http.Request) {
	p := Page{}

	p.Id = r.FormValue("id")
	p.URL = NewPath(r.FormValue("url"))
	p.Title = r.FormValue("title")
	p.Content = r.FormValue("content")

	Store.AddPage(p)

	http.Redirect(w, r, "/admin/", 301)
	return

}
