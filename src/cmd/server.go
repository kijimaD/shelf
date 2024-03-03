package cmd

import (
	"github.com/urfave/cli/v2"
)

var CmdServer = &cli.Command{
	Name:        "server",
	Usage:       "start shelf server",
	Description: "start shelf server",
	Action:      runServer,
	Flags:       []cli.Flag{},
}

func runServer(_ *cli.Context) error {
	return nil
}
