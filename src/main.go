package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Sysop struct {
	Name     string
	Password string
}

type Config struct {
	Port int
	//	Sysops   []Sysop
	//	AdminUrl string
}

func SetupLogging() {
	// logfile
	logfile, err := os.OpenFile("peasy.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	// log to both logfile and terminal
	multi := io.MultiWriter(logfile, os.Stdout)
	log.SetOutput(multi)
}

func main() {
	SetupLogging()

	// load config
	config, err := LoadToml[Config]("./config.toml")
	if err != nil {
		log.Fatalln("Unable to load config", err)
	}

	static_path := "./assets/"

	Store.LoadAll()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	//	r.Use(middleware.Logger)

	// serve static assets (css, js, jpg, png)
	fs := http.FileServer(http.Dir(static_path))
	r.Handle("/assets/*", http.StripPrefix("/assets/", fs))

	// server admin
	//r.Get("/admin/", AdminDashboardView)
	//r.Post("/admin/", AdminDashboardView)
	r.HandleFunc("/admin/", AdminDashboardView)
	r.HandleFunc("/admin/page/new/", AdminPageEditView)
	r.HandleFunc("/admin/page/save/", AdminPageSave)

	// serve pages
	// r.Get("/", PageView)
	// r.Post("/", PageView)
	r.Get("/*", PageView)
	// r.Post("/*", PageView)

	addr := fmt.Sprintf("0.0.0.0:%d", config.Port)
	log.Printf("Starting server at http://" + addr + "\n")
	http.ListenAndServe(addr, r)
}
