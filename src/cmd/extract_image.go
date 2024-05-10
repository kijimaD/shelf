package cmd

import (
	"errors"
	"fmt"
	"os"

	shelf "github.com/kijimaD/shelf/src"
	"github.com/urfave/cli/v2"
)

var CmdExtractImage = &cli.Command{
	Name:        "extract",
	Usage:       "extract [PDF]",
	Description: "extract image",
	Action:      extractImage,
	Flags:       []cli.Flag{},
}

func extractImage(c *cli.Context) error {
	if c.NArg() < 1 {
		return errors.New("extract [PDF]")
	}
	file := c.Args().Get(0)

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	book := shelf.NewBook(*f)
	base64, err := book.ExtractImageBase64()
	if err != nil {
		return err
	}

	fmt.Println(base64)

	return nil
}
