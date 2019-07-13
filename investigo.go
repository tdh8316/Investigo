package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

import (
	array "./array"
	http "./http"
	color "github.com/fatih/color"
)

const (
	dataFileName string = "data.json"
)

var siteData = map[string]SiteData{}

// A SiteData struct for json datatype
type SiteData struct {
	ErrorType string `json:"errorType"`
	ErrorMsg  string `json:"errorMsg"`
	// Rank int`json:"rank"`
	URL      string `json:"url"`
	URLMain  string `json:"urlMain"`
	URLProbe string `json:"urlProbe"`
	// UsedUsername   string `json:"username_claimed"`
	// UnusedUsername string `json:"username_unclaimed"`
	// RegexCheck string `json:"regexCheck"`
}

// Result of Investigo function
type Result struct {
	exist   bool
	link    string
	message string
}

// Options contains command line arguments data object
type Options struct {
	color           bool
	updateBeforeRun bool
	verbose         bool
}

// Command line options (arguments)
var options Options

// Investigo investigate if username exists on social media.
func Investigo(username string, site string, data SiteData) Result {
	var url, urlProbe string

	// string to display
	url = strings.Replace(data.URL, "{}", username, 1)

	if data.URLProbe != "" {
		urlProbe = strings.Replace(data.URLProbe, "{}", username, 1)
	} else {
		urlProbe = strings.Replace(data.URL, "{}", username, 1)
	}

	r, err := http.Request(urlProbe)
	if err != nil {
		panic(err)
	}

	if data.ErrorType == "status_code" {
		if r.StatusCode <= 300 || r.StatusCode < 200 {
			return Result{
				exist: true, link: url,
			}
		}
		return Result{exist: false, message: color.HiYellowString("Not Found!")}
	} else if data.ErrorType == "message" {
		if !strings.Contains(http.ReadResponseBody(r), data.ErrorMsg) {
			return Result{
				exist: true, link: url,
			}
		}
		return Result{exist: false, message: color.HiYellowString("Not Found!")}
	} else if data.ErrorType == "response_url" {

	} else {
		return Result{
			exist: false, message: "ERROR: Unsupported error type",
		}
	}

	return Result{
		exist: false, message: "ERROR: No return value",
	}
}

// Write result
func resultWriter(site string, exist bool, detail string) {
	if exist {
		fmt.Fprintf(
			color.Output,
			"[%s] %s: %s\n",
			color.HiGreenString("+"), color.HiWhiteString(site), detail,
		)
	} else {
		if options.verbose {
			fmt.Fprintf(
				color.Output,
				"[%s] %s: %s\n",
				color.HiRedString("-"), color.HiWhiteString(site), detail,
			)
		}
	}
	return
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
	args := os.Args[1:]
	var argIndex int

	options.color, argIndex = array.Contains(args, "--no-color")
	if options.color {
		args = append(args[:argIndex], args[argIndex+1:]...)
	}

	options.verbose, argIndex = array.Contains(args, "-v", "--verbose")
	if options.verbose {
		args = append(args[:argIndex], args[argIndex+1:]...)
	}

	if help, _ := array.Contains(args, "-h", "--help"); help || len(args) < 1 {
		fmt.Println(`Investigo - Investigate User Across Social Networks.`)
		os.Exit(0)
	}

	initializeSiteData()

	for _, username := range args {
		for site := range siteData {
			investigo := Investigo(username, site, siteData[site])
			if investigo.exist {
				resultWriter(site, true, investigo.link)
			} else {
				resultWriter(site, false, color.HiMagentaString(investigo.message))
			}
		}
	}

	return
}
