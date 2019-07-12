package main

import (
	"encoding/json"
	// "fmt"
	"io/ioutil"
	// "net/http"
	"os"
	// "strings"
)

import (
// color "github.com/fatih/color"
)

const (
	dataFileName string = "data.json"
	userAgent    string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36"
)

var siteData = map[string]SiteData{}

// A SiteData struct for json datatype
type SiteData struct {
	ErrorMsg  string `json:"errorMsg"`
	ErrorType string `json:"errorType"`
	// Rank int`json:"rank"`
	URL            string `json:"url"`
	URLMain        string `json:"urlMain"`
	UsedUsername   string `json:"username_claimed"`
	UnusedUsername string `json:"username_unclaimed"`
}

func initializeSiteData() {
	jsonFile, err := os.Open(dataFileName)
	if err != nil {
		panic("Failed to open " + dataFileName)
	} else {
		defer jsonFile.Close()
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic("Error while read " + dataFileName)
	} else {
		json.Unmarshal([]byte(byteValue), &siteData)
	}
	return
}

func main() {
	initializeSiteData()
}
