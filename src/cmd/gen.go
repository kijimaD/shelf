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
	Usage:       "generate [DIRECTORY]",
	Description: "generate files in directory",
	Action:      runGen,
	Flags:       []cli.Flag{},
}

func runGen(c *cli.Context) error {
	if c.NArg() < 1 {
		return errors.New("gen [DIRECTORY]")
	}
	dir := c.Args().Get(0)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) != shelf.DocExtension {
			continue
		}
		filePath := filepath.Join(dir, file.Name())
		if err := shelf.Register(filePath); err != nil {
			return err
		}
	}

	return nil
}
