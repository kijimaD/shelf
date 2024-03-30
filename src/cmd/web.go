package cmd

import (
	"github.com/kijimaD/shelf/src/routers"
	"github.com/urfave/cli/v2"
)

var CmdWeb = &cli.Command{
	Name:        "web",
	Usage:       "web",
	Description: "start shelf app server",
	Action:      runWeb,
	Flags:       []cli.Flag{},
}

func runWeb(_ *cli.Context) error {
	if err := routers.RunServer(); err != nil {
		return err
	}

	return nil
}
