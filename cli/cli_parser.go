// Command line parser.
// The below package, parses the command line arguments for go gpdb software.
package cli

import (
	"github.com/urfave/cli"
	"github.com/ielizaga/piv-go-gpdb/download"
	"github.com/ielizaga/piv-go-gpdb/config"
	"github.com/op/go-logging"
	"fmt"
	"os"
	"github.com/ielizaga/piv-go-gpdb/environment"
	"github.com/ielizaga/piv-go-gpdb/install"
)

// Logger
var (
	log = logging.MustGetLogger("gpdb")
)

type ParserOptions struct {
	download download.DownloadOptions
	install install.InstallList
	env environment.EnvList
}

// The command line parser
func (p *ParserOptions) Parser() error {

	// New cli Instance
	app := cli.NewApp()

	// cli Metadata
	app.Name = "gpdb"
	app.Version = "3.0"
	app.Authors = []cli.Author{
		cli.Author{
			Name: "Ignacio Elizaga & Faisal Ali",
		},
	}
	app.Copyright = "MIT License"
	app.Usage = "Pivotal Greenplum Database Downloader and Installer based in Golang"

	// Reading & setting up the configuration file parameters
	e := config.EnvObjects{}
	err := e.Config()
	if err != nil {
		log.Errorf("%s", err)
	}

	// cli command and option
	app.Commands = []cli.Command{
		p.downloadParser(e),
		p.InstallParser(e),
		p.environmentParser(e)}

	// parse the cli arguments and error out if any error
	err = app.Run(os.Args)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	return nil
}
