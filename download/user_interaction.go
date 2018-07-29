package download

import (
	"github.com/ielizaga/piv-go-gpdb/core"
	"fmt"
	"strings"
	"strconv"
)

// provide choice of which version to download
func (r *ResponseBody) ShowAvailableVersion() (error) {

	// Local storehouse
	var ReleaseOutputMap = make(map[string]string)
	var ReleaseVersion []string

	// Get all the releases from the ReleaseJson
	for _, release := range r.ReleaseList.Release {
		ReleaseOutputMap[release.Version] = release.Links.Self.Href
		ReleaseVersion = append(ReleaseVersion, release.Version)
	}

	// Check if the user provided version is on the list we have just extracted
	if core.Contains(ReleaseVersion, r.Flags.Version) {
		log.Info("Found the GPDB version \"" + r.Flags.Version + "\" on PivNet, continuing..")
		r.UserRequest.releaseLink = ReleaseOutputMap[r.Flags.Version]
		r.UserRequest.versionChoosen = r.Flags.Version

	} else { // If its not on the list then fallback to interactive mode

		// Print warning if the user did provide a value of the version
		if r.Flags.Version != "" {
			log.Warning("Unable to find the GPDB version \"" + r.Flags.Version + "\" on PivNet, failing back to interactive mode..\n")
		} else { // print a blank line
			fmt.Println()
		}

		// Sort all the keys
		var data [][]string
		var header = []string{"Index", "Product Version"}
		for index, version := range ReleaseVersion {
			temp :=[]string{strconv.Itoa(index+1), version}
			data = append(data, temp)
		}
		core.TableOutput(data, header)

		// Total accepted values that user can enter
		TotalOptions := len(ReleaseVersion)

		// Ask user for choice
		users_choice := core.Prompt_choice(TotalOptions)

		// Selected by the user
		r.UserRequest.releaseLink = ReleaseOutputMap[ReleaseVersion[users_choice-1]]
		r.UserRequest.versionChoosen = ReleaseVersion[users_choice-1]
	}

	return nil
}

// Ask user what file in that version are they interested in downloading
// Default is to download GPDB, GPCC and other with a option from parser
func (r *ResponseBody) WhichProduct() error {

	// Clearing up the buffer to ensure we are using a clean array and MAP
	var ProductOutputMap = make(map[string]string)
	var ProductOptions = []string{}

	// This is the correct API, all the files are inside the file group MAP
	for _, fg := range r.VersionList.File_groups {

		// GPDB Options
		if r.Flags.Product == "gpdb" {
			// Default Download which is download the GPDB for Linux (no choice)
			if strings.Contains(fg.Name, DBServer) {
				for _, file := range fg.Product_files {
					BreakLoop:
					for _, mustContain := range FileNameContains { // Check if it matches the criteria
						for _, mustNotContain := range IgnoreFileExtensions { // Ignore certain criteria
							if strings.Contains(file.Name, mustNotContain) {
								break BreakLoop // If found break the outerloop
							} else {
								if strings.Contains(file.Name, mustContain) {
									r.UserRequest.ProductFileURL = file.Links.Self.Href
									r.UserRequest.DownloadURL = file.Links.Download.Href
								}
							}
						}
					}
				}
			}
		}

		// GPCC option
		if r.Flags.Product == "gpcc" {
			if strings.Contains(fg.Name, GPCC) {
				for _, file := range fg.Product_files {
					ProductOutputMap[file.Name] = file.Links.Self.Href
					ProductOptions = append(ProductOptions, file.Name)
				}
			}
		}

		// Others or fallback method
		if r.Flags.Product == "gpextras" {
			for _, file := range fg.Product_files {
				ProductOutputMap[file.Name] = file.Links.Self.Href
				ProductOptions = append(ProductOptions, file.Name)
			}
		}

	}

	// If its GPCC or GPextras, then ask users for choice.
	if (r.Flags.Product == "gpextras" || r.Flags.Product == "gpcc") && len(ProductOptions) != 0 {
		fmt.Println()

		// Show all the prduct available
		var data [][]string
		var header = []string{"Index", "Product Description"}
		for index, product := range ProductOptions {
			temp :=[]string{strconv.Itoa(index+1), product}
			data = append(data, temp)
		}
		core.TableOutput(data, header)

		// Record the user choice
		TotalOptions := len(ProductOptions)
		users_choice := core.Prompt_choice(TotalOptions)
		version_selected_url := ProductOutputMap[ProductOptions[users_choice-1]]

		// Store the user selection
		r.UserRequest.ProductFileURL = version_selected_url
		r.UserRequest.DownloadURL = version_selected_url + "/download"
	}

	return nil
}
