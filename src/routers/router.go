package routers

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	shelf "github.com/kijimaD/shelf/src"
)

//go:embed static/*
var staticFS embed.FS

var tmpl *template.Template

type Content struct {
	Views []shelf.View
	Tags  []string // ユニークなタグ一覧
}

func RunServer() error {
	http.HandleFunc("/", indexHandler)
	http.Handle("/static/", http.FileServer(http.FS(staticFS)))

	dirname, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	fileServer := http.StripPrefix("/dir/", http.FileServer(http.Dir(dirname)))
	http.Handle("/dir/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/dir/") {
			fileServer.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	}))

	tmpl, err = template.ParseFS(staticFS, "static/templates/*.html")
	if err != nil {
		return err
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", Config.Port), nil); err != nil {
		return err
	}

	return nil
}

const (
	defaultLimit = 20
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	limit := r.URL.Query().Get("limit")

	views := shelf.GenerateViews(Config.ServeBase)
	if q != "" {
		views = shelf.FilterViewsByTag(q, views)
	}

	limitLen := defaultLimit
	if limit != "" {
		i, err := strconv.Atoi(limit)
		if err != nil {
			log.Print(err)
		}
		limitLen = i
	}
	idx := min(len(views), limitLen)
	views = views[:idx]

	content := Content{
		Views: views,
		Tags:  shelf.UniqTags(views),
	}
	if err := tmpl.ExecuteTemplate(w, "index.html", content); err != nil {
		log.Fatal(err)
	}
}
