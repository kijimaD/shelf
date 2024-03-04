package cmd

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	shelf "github.com/kijimaD/shelf/src"
	"github.com/urfave/cli/v2"
)

var CmdGen = &cli.Command{
	Name:        "gen",
	Usage:       "generate [DIRECTORY] [EXT]",
	Description: "generate files in directory",
	Action:      runGen,
	Flags:       []cli.Flag{},
}

func runGen(c *cli.Context) error {
	if c.NArg() < 2 {
		return errors.New("gen [DIRECTORY] [EXT]")
	}
	dir := c.Args().Get(0)
	ext := c.Args().Get(1)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if filepath.Ext(file.Name())[1:] != ext {
			continue
		}
		filePath := filepath.Join(dir, file.Name())
		if err := shelf.Register(filePath); err != nil {
			return err
		}
	}

	return nil
}
