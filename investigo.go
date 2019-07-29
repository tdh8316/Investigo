package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

import (
	color "github.com/fatih/color"
)

const (
	dataFileName  string = "data.json"
	userAgent     string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36"
	maxGoroutines int    = 64
)

// Result of Investigo function
type Result struct {
	exist  bool
	site   string
	link   string
	err    bool
	errMsg string
}

var (
	guard     = make(chan int, maxGoroutines)
	waitGroup = &sync.WaitGroup{}
	logger    = log.New(color.Output, "", 0)
	siteData  = map[string]SiteData{}
	options   struct {
		noColor         bool
		updateBeforeRun bool
		verbose         bool
	}
)

// A SiteData struct for json datatype
type SiteData struct {
	ErrorType string `json:"errorType"`
	ErrorMsg  string `json:"errorMsg"`
	URL       string `json:"url"`
	URLMain   string `json:"urlMain"`
	URLProbe  string `json:"urlProbe"`
	URLError  string `json:"errorUrl"`
	// UsedUsername   string `json:"username_claimed"`
	// UnusedUsername string `json:"username_unclaimed"`
	// RegexCheck string `json:"regexCheck"`
	// Rank int`json:"rank"`
}

// RequestError interface
type RequestError interface {
	Error() string
}

func initializeSiteData() {
	jsonFile, err := os.Open(dataFileName)
	if err != nil {
		if options.noColor {
			fmt.Printf(
				"%s Failed to read %s from current directory. %s",
				("->"),
				dataFileName,
				("Downloading..."),
			)
		} else {
			fmt.Fprintf(
				color.Output,
				"%s Failed to read %s from current directory. %s",
				color.HiRedString("->"),
				dataFileName,
				color.HiYellowString("Downloading..."),
			)
		}

		r, err := Request("https://raw.githubusercontent.com/tdh8316/Investigo/master/data.json")
		if err != nil || r.StatusCode != 200 {
			if options.noColor {
				fmt.Printf(" [%s]\n", ("Failed"))
			} else {
				fmt.Fprintf(color.Output, " [%s]\n", color.HiRedString("Failed"))
			}
			panic("Failed to connect to Investigo repository.")
		} else {
			defer r.Body.Close()
		}
		if _, err := os.Stat(dataFileName); !os.IsNotExist(err) {
			if err = os.Remove(dataFileName); err != nil {
				panic(err)
			}
		}
		_updateFile, _ := os.OpenFile(dataFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if _, err := _updateFile.WriteString(ReadResponseBody(r)); err != nil {
			if options.noColor {
				fmt.Printf("Failed to update data.\n")
			} else {
				fmt.Fprintf(color.Output, color.RedString("Failed to update data.\n"))
			}
			panic(err)
		}

		_updateFile.Close()
		jsonFile, _ = os.Open(dataFileName)

		fmt.Println(" [Done]")
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic("Error while read " + dataFileName)
	} else {
		json.Unmarshal([]byte(byteValue), &siteData)
	}
	return
}

// Initialize sites banned from Sherlock
func initializeExtraSiteData() {
	siteData["Pornhub"] = SiteData{
		ErrorType: "status_code",
		URLMain: "https://www.pornhub.com/",
		URL: "https://www.pornhub.com/users/{}",
	}
}

func main() {
	fmt.Println(`Investigo - Investigate User Across Social Networks.`)

	args := os.Args[1:]
	var argIndex int

	options.noColor, argIndex = HasElement(args, "--no-color")
	if options.noColor {
		logger = log.New(os.Stdout, "", 0)
		args = append(args[:argIndex], args[argIndex+1:]...)
	}

	options.verbose, argIndex = HasElement(args, "-v", "--verbose")
	if options.verbose {
		args = append(args[:argIndex], args[argIndex+1:]...)
	}

	if help, _ := HasElement(args, "-h", "--help"); help || len(args) < 1 {
		os.Exit(0)
	}

	// Loads site data from sherlock database and assign to a variable.
	initializeSiteData()

	// Loads extra site data
	initializeExtraSiteData()

	for _, username := range args {
		if options.noColor {
			fmt.Printf("Investigating %s on:\n", username)
		} else {
			fmt.Fprintf(color.Output, "Investigating %s on:\n", color.HiGreenString(username))
		}
		waitGroup.Add(len(siteData))
		for site := range siteData {
			guard <- 1
			go func(site string) {
				defer waitGroup.Done()
				WriteResult(
					Investigo(username, site, siteData[site]),
				)
				<-guard
			}(site)
		}
		waitGroup.Wait()
	}

	return
}

// Request makes HTTP request
func Request(url string) (*http.Response, RequestError) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", userAgent)
	client := &http.Client{}

	return client.Do(request)
}

// ReadResponseBody reads response body and return string
func ReadResponseBody(response *http.Response) string {
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	return string(bodyBytes)
}

// HasElement reports whether elements is within array.
func HasElement(array []string, targets ...string) (bool, int) {
	for index, item := range array {
		for _, target := range targets {
			if item == target {
				return true, index
			}
		}
	}
	return false, -1
}

// Investigo investigate if username exists on social media.
func Investigo(username string, site string, data SiteData) Result {
	var url, urlProbe string
	result := Result{
		exist:  false,
		site:   site,
		err:    true,
		errMsg: "No return value",
	}

	// string to display
	url = strings.Replace(data.URL, "{}", username, 1)

	if data.URLProbe != "" {
		urlProbe = strings.Replace(data.URLProbe, "{}", username, 1)
	} else {
		urlProbe = url
	}

	r, err := Request(urlProbe)

	if err != nil {
		if r != nil {
			r.Body.Close()
		}
		return Result{
			exist: false, site: site, err: true, errMsg: err.Error(),
		}
	}

	switch data.ErrorType {
	case "status_code":
		if r.StatusCode <= 300 || r.StatusCode < 200 {
			result = Result{
				exist: true, link: url, site: site,
			}
		} else {
			result = Result{site: site}
		}
	case "message":
		if !strings.Contains(ReadResponseBody(r), data.ErrorMsg) {
			result = Result{
				exist: true, link: url, site: site,
			}
		} else {
			result = Result{site: site}
		}
	case "response_url":
		// In the original Sherlock implementation,
		// the error type `response_url` works as `status_code`.
		if (r.StatusCode <= 300 || r.StatusCode < 200) && r.Request.URL.String() == url {
			result = Result{
				exist: true, link: url, site: site,
			}
		} else {
			result = Result{site: site}
		}
	default:
		result = Result{
			exist: false, err: true, errMsg: "Unsupported error type `" + data.ErrorType + "`", site: site,
		}
	}

	r.Body.Close()

	return result
}

// WriteResult writes investigation result to stdout and file
func WriteResult(result Result) {
	if options.noColor {
		if result.exist {
			logger.Printf("[%s] %s: %s\n", ("+"), result.site, result.link)
		} else {
			if result.err {
				logger.Printf("[%s] %s: ERROR: %s", ("!"), result.site, (result.errMsg))
			} else if options.verbose {
				logger.Printf("[%s] %s: %s", ("-"), result.site, ("Not Found!"))
			}
		}
	} else {
		if result.exist {
			logger.Printf("[%s] %s: %s\n", color.HiGreenString("+"), color.HiWhiteString(result.site), result.link)
		} else {
			if result.err {
				logger.Printf("[%s] %s: %s: %s", color.HiRedString("!"), result.site, color.HiMagentaString("ERROR"), color.HiRedString(result.errMsg))
			} else if options.verbose {
				logger.Printf("[%s] %s: %s", color.HiRedString("-"), result.site, color.HiYellowString("Not Found!"))
			}
		}
	}

	return
}
