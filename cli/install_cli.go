package cli

import (
	"github.com/ielizaga/piv-go-gpdb/config"
	"github.com/urfave/cli"
	"fmt"
	"github.com/ielizaga/piv-go-gpdb/core"
)

var AcceptedInstallProduct = []string{"gpdb", "gpcc"}

// Parse the download flag
func (d *ParserOptions) InstallParser(e config.EnvObjects) cli.Command {
	return 		cli.Command{
		Name:    "install",
		Aliases: []string{"i"},
		Usage:   "Install the software that is downloaded using the download option (only works for gpdb & gpcc)",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "product, p",
				Value: "gpdb",
				Usage: "What product do you want to install? [OPTIONS: gpdb, gpcc]",
				Destination: &d.install.Flags.Product,
			},
			cli.StringFlag{
				Name: "version, v",
				Value: "",
				Usage: "Which GPDB version software do you want the program to choose from the list?",
				Destination: &d.install.Flags.Version,
			},
			cli.StringFlag{
				Name: "ccversion, c",
				Value: "",
				Usage: "Which GPCC version do you want to install?",
				Destination: &d.install.Flags.CCVersion,
			},
		},
		Before: func(c *cli.Context) error {

			// Check if the provided parameter for product is valid
			if !core.Contains(AcceptedInstallProduct, d.install.Flags.Product) {
				return fmt.Errorf("invalid parameter for product, please type in valid options")
			}

			// If the version flag is blank, then error out
			if core.IsValueEmpty(d.install.Flags.Version) {
				return fmt.Errorf("GPDB Version missing, Please provide the version that needs to be used for installation")
			}

			// If summary and version both are set together then error out
			if !core.IsValueEmpty(d.install.Flags.CCVersion) && core.IsValueEmpty(d.install.Flags.Version) {
				return fmt.Errorf("to install command center, a database version is required")
			}

			return nil
		},
		After: func(c *cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			// Time to run the download code
			err := d.install.InstallSingleNodeGPDB(e)
			if err != nil {
				return err
			}
			return nil
		},
	}
}

