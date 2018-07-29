package config

import (
	"github.com/op/go-logging"
	"fmt"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/ielizaga/piv-go-gpdb/core"
)

// Logger
var (
	log = logging.MustGetLogger("gpdb")
)

type coreType struct {
	BaseDir string `yaml:"BASE_DIR"`
	TempDir string `yaml:"TEMP_DIR"`
	AppName string `yaml:"APPLICATION_NAME"`
}

type downloadType struct {
	ApiToken string `yaml:"API_TOKEN"`
	DownloadDir string `yaml:"DOWNLOAD_DIR"`
}

type installType struct {
	EnvDir string `yaml:"ENV_DIR"`
	UninstallDir string `yaml:"UNINTSALL_DIR"`
	FutureRef string `yaml:"FUTUREREF_DIR"`
	MasterHost string `yaml:"MASTER_HOST"`
	MasterUser string `yaml:"MASTER_USER"`
	MasterPass string `yaml:"MASTER_PASS"`
	GpMonPass string `yaml:"GPMON_PASS"`
	MasterDataDirectory string `yaml:"MASTER_DATA_DIRECTORY"`
	SegmentDataDirectory string `yaml:"SEGMENT_DATA_DIRECTORY"`
	TotalSegments int `yaml:"TOTAL_SEGMENT"`
}

type EnvObjects struct {
	Core coreType `yaml:"CORE"`
	Download downloadType `yaml:"DOWNLOAD"`
	Install installType `yaml:"INSTALL"`
}

// Check if the value exists else place default values for them
func (e *EnvObjects) valuesExistsOrSetDefault() error {

	// Base Directory
	if core.IsValueEmpty(e.Core.BaseDir) {
		log.Warning(fmt.Errorf("BASE_DIR parameter missing in the config file, setting to default"))
		e.Core.BaseDir = "/home/gpadmin/"
	}

	// App name
	if core.IsValueEmpty(e.Core.AppName) {
		log.Warning(fmt.Errorf("APPLICATION_NAME parameter missing in the config file, setting to default"))
		e.Core.AppName = "gpdbinstall"
	}

	// Temp Directory
	if core.IsValueEmpty(e.Core.TempDir) {
		log.Warning(fmt.Errorf("TEMP_DIR parameter missing in the config file, setting to default"))
		e.Core.TempDir = "/temp/"
	}

	// Download Directory
	if core.IsValueEmpty(e.Download.DownloadDir) {
		log.Warning(fmt.Errorf("DOWNLOAD_DIR parameter missing in the config file, setting to default"))
		e.Download.DownloadDir = "/download/"
	}

	// Env Directory
	if core.IsValueEmpty(e.Install.EnvDir) {
		log.Warning(fmt.Errorf("ENV_DIR parameter missing in the config file, setting to default"))
		e.Install.EnvDir = "/env/"
	}

	// Uninstall Directory
	if core.IsValueEmpty(e.Install.UninstallDir) {
		log.Warning(fmt.Errorf("UNINTSALL_DIR parameter missing in the config file, setting to default"))
		e.Install.UninstallDir = "/uninstall/"
	}

	// Future Reference Directory
	if core.IsValueEmpty(e.Install.FutureRef) {
		log.Warning(fmt.Errorf("UNINTSALL_DIR parameter missing in the config file, setting to default"))
		e.Install.FutureRef = "/future_reference/"
	}

	// If there is no API Token then error out
	if core.IsValueEmpty(e.Download.ApiToken) {
		return fmt.Errorf("cannot find value for \"API_TOKEN\", check \"config.yml\"")
	}

	return nil
}

func (e *EnvObjects) createDirectory(location string) error {
	existsOrNot, err := core.DoesFileOrDirExists(location)
	if err != nil {
		return err
	}
	if !existsOrNot {
		log.Warning("Directory \""+ location + "\" does not exists, creating one")
		err:= os.MkdirAll(location, 0755)
		if err != nil {return err}
	}

	return nil
}

// Read the configuration file and create directory if not exists
// or set default values if values are missing
func (e *EnvObjects) checkConfig() error {

	// Check if the config have the parameters or values set
	// else set it to default
	log.Info("Checking if the parameters needed for the program exists")
	err := e.valuesExistsOrSetDefault()
	if err != nil {
		return err
	}

	// All Directories needed by the program
	log.Info("Checking if the directory needed for the program exists")
	e.Core.BaseDir = e.Core.BaseDir + e.Core.AppName                    // Base Directory
	e.Core.TempDir =  e.Core.BaseDir + e.Core.TempDir					// Temp Directory
	e.Download.DownloadDir =  e.Core.BaseDir + e.Download.DownloadDir   // Product Download directory
	e.Install.EnvDir = e.Core.BaseDir + e.Install.EnvDir				// Install Directory
	e.Install.UninstallDir = e.Core.BaseDir + e.Install.UninstallDir	// Uninstall Directory
	e.Install.FutureRef = e.Core.BaseDir + e.Install.FutureRef			// Future Reference Directory

	// Check if all the directory exists or create one
	allDirectory := []string{
		e.Core.TempDir,
		e.Download.DownloadDir,
		e.Install.EnvDir,
		e.Install.UninstallDir,
		e.Install.FutureRef,
	}
	for _, elem := range allDirectory {
		err := e.createDirectory(elem)
		if err != nil {
			return err
		}
	}

	return nil
}


// Configuration file reader.
func (e *EnvObjects) Config() error {

	home := os.Getenv("HOME")
	// TODO: make this a environment variable or default location
	configFile := home + "/.config.yml"

	log.Infof("Reading the configuration file: %s", configFile)

	// Read the config file and store the value on a struct
	source, err := ioutil.ReadFile(configFile)
	if err != nil {return err}
	yaml.Unmarshal(source, &e)

	// Creating Directory needed for the program
	err = e.checkConfig()
	if err != nil {return err}

	return nil
}