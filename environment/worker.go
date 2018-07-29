package environment

import (
	"io/ioutil"
	"strings"
)

// Create environment file of this installation
//func (env *EnvList) CreateEnvFile(t string) error {
//
//	// Environment file fully qualified path
//	EnvFileName := env.Env.Install.EnvDir + "env_" + core.RequestedInstallVersion + "_" + t
//	log.Info("Creating environment file for this installation at: " + EnvFileName)
//
//	// Create the file
//	err := CreateFile(EnvFileName)
//	if err != nil {
//		return err
//	}
//
//	// Environment file content
//	var EnvFileContents []string
//	EnvFileContents = append(EnvFileContents, "export GPHOME=" + BinaryInstallLocation)
//	EnvFileContents = append(EnvFileContents, "export PYTHONPATH=$GPHOME/lib/python")
//	EnvFileContents = append(EnvFileContents, "export PYTHONHOME=$GPHOME/ext/python")
//	EnvFileContents = append(EnvFileContents, "export PATH=$GPHOME/bin:$PYTHONHOME/bin:$PATH")
//	EnvFileContents = append(EnvFileContents, "export LD_LIBRARY_PATH=$GPHOME/lib:$PYTHONHOME/lib:$LD_LIBRARY_PATH")
//	EnvFileContents = append(EnvFileContents, "export OPENSSL_CONF=$GPHOME/etc/openssl.cnf")
//	EnvFileContents = append(EnvFileContents, "export MASTER_DATA_DIRECTORY=" + GpInitSystemConfig.MasterDir + "/" + GpInitSystemConfig.ArrayName + "-1")
//	EnvFileContents = append(EnvFileContents, "export PGPORT=" + strconv.Itoa(GpInitSystemConfig.MasterPort))
//	EnvFileContents = append(EnvFileContents, "export PGDATABASE=" + GpInitSystemConfig.DBName)
//
//	// Write to EnvFile
//	err = core.WriteFile(EnvFileName, EnvFileContents)
//	if err != nil { return err }
//
//	return nil
//}


func (env *EnvList) ListInstallationContent(version string) error {

	log.Info("Found below matching environment file for the version: %v", version)

	return nil
}

func (env *EnvList) ListPrevInstallation(version string) error {

	log.Info("Checking if there is previous installation for the version: %v", version)

	// Listing all the files in the directory
	allfiles, err := ioutil.ReadDir(env.Env.Install.EnvDir)
	if err != nil {
		return err
	}

	// Match the version on this directory
	for _, file := range allfiles {

		if strings.Contains(file.Name(), version) {
			MatchingFilesInDir = append(MatchingFilesInDir, file.Name())
		}

	}


	return nil
}