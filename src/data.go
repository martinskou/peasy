package main

import (
	"errors"
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

func (p Path) Parent() (Path, error) {
	n := Path{}
	if p.Len() > 0 {
		n.Parts = p.Parts[0 : len(p.Parts)-1]
		return n, nil
	} else {
		return n, errors.New("root element")
	}
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
	return "data/p" + strings.Join(p.Parts, "_") + ".toml"
}
func (p Path) ToLink() string {
	if len(p.Parts) == 0 {
		return "/"
	}
	return "/" + strings.Join(p.Parts, "/") + "/"
}

type Page struct {
	Id      string
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

type PageStore struct {
	Pages []Page
}

func (s *PageStore) AddPage(p Page) {
	s.Pages = append(s.Pages, p)
	p.Save()
}

func (s *PageStore) LoadAll() {

	files, err := ioutil.ReadDir("data/")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {

		n, err := LoadToml[Page]("data/" + f.Name())
		if err != nil {
			log.Println("Unable to load page", f.Name(), err)
			continue
		}

		s.Pages = append(s.Pages, n)

	}

	log.Printf("Loaded %d Pages\n", len(s.Pages))
}

func (s *PageStore) TopPages() []Page {
	tp := []Page{}
	for _, p := range s.Pages {
		if p.URL.Len() <= 1 {
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
		if p.URL.Contains(pp.URL) && p.Id != pp.Id {
			tp = append(tp, p)
		}
	}
	return tp
}

func (s *PageStore) FindParent(pp Page) *Page {
	pp_parent_path, err := pp.URL.Parent()
	if err != nil {
		return nil
	} else {
		page, found := s.FindPage(pp_parent_path)
		if found {
			return page
		} else {
			return nil
		}
	}
}

var Store PageStore
