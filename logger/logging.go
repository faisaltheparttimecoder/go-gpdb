package logger

import (
	"github.com/op/go-logging"
	"os"
)

var (
	log    = logging.MustGetLogger("gpdb")
	format = logging.MustStringFormatter(
		`%{color}%{time:2006-01-02 15:04:05.000}:%{level:s} > %{color:reset}%{message}`,
	)
)

// Logger for go-logging package
// create backend for os.Stderr, set the format and update the logger to what logger to be used
func LoggerInit() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
}