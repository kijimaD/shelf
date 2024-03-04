package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func NewMainApp() *cli.App {
	app := cli.NewApp()
	app.Name = "shelf"
	app.Usage = "shelf -h"
	app.Description = "shelf is simple document management tool."
	app.Version = "v0.0.0"
	app.EnableBashCompletion = true
	app.DefaultCommand = CmdWeb.Name
	app.Commands = []*cli.Command{
		CmdWeb,
		CmdGen,
		CmdGenSingle,
	}

	return app
}

func RunMainApp(app *cli.App, args ...string) error {
	err := app.Run(args)
	if err != nil {
		return fmt.Errorf("メインコマンドの起動に失敗した: %w", err)
	}

	return nil
}
