package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/urfave/cli/v2"
)

var CmdPrint = &cli.Command{
	Name:        "print",
	Usage:       "print",
	Description: "print all command description",
	Action:      runPrint,
	Flags:       []cli.Flag{},
}

func runPrint(c *cli.Context) error {
	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
	for _, cmd := range c.App.Commands {
		if cmd.Name == "help" {
			continue
		}
		fmt.Fprintln(writer, fmt.Sprintf("%s\t%s\t%s", cmd.Name, cmd.Description, cmd.Usage))
	}
	writer.Flush()

	return nil
}
