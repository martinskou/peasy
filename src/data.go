package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type Page struct {
	URL     string
	Title   string
	Content string
}

func (n Page) URLToFilename(url string) string {
	filename := strings.ReplaceAll(strings.ToLower(url), "/", "_")
	return "data/" + filename + ".toml"
}

func (n Page) Save() {

	f, err := os.OpenFile(n.URLToFilename(n.URL), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
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

func (n *Page) Load(url string) error {
	if _, err := toml.DecodeFile(n.URLToFilename(url), &n); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

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

var Store PageStore
