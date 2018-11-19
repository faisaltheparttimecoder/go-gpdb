package main

import (
	"encoding/json"
	"strings"
)

// All the PivNet Url's & Constants
const (
	EndPoint = "https://network.pivotal.io"
	Authentication = EndPoint + "/api/v2/authentication"
	RefreshToken = EndPoint + "/api/v2/authentication/access_tokens"
	Products =  EndPoint + "/api/v2/products"
	ProductSlug = "pivotal-gpdb" // we only care about this slug rest we ignore
	DBServer = "Database Server"
	GPCC = "Greenplum Command Center"
)

var (
	rx_gpcc = `Greenplum Command Center`
	rx_gpdb = `Greenplum Database (\d+\.)(\d+\.)(\d+)?(\.\d)?( Binary Installer)? ` +
		`for ((Red Hat Enterprise|RedHat Enterprise) Linux|RHEL).*7`
)

// Struct to where all the API response will be stored
type HrefType struct {
	Href string `json:"href"`
}

type LinksType struct {
	Self   HrefType `json:"self"`
	Releases   HrefType `json:"releases"`
	Product_files   HrefType `json:"product_files"`
	File_groups   HrefType `json:"file_groups"`
	Signature_file_download HrefType `json:"signature_file_download"`
	Eula_acceptance HrefType `json:"eula_acceptance"`
	User_groups HrefType `json:"user_groups"`
	Download HrefType `json:"download"`
}

type ProductObjType struct {
	Id int `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
	Logo_url string `json:"logo_url"`
	Links   LinksType `json:"_links"`
}

type ProductObjects struct {
	Products []ProductObjType `json:"products"`
	Links LinksType `json:"_links"`
}

type eulaType struct {
	Id int `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
	Links LinksType `json:"_links"`
}

type releaseObjType struct {
	Id int `json:"id"`
	Version string `json:"version"`
	Release_type string `json:"release_type"`
	Release_date string `json:"release_date"`
	Release_notes_url string `json:"release_notes_url"`
	Availability string `json:"availability"`
	Description string `json:"description"`
	Eula   eulaType `json:"eula"`
	Eccn string `json:"eccn"`
	License_exception string `json:"license_exception"`
	Controlled bool `json:"controlled"`
	Updated_at   string `json:"updated_at"`
	Software_files_updated_at string `json:"software_files_updated_at"`
	Links LinksType `json:"_links"`
}

type ReleaseObjects struct {
	Release []releaseObjType `json:"releases"`
	Links LinksType `json:"_links"`
}

type verProdType struct {
	Id int `json:"id"`
	Aws_object_key string `json:"aws_object_key"`
	File_version string `json:"file_version"`
	Sha256 string `json:"sha256"`
	Name string `json:"name"`
	Links LinksType `json:"_links"`
}

type verFileGroupType struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Product_files []verProdType `json:"product_files"`
}

type VersionObjType struct {
	Id int `json:"id"`
	Version string `json:"version"`
	Release_type string `json:"release_type"`
	Release_date string `json:"release_date"`
	Availability string `json:"availability"`
	Eula   eulaType `json:"eula"`
	End_of_support_date string `json:"end_of_support_date"`
	End_of_guidance_date string `json:"end_of_guidance_date"`
	Eccn string `json:"eccn"`
	License_exception string `json:"license_exception"`
	Controlled bool `json:"controlled"`
	Product_files []verProdType `json:"product_files"`
	File_groups []verFileGroupType `json:"file_groups"`
	Updated_at   string `json:"updated_at"`
	Software_files_updated_at string `json:"software_files_updated_at"`
	Links LinksType `json:"_links"`
}

type VersionObjects struct {
	VersionObjType
}

type ProductFilesObjType struct {
	Id int `json:"id"`
	Aws_object_key string `json:"aws_object_key"`
	Description string `json:"description"`
	Docs_url string `json:"docs_url"`
	File_transfer_status string `json:"file_transfer_status"`
	File_type string `json:"file_version"`
	Has_signature_file string `json:"has_signature_file"`
	Included_files []string `json:"included_files"`
	Md5 string `json:"md5"`
	Sha256 string `json:"sha256"`
	Name string `json:"name"`
	Ready_to_serve bool `json:"ready_to_serve"`
	Released_at string `json:"released_at"`
	Size int64 `json:"size"`
	System_requirements []string `json:"system_requirements"`
	Links LinksType `json:"_links"`
}


type ProductFilesObjects struct {
	Product_file ProductFilesObjType `json:"product_file"`
}

type userChoice struct {
	versionChoosen string
	releaseLink string
	DownloadURL string
	ProductFileURL string
	ProductFileName string
	ProductFileSize int64
}

type Responses struct {
	ProductList ProductObjects
	ReleaseList ReleaseObjects
	VersionList VersionObjects
	EULALink string
	productFiles ProductFilesObjects
	UserRequest userChoice
}

// Extract all the Pivotal Network Product from the Product API page.
func (r *Responses) extractProduct(token string) {
	Info("Obtaining the product ID")

	// Get the API from the Pivotal Products URL
	products :=  get(Products, token)

	// Load the information to the tile struct
	err := json.Unmarshal(products, &r.ProductList)
	if err != nil {
		Fatalf("Unable to unmarshal the products list: %v", err)
	}

	// Extract the releases
	r.extractRelease(token)

}

// Extract all the Releases of the product with slug : pivotal-gpdb
func (r *Responses) extractRelease(token string) {

	var ReleaseURL string
	var PivotalProduct string

	// Check what is the URL for the all the releases for product of our interest
	for _, product := range r.ProductList.Products {
		if product.Slug == ProductSlug {
			ReleaseURL = product.Links.Releases.Href
			PivotalProduct = product.Name
			Debugf("Release URL: %s", ReleaseURL)
			Debugf("Pivotal Product: %s", PivotalProduct)
		}
	}

	Infof("Obtaining the releases for product: %s", PivotalProduct)

	// If we do find the release URL, lets continue
	if ReleaseURL != "" {

		// Extract all the releases
		releases := get(ReleaseURL, token)

		// Store all the releases on the release struct
		err := json.Unmarshal(releases, &r.ReleaseList)
		if err != nil {
			Fatalf("Unable to unmarshal the release list: %v", err)
		}

	} else { // Else lets error out
		Fatalf("cannot find any release URL for slug ID: %s", ProductSlug)
	}

	// If the user has not provided the version, prompt to choose it
	// else if provided continue with download
	r.ShowAvailableVersion(token)
}

// From the user choice extract all the files available on that version
func (r *Responses) ExtractDownloadURL(token string) {

	Infof("Obtaining the files under the greenplum version: %s", r.UserRequest.versionChoosen)

	// Extract all the files from the API
	versionApiResponse := get(r.UserRequest.releaseLink, token)

	// Load it to struct
	err := json.Unmarshal(versionApiResponse, &r.VersionList)
	if err != nil {
		Fatalf("Unable to unmarshal the download url list: %v", err)
	}

	// Updating the EULA Acceptance link
	r.EULALink = r.VersionList.Links.Eula_acceptance.Href
	Debugf("EULA Link: %s", r.EULALink)

	// Extract the download URL if we can find it, else prompt it to user
	r.WhichProduct(token)

}

// extract the filename and the size of the product that the user has choosen
func (r *Responses) ExtractFileNamePlusSize (token string) {

	Info("Extracting the filename and the size of the product file.")

	// Obtain the JSON response of the product file API
	ProductFileApiResponse:= get(r.UserRequest.ProductFileURL, token)

	// Store it on JSON
	err := json.Unmarshal(ProductFileApiResponse, &r.productFiles)
	if err != nil {
		Fatalf("Unable to unmarshal the extract filename url list: %v", err)
	}

	// Get the filename and the file size
	filename := strings.Split(r.productFiles.Product_file.Aws_object_key, "/")
	r.UserRequest.ProductFileName = filename[len(filename)-1]
	r.UserRequest.ProductFileSize = r.productFiles.Product_file.Size
	Debugf("Product File Name: %s", r.UserRequest.ProductFileName)
	Debugf("Product File Size: %v", r.UserRequest.ProductFileSize)
}

// Download the the product from PivNet
func Download() {

	Info("Starting the program to download the product")

	// Get the authentication token
	token := getToken()

	// Initialize the struct
	r := new(Responses)

	// Extract all the product list / releases information
	r.extractProduct(token)

	// Accept the EULA
	Infof("Accepting the EULA (End User License Agreement): %s", r.EULALink)
	post(r.EULALink, token)

	// All hard work is now done, lets download the version
	Infof("Starting downloading of file: %s", r.UserRequest.ProductFileName)
	if r.UserRequest.DownloadURL != "" {
		downloadProduct(r.UserRequest.DownloadURL, token, *r)
		Infof("Downloaded file available at: %s", r.UserRequest.ProductFileName)
	} else {
		Fatalf("download URL is blank, cannot download the product")
	}

}