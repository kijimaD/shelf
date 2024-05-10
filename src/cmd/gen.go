package cmd

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
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

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	metafile, err := os.OpenFile(shelf.MetaPath2, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(metafile)
	if err != nil {
		return err
	}
	oldMetas, err := shelf.GetMetas(string(bytes))
	if err != nil {
		return err
	}
	newMetas := shelf.Metas{}

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
		book, err := shelf.Register(f)
		if err != nil {
			return err
		}
		if err := book.AppendMeta(newMetas); err != nil {
			return err
		}
	}

	// 新しく追加された分だけにする
	for k, _ := range oldMetas {
		delete(newMetas, k)
	}

	if err := toml.NewEncoder(metafile).Encode(newMetas); err != nil {
		return err
	}

	return nil
}
