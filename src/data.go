package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type Path struct {
	Parts []string
}

func NewPath(URL string) Path {
	p := Path{}
	for _, i := range strings.Split(strings.ToLower(URL), "/") {
		if len(i) > 0 {
			p.Parts = append(p.Parts, i)
		}
	}
	return p
}

func (p Path) Len() int {
	return len(p.Parts)
}

func (p Path) Contains(a Path) bool {
	if a.Len() > p.Len() {
		return false
	}
	for i := 0; i < a.Len(); i++ {
		if p.Parts[i] != a.Parts[i] {
			return false
		}
	}
	return true
}

func (p Path) Equals(a Path) bool {
	if a.Len() != p.Len() {
		return false
	}
	for i := 0; i < a.Len(); i++ {
		if p.Parts[i] != a.Parts[i] {
			return false
		}
	}
	return true
}

func (p Path) ToFilename() string {
	return "data/" + strings.Join(p.Parts, "_")
}
func (p Path) ToLink() string {
	return "/" + strings.Join(p.Parts, "/") + "/"
}

type Page struct {
	URL     Path
	Title   string
	Content string
}

func (n Page) URLToFilename(url string) string {
	filename := strings.ReplaceAll(strings.ToLower(url), "/", "_")
	return "data/" + filename + ".toml"
}

func (n Page) Save() {

	f, err := os.OpenFile(n.URL.ToFilename(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}

	err = toml.NewEncoder(f).Encode(n)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}
}

/*
func (n *Page) Load(url string) error {
	if _, err := toml.DecodeFile(n.URLToFilename(url), &n); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
*/

type PageStore struct {
	Pages []Page
}

func (s *PageStore) LoadAll() {

	files, err := ioutil.ReadDir("data/")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {

		n := Page{}
		if _, err := toml.DecodeFile("data/"+f.Name(), &n); err != nil {
			fmt.Println(err)
		}

		s.Pages = append(s.Pages, n)

	}

	log.Printf("Loaded %d Pages\n", len(s.Pages))
}

func (s *PageStore) TopPages() []Page {
	tp := []Page{}
	for _, p := range s.Pages {
		if p.URL.Len() == 2 {
			tp = append(tp, p)
		}
	}
	return tp
}

func (s *PageStore) FindPage(path Path) (*Page, bool) {
	for _, p := range s.Pages {
		if p.URL.Equals(path) {
			return &p, true
		}
	}
	return nil, false
}

func (s *PageStore) FindChildren(pp Page) []Page {
	tp := []Page{}
	for _, p := range s.Pages {
		if p.URL.Contains(pp.URL) {
			tp = append(tp, p)
		}
	}
	return tp
}

func (s *PageStore) FindParent(pp Page) *Page {
	if len(s.Pages) > 0 {
		return &s.Pages[0]
	}
	return nil
}

var Store PageStore
