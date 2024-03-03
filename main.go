package main

import (
	"os"

	"github.com/kijimaD/shelf/src/cmd"
)

func main() {
	app := cmd.NewMainApp()
	err := cmd.RunMainApp(app, os.Args...)
	panic(err)
}
