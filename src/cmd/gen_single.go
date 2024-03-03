package cmd

import (
	"errors"

	shelf "github.com/kijimaD/shelf/src"
	"github.com/urfave/cli/v2"
)

var CmdGenSingle = &cli.Command{
	Name:        "gensingle",
	Usage:       "generate file",
	Description: "generate file",
	Action:      runGenSingle,
	Flags:       []cli.Flag{},
}

func runGenSingle(c *cli.Context) error {
	if c.Args().Get(0) == "" {
		return errors.New("need argument")
	}
	filename := c.Args().Get(0)
	err := shelf.Register(filename)
	if err != nil {
		return err
	}

	return nil
}
