package install

import (
	"github.com/ielizaga/piv-go-gpdb/core"
	"errors"
	"os"
)

// Check if the directory provided is readable and writeable
func (i *InstallList) DirValidator() error {

	// Check if the master & segment location is provided.
	if core.IsValueEmpty(i.Env.Install.MasterDataDirectory) {
		return errors.New("MASTER_DATA_DIRECTORY parameter missing in the config file, please set it and try again")
	}

	if core.IsValueEmpty(i.Env.Install.SegmentDataDirectory) {
		return errors.New("SEGMENT_DATA_DIRECTORY parameter missing in the config file, please set it and try again")
	}

	// Check if the master & segment location is readable and writable
	master_dir, err := core.DoesFileOrDirExists(i.Env.Install.MasterDataDirectory)
	if err != nil {
		return err
	}

	segment_dir, err := core.DoesFileOrDirExists(i.Env.Install.SegmentDataDirectory)
	if err != nil {
		return err
	}

	// if the file doesn't exists then let try creating it ...
	if !master_dir || !segment_dir {
		err := os.MkdirAll(i.Env.Install.MasterDataDirectory, 0775)
		if err != nil {
			return err
		}
		err = os.MkdirAll(i.Env.Install.SegmentDataDirectory, 0775)
		if err != nil {
			return err
		}
	}

	return nil
}
