package cmd

import (
	"errors"
	"io/ioutil"
	"log"
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
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) != shelf.DocExtension {
			continue
		}
		filePath := filepath.Join(dir, file.Name())
		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		if _, err := shelf.Register(f); err != nil {
			switch err {
			case shelf.ErrAlreadyFormatted:
				log.Printf("スキップ: %v %s", err, f.Name())
				continue
			}
			return err
		}
	}

	return nil
}
