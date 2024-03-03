package cmd

import "github.com/urfave/cli/v2"

var CmdGen = &cli.Command{
	Name:        "gen",
	Usage:       "generate file",
	Description: "generate file",
	Action:      runGen,
	Flags:       []cli.Flag{},
}

func runGen(_ *cli.Context) error {
	return nil
}
