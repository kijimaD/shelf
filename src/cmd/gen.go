package cmd

import (
	"errors"
	"io/ioutil"
	"os"
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
		filePath := filepath.Join(dir, file.Name())
		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		if _, err := shelf.Register(f); err != nil {
			return err
		}
	}

	return nil
}
