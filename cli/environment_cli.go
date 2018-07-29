package cli

import (
	"github.com/ielizaga/piv-go-gpdb/config"
	"github.com/urfave/cli"
	"fmt"
	"github.com/ielizaga/piv-go-gpdb/core"
)

// Parse the download flag
func (d *ParserOptions) environmentParser(e config.EnvObjects) cli.Command {
	return 		cli.Command{
		Name:    "env",
		Aliases: []string{"e"},
		Usage:   "Show all the software downloaded / installed on this machine",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "version, v",
				Value: "",
				Usage: "Which GPDB version software do you want the program to choose from the list?",
				Destination: &d.env.Flags.Version,
			},
			cli.BoolFlag{
				Name: "list, l",
				Usage: "List of all the product installed / downloaded & configuration",
				Destination: &d.env.Flags.Listing,
			},
		},
		Before: func(c *cli.Context) error {
			// If summary and version both are set together then error out
			if d.env.Flags.Listing && !core.IsValueEmpty(d.env.Flags.Version) {
				return fmt.Errorf("listing & version cannot be set together for environment listing")
			}
			return nil
		},
		After: func(c *cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			// Time to run the download code
			err := d.env.Environment(e)
			if err != nil {
				return err
			}
			return nil
		},
	}
}

