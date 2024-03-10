package cmd

import (
	"errors"
	"os"

	shelf "github.com/kijimaD/shelf/src"
	"github.com/urfave/cli/v2"
)

var CmdGenSingle = &cli.Command{
	Name:        "gensingle",
	Usage:       "gensingle [FILE]",
	Description: "generate file",
	Action:      runGenSingle,
	Flags:       []cli.Flag{},
}

func runGenSingle(c *cli.Context) error {
	if c.Args().Get(0) == "" {
		return errors.New("gensingle [DOC PATH]")
	}
	filename := c.Args().Get(0)
	f, err := os.Open(filename)
	_, err = shelf.Register(f)
	if err != nil {
		return err
	}

	return nil
}
