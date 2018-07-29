package download

import (
	"github.com/op/go-logging"
	"github.com/ielizaga/piv-go-gpdb/config"
	"fmt"
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

// Logger
var (
	log = logging.MustGetLogger("gpdb")
	FileNameContains = []string{
		"Red Hat Enterprise Linux",
		"RedHat Entrerprise Linux",
		"RedHat Enterprise Linux",
		"REDHAT ENTERPRISE LINUX",
		"RHEL",
		"Binary Installer for RHEL",
		"Binary Installer for RHEL 7"}
	IgnoreFileExtensions = []string{
		"RPM for RHEL",
		"Binary Installer for RHEL 6"}
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

type AuthToken struct {
	RefreshToken string `json:"refresh_token"`
	ReceivedToken string `json:"access_token"`
}

type DownloadOptions struct {
	Product string
	Version string
	Install bool
}

type userChoice struct {
	versionChoosen string
	releaseLink string
	DownloadURL string
	ProductFileURL string
	ProductFileName string
	ProductFileSize int64
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

type ResponseBody struct {
	Token AuthToken
	ProductList ProductObjects
	ReleaseList ReleaseObjects
	VersionList VersionObjects
	Env config.EnvObjects
	Flags DownloadOptions
	UserRequest userChoice
	EULALink string
	productFiles ProductFilesObjects
	download bool
}

// The main download web crawler
func (d *DownloadOptions) Download(e config.EnvObjects) error {
	log.Infof("Checking if the user is a valid user")

	// Initializing the response struct & store all the env configuration
	r := ResponseBody{}
	r.Env = e
	r.Flags = *d

	// Check the Authentication Token
	_, err := r.GetApi("GET", Authentication)
	if err != nil {
		return err
	}

	// Product list
	err = r.ExtractProduct()
	if err != nil {
		return err
	}

	// Release list
	err = r.ExtractRelease()
	if err != nil {
		return err
	}

	// Ask the user on what version do they want
	err = r.ShowAvailableVersion()
	if err != nil {
		return err
	}

	// Get all the files under that version
	err = r.ExtractDownloadURL()
	if err != nil {
		return err
	}

	// The users choice to what to download from that version
	err = r.WhichProduct()
	if err != nil {
		return err
	}

	// Extract the filename and the size of the file
	err = r.ExtractFileNamePlusSize()
	if err != nil {
		return err
	}

	// Accept the EULA
	log.Info("Accepting the EULA (End User License Agreement): " + r.EULALink)
	_, err = r.GetApi("POST", r.EULALink)
	if err != nil {
		return err
	}

	// All hard work is now done, lets download the version
	log.Info("Starting downloading of file: " + r.UserRequest.ProductFileName)
	if r.UserRequest.DownloadURL != "" {
		r.download = true
		_, err := r.GetApi("GET", r.UserRequest.DownloadURL)
		if err != nil {
			r.download = false
			return err
		}
		r.download = false
		log.Info("Downloaded file available at: " + r.UserRequest.ProductFileName)
	} else {
		return fmt.Errorf("download URL is blank, cannot download the product")
	}

	return nil
}
