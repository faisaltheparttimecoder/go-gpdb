package environment

import (
	"github.com/ielizaga/piv-go-gpdb/core"
	"errors"
	"github.com/op/go-logging"
	"github.com/ielizaga/piv-go-gpdb/config"
)

var (
	log = logging.MustGetLogger("gpdb")
)

// All the struct (the placeholders)
type EnvOptions struct {
	Version string
	Listing bool
}

type EnvList struct {
	Flags EnvOptions
	Env config.EnvObjects
}

// The main environment program
func (env *EnvList) Environment(e config.EnvObjects) error {

	// Store all the configuration on the struct so that they can be used by rest of the program
	env.Env = e

	// If no version is provided then the user is requesting to list all the environment installed
	if env.Flags.Version == "" {

		log.Info("listing all the environment version installed on this cluster")

		// Show all the files that was downloaded
		err := env.downloadDirectoryListing()
		if err != nil {
			return err
		}

		// set the chosen environment
		if core.IsValueEmpty("") {
			return errors.New("there is no any installation of GPDB, please install the product to list the environment here")
		} else {
			//err = SettheChoosenEnv(chosenEnvFile, e.Flags.Version)
			//if err != nil { return err }
		}


	} else { // he is checking for a specific version

		log.Info("listing all the environment that has been installed with version: " + env.Flags.Version)

	}

	log.Info("Exiting ..... ")

	return nil
}
