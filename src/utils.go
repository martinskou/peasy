package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/BurntSushi/toml"
)

// Load a TOML file into an instance to T
func LoadToml[T any](filename string) (T, error) {
	var c T
	if _, err := toml.DecodeFile(filename, &c); err != nil {
		return c, err
	}
	return c, nil
}

// RenderTemplate tries to render a template consisting of one or more files.
// Data is accessible from the template. The result is written to w.
// Beware that only files mentioned is used! This func does not scan
// any folders for parts of templates...
func RenderTemplate(w http.ResponseWriter, data any, files ...string) error {

	af := []string{}
	for _, f := range files {
		af = append(af, "templates/"+f)
	}

	tmpl, err := template.ParseFiles(af...)
	if err != nil {
		log.Printf("Error parsing template : %s", err)
		return err
	}
	w.Header().Add("content-type", "text/html")
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template : %s", err)
		return err
	}
	return nil
}

// Pretty-print anything which can be marshalled as json
// otherwise print using printf
func Pprint(v any) {
	s, e := json.MarshalIndent(v, "", "  ")
	if e != nil {
		fmt.Printf("%#v\n", v)
	} else {
		fmt.Printf("%s\n", s)
	}
}

func Uuid() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		panic(err)
	}
	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}
