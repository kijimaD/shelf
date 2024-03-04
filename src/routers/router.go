package routers

import (
	"embed"
	"fmt"
	"net/http"
)

//go:embed static/*
var staticFS embed.FS

func RunServer() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "top page")
	})

	http.Handle("/static/", http.FileServer(http.FS(staticFS)))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", Config.Port), nil); err != nil {
		return err
	}

	return nil
}
