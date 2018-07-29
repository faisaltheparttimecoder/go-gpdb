package cli

import (
	"github.com/ielizaga/piv-go-gpdb/config"
	"github.com/urfave/cli"
	"github.com/ielizaga/piv-go-gpdb/core"
	"fmt"
)

var (
	acceptedDownloadProduct = []string{"gpdb", "gpcc", "gpextras"}
)

// Parse the download flag
func (d *ParserOptions) downloadParser(e config.EnvObjects) cli.Command {
	return 		cli.Command{
		Name:    "download",
		Aliases: []string{"d"},
		Usage:   "Download the products from pivotal network (pivotal software hosting site)",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "product, p",
				Value: "gpdb",
				Usage: "What product do you want to download? [OPTIONS: gpdb, gpcc, gpextras]",
				Destination: &d.download.Product,
			},
			cli.StringFlag{
				Name: "version, v",
				Value: "",
				Usage: "OPTIONAL: Which GPDB version software do you want the program to choose from the list?",
				Destination: &d.download.Version,
			},
			cli.BoolFlag{
				Name: "install, i",
				Usage: "OPTIONAL: Install after downloaded (Only works with \"gpdb\")",
				Destination: &d.download.Install,
			},
		},
		Before: func(c *cli.Context) error {
			// If its valid then we are going to store it
			if !core.Contains(acceptedDownloadProduct, d.download.Product) {
				return fmt.Errorf("the product value \"%s\" is not a valid option", d.download.Product)
			}
			// Checking if users have used wrong parameter related to install
			if d.download.Install && d.download.Product != "gpdb" {
				return fmt.Errorf("install flag only works with product flag \"gpdb\"")
			}
			return nil
		},
		After: func(c *cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			// Time to run the download code
			err := d.download.Download(e)
			if err != nil {
				return err
			}
			return nil
		},
	}
}

