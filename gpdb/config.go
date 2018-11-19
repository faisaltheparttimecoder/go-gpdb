package main

import (
	"github.com/jinzhu/configor"
)

// Struct that store the configuration file for the program to run
var Config = struct {
	CORE struct {
		APPLICATIONNAME string `yaml:"APPLICATION_NAME"`
		OS              string `yaml:"OS"`
		ARCH            string `yaml:"ARCH"`
		GOBUILD         string `yaml:"GO_BUILD"`
		BASEDIR         string `yaml:"BASE_DIR"`
		TEMPDIR         string `yaml:"TEMP_DIR"`
	} `yaml:"CORE"`
	DOWNLOAD struct {
		APITOKEN    string `yaml:"API_TOKEN"`
		DOWNLOADDIR string `yaml:"DOWNLOAD_DIR"`
	} `yaml:"DOWNLOAD"`
	INSTALL struct {
		ENVDIR               string `yaml:"ENV_DIR"`
		UNINTSALLDIR         string `yaml:"UNINTSALL_DIR"`
		FUTUREREFDIR         string `yaml:"FUTUREREF_DIR"`
		MASTERHOST           string `yaml:"MASTER_HOST"`
		MASTERUSER           string `yaml:"MASTER_USER"`
		MASTERPASS           string `yaml:"MASTER_PASS"`
		GPMONPASS            string `yaml:"GPMON_PASS"`
		MASTERDATADIRECTORY  string `yaml:"MASTER_DATA_DIRECTORY"`
		SEGMENTDATADIRECTORY string `yaml:"SEGMENT_DATA_DIRECTORY"`
		TOTALSEGMENT         int    `yaml:"TOTAL_SEGMENT"`
	} `yaml:"INSTALL"`
}{}

// Load the configuration file to the memory
func config() {
	configFile := "config.yml"
	Debugf("Reading the configuration file and loading to the memory: %s", configFile)
	configor.Load(&Config, configFile)
}