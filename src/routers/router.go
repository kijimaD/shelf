package routers

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	shelf "github.com/kijimaD/shelf/src"
)

//go:embed static/*
var staticFS embed.FS

type Content struct {
	Views []shelf.View
}

func RunServer() error {
	http.HandleFunc("/", indexHandler)
	http.Handle("/static/", http.FileServer(http.FS(staticFS)))

	dirname, err := filepath.Abs(".")
	if err != nil {
		log.Fatalf("Could not get absolute path to directory: %s: %s", dirname, err.Error())
	}
	fileServer := http.StripPrefix("/dir/", http.FileServer(http.Dir(dirname)))
	http.Handle("/dir/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/dir/") {
			fileServer.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	}))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", Config.Port), nil); err != nil {
		return err
	}

	return nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/index.html"))
	q := r.URL.Query().Get("q")

	views := shelf.GenerateViews(Config.ServeBase)
	if q != "" {
		views = shelf.FilterViewsByTag(q, views)
	}

	content := Content{Views: views}
	if err := t.Execute(w, content); err != nil {
		log.Fatal(err)
	}
}
