package main

import (
	"github.com/op/go-logging"
	"github.com/ielizaga/piv-go-gpdb/logger"
	"github.com/ielizaga/piv-go-gpdb/cli"
)

var (
	log = logging.MustGetLogger("gpdb")
)

func main() {

	// Setup the logger
	logger.LoggerInit()

	// Parse the command line arguments and run appropriate program
	p := cli.ParserOptions{}
	err := p.Parser()
	if err != nil {
		log.Errorf("%s", err)
	}

}
