package environment

import (
	"io/ioutil"
	"strconv"
	"github.com/ielizaga/piv-go-gpdb/core"
	"fmt"
)

func (env *EnvList) configurationListing() {
	var header = []string{"Configuration Parameter", "Value"}
	data := [][]string{
		[]string{"Download Directory", env.Env.Download.DownloadDir},
		[]string{"Environment Directory", env.Env.Install.EnvDir},
		[]string{"Master Data Directory", env.Env.Install.MasterDataDirectory},
		[]string{"Segment Data Directory", env.Env.Install.SegmentDataDirectory},
		[]string{"Uninstall Directory", env.Env.Install.UninstallDir},
		[]string{"Program Parameters Directory", env.Env.Install.FutureRef},
		[]string{"Master Hostname", env.Env.Install.MasterHost},
		[]string{"Default Master Username", env.Env.Install.MasterUser},
		[]string{"Default Master Password", env.Env.Install.MasterPass},
		[]string{"Default GPMON Password", env.Env.Install.GpMonPass},
		[]string{"Total Segments", strconv.Itoa(env.Env.Install.TotalSegments)},
	}
	fmt.Println()
	// TODO: fix this after fixing the config path
	fmt.Printf("CONFIGURATION OPTIONS: \n\n")
	core.TableOutput(data, header)
}

// Program to list all the files in the download directory
func (env *EnvList) downloadDirectoryListing() error {

	// Reading the directory and looking through all the files
	files, err := ioutil.ReadDir(env.Env.Download.DownloadDir)
	if err != nil {
		return fmt.Errorf("error in reading directory: %v", err)
	}
	env.configurationListing()
	// Displaying all the output
	var data [][]string
	var header = []string{"Downloaded Files", "Mode", "Size", "Last Modified", "Directory"}
	for _, f := range files {
		t := f.ModTime()
		temp := []string{
			f.Name(),
			f.Mode().String(),
			strconv.FormatInt(f.Size(), 10),
			t.Format("2006-01-02 15:04:05") ,
			strconv.FormatBool(f.IsDir())}
		data = append(data, temp)
	}
	fmt.Println()
	fmt.Printf("DOWNLOADED FILES: %v \n\n", env.Env.Download.DownloadDir)
	core.TableOutput(data, header)
	fmt.Println()

	return nil
}


// Program to list all the files on the environment directory
// and what version of software is it installed
//func (env *EnvList) environmentFileListing() error {
//
//	// Reading the directory and looking through all the files
//	files, err := ioutil.ReadDir(env.Env.Install.EnvDir)
//	if err != nil {
//		return fmt.Errorf("error in reading directory: %v", err)
//	}
//
//	//// Displaying all the output
//	//var data [][]string
//	//var header = []string{"Downloaded Files", "Mode", "Size", "Last Modified", "IsItDirectory"}
//	//for _, f := range files {
//	//	t := f.ModTime()
//	//	temp := []string{
//	//		f.Name(),
//	//		f.Mode().String(),
//	//		strconv.FormatInt(f.Size(), 10),
//	//		t.Format("2006-01-02 15:04:05") ,
//	//		strconv.FormatBool(f.IsDir())}
//	//	data = append(data, temp)
//	//}
//	//fmt.Println()
//	//fmt.Printf("LISTING DOWNLOADED FILES FROM: %v \n\n", env.Env.Download.DownloadDir)
//	//core.TableOutput(data, header)
//	//fmt.Println()
//
//	return nil
//}
