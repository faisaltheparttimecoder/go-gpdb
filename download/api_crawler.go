package download

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Extract all the Pivotal Network Product from the Product API page.
func (r *ResponseBody) ExtractProduct() error {

	log.Info("Obtaining the product ID")

	// Get the API from the Pivotal Products URL
	ProductApiResponse, err := r.GetApi("GET", Products)
	if err != nil {
		return err
	}

	// Store all the JSON on the Product struct
	json.Unmarshal(ProductApiResponse, &r.ProductList)

	// Return the struct
	return nil
}

// Extract all the Releases of the product with slug : pivotal-gpdb
func (r *ResponseBody) ExtractRelease() error {

	var ReleaseURL string
	var PivotalProduct string

	// Check what is the URL for the all the releases for product of our interest
	for _, product := range r.ProductList.Products {
		if product.Slug == ProductSlug {
			ReleaseURL = product.Links.Releases.Href
			PivotalProduct = product.Name
		}
	}

	log.Info("Obtaining the releases for product: " + PivotalProduct)

	// If we do find the release URL, lets continue
	if ReleaseURL != "" {

		// Extract all the releases
		ReleaseApiResponse, err := r.GetApi("GET", ReleaseURL)
		if err != nil {
			return err
		}

		// Store all the releases on the release struct
		json.Unmarshal(ReleaseApiResponse, &r.ReleaseList)

	} else { // Else lets error out
		return errors.New("cannot find any release URL for slug ID: " + ProductSlug)
	}

	// Return the release struct
	return nil
}


// From the user choice extract all the files available on that version
func (r *ResponseBody) ExtractDownloadURL() error{

	log.Info("Obtaining the files under the greenplum version: " + r.UserRequest.versionChoosen)

	// Extract all the files from the API
	VersionApiResponse, err := r.GetApi("GET", r.UserRequest.releaseLink)
	if err != nil {
		return fmt.Errorf("unable to extract download URL: %v", err)
	}

	// Load it to struct
	json.Unmarshal(VersionApiResponse, &r.VersionList)

	// Updating the EULA Acceptance link
	r.EULALink = r.VersionList.Links.Eula_acceptance.Href

	// Return the result
	return nil
}

// extract the filename and the size of the product that the user has choosen
func (r *ResponseBody) ExtractFileNamePlusSize () error {

	log.Info("Extracting the filename and the size of the product file.")

	// Obtain the JSON response of the product file API
	ProductFileApiResponse, err := r.GetApi("GET", r.UserRequest.ProductFileURL)
	if err != nil {
		return err
	}

	// Store it on JSON
	json.Unmarshal(ProductFileApiResponse, &r.productFiles)

	// Get the filename and the file size
	filename := strings.Split(r.productFiles.Product_file.Aws_object_key, "/")
	r.UserRequest.ProductFileName = filename[len(filename)-1]
	r.UserRequest.ProductFileSize = r.productFiles.Product_file.Size

	return err

}