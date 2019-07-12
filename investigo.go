package main

import (
	"encoding/json"
	"fmt"
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
	ErrorType string `json:"errorType"`
	ErrorMsg  string `json:"errorMsg"`
	// Rank int`json:"rank"`
	URL     string `json:"url"`
	URLMain string `json:"urlMain"`
	// UsedUsername   string `json:"username_claimed"`
	// UnusedUsername string `json:"username_unclaimed"`
}

// Options contains command line arguments data object
type Options struct {
	color           bool
	updateBeforeRun bool
	verbose         bool
}

var options Options

// Investigo investigate if username exists on social media.
func Investigo(name string) {

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

func contains(array []string, targets ...string) (bool, int) {
	for index, item := range array {
		for _, target := range targets {
			if item == target {
				return true, index
			}
		}
	}
	return false, -1
}

func main() {
	args := os.Args[1:]
	var argIndex int

	options.color, argIndex = contains(args, "--no-color")
	if options.color {
		args = append(args[:argIndex], args[argIndex+1:]...)
	}

	if help, _ := contains(args, "-h", "--help"); help || len(args) < 1 {
		fmt.Println(`Investigo - Investigate User Across Social Networks.`)
		os.Exit(0)
	}

	initializeSiteData()

	for _, username := range args {
		fmt.Println(username)
	}
}
