package install

import (
	"github.com/op/go-logging"
	"github.com/ielizaga/piv-go-gpdb/config"
	"github.com/ielizaga/piv-go-gpdb/environment"
)

var (
	log = logging.MustGetLogger("gpdb")
)

// All the struct (the placeholders)
type InstallOptions struct {
	Product string
	Version string
	CCVersion string
}

type InstallList struct {
	Flags InstallOptions
	Env config.EnvObjects
	EnvManager environment.EnvList
}

// GPDB Single Node Installer
func (i *InstallList) InstallSingleNodeGPDB(e config.EnvObjects) error {

	// Store all the environment configuration
	i.Env = e

	// TODO: calling install after download
	//// If the install is called from download command the set default values
	//if !core.IsValueEmpty(core.RequestedDownloadVersion) {
	//	core.RequestedInstallVersion = core.RequestedDownloadVersion
	//}

	// Validate the master & segment exists and is readable
	//err := i.DirValidator()
	//if err != nil {
	//	return err
	//}

	// Check if there is already a previous version of the same version
	err := i.EnvManager.ListPrevInstallation(i.Flags.Version)
	if err != nil {
		return err
	}

	log.Info("Installation of GPDB software has been completed successfully")

	return nil

}