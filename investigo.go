package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	color "github.com/fatih/color"
	"golang.org/x/net/proxy"
)

const (
	dataFileName  string = "data.json"
	userAgent     string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36"
	maxGoroutines int    = 64
)

// Initialize sites not included in Sherlock
func initializeExtraSiteData() {
	siteData["Pornhub"] = SiteData{
		ErrorType: "status_code",
		URLMain:   "https://www.pornhub.com/",
		URL:       "https://www.pornhub.com/users/{}",
	}
	siteData["NAVER"] = SiteData{
		ErrorType: "status_code",
		URLMain:   "https://www.naver.com/",
		URL:       "https://blog.naver.com/{}",
	}
	siteData["xvideos"] = SiteData{
		ErrorType: "status_code",
		URLMain:   "https://xvideos.com/",
		URL:       "https://xvideos.com/profiles/{}",
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

	options.withTor, argIndex = HasElement(args, "-t", "--tor")
	if options.withTor {
		args = append(args[:argIndex], args[argIndex+1:]...)
	}

	options.verbose, argIndex = HasElement(args, "-v", "--verbose")
	if options.verbose {
		args = append(args[:argIndex], args[argIndex+1:]...)
	}

	options.checkForUpdate, argIndex = HasElement(args, "--update")
	if options.checkForUpdate {
		args = append(args[:argIndex], args[argIndex+1:]...)
	}

	// Loads site data from sherlock database and assign to a variable.
	initializeSiteData(options.checkForUpdate)

	if help, _ := HasElement(args, "-h", "--help"); help || len(args) < 1 {
		os.Exit(0)
	}

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
				res := Investigo(username, site, siteData[site])
				WriteResult(res)
				<-guard
			}(site)
		}
		waitGroup.Wait()
	}
	return
}

// Result of Investigo function
type Result struct {
	Usernane string
	Exist    bool
	Proxied  bool
	Site     string
	URL      string
	URLProbe string
	Link     string
	Err      bool
	ErrMsg   string
}

var (
	guard     = make(chan int, maxGoroutines)
	waitGroup = &sync.WaitGroup{}
	logger    = log.New(color.Output, "", 0)
	siteData  = map[string]SiteData{}
	options   struct {
		noColor         bool
		updateBeforeRun bool
		withTor         bool
		withAdmin       bool
		withExport      bool
		verbose         bool
		checkForUpdate  bool
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

func initializeSiteData(forceUpdate bool) {
	jsonFile, err := os.Open(dataFileName)
	if err != nil || forceUpdate {
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

		if forceUpdate {
			jsonFile.Close()
		}

		r, err := Request("https://raw.githubusercontent.com/sherlock-project/sherlock/master/data.json")
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

// Specify Tor proxy ip and port
// var torProxy string = "socks5://127.0.0.1:9050" // 9150 w/ Tor Browser
// var UseTor bool = true

// Request makes HTTP request
func Request(target string) (*http.Response, RequestError) {
	request, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", userAgent)

	client := &http.Client{}
	// client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
	//	return errors.New("Redirect")
	// }

	if options.withTor {

		// fmt.Println("using tor... ")
		tbProxyURL, err := url.Parse("socks5://127.0.0.1:9050")
		if err != nil {
			return nil, err
		}
		tbDialer, err := proxy.FromURL(tbProxyURL, proxy.Direct)
		if err != nil {
			return nil, err
		}
		tbTransport := &http.Transport{
			Dial: tbDialer.Dial,
		}
		client.Transport = tbTransport
	}

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
		Usernane: username,
		URL:      data.URL,
		URLProbe: data.URLProbe,
		Proxied:  options.withTor,
		Exist:    false,
		Site:     site,
		Err:      true,
		ErrMsg:   "No return value",
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
			Usernane: username,
			URL:      data.URL,
			URLProbe: data.URLProbe,
			Proxied:  options.withTor,
			Exist:    false,
			Site:     site,
			Err:      true,
			ErrMsg:   err.Error(),
		}
	}

	// check error types
	switch data.ErrorType {
	case "status_code":
		if r.StatusCode <= 300 || r.StatusCode < 200 {
			result = Result{
				Usernane: username,
				URL:      data.URL,
				URLProbe: data.URLProbe,
				Proxied:  options.withTor,
				Exist:    true,
				Link:     url,
				Site:     site,
			}
		} else {
			result = Result{
				Site:     site,
				Usernane: username,
			}
		}
	case "message":
		if !strings.Contains(ReadResponseBody(r), data.ErrorMsg) {
			result = Result{
				Usernane: username,
				URL:      data.URL,
				URLProbe: data.URLProbe,
				Proxied:  options.withTor,
				Exist:    true,
				Link:     url,
				Site:     site,
			}
		} else {
			// check if 404
			result = Result{
				URL:      data.URL,
				URLProbe: data.URLProbe,
				Proxied:  options.withTor,
				Usernane: username,
				Site:     site,
			}
		}
	case "response_url":
		// In the original Sherlock implementation,
		// the error type `response_url` works as `status_code`.
		if (r.StatusCode <= 300 || r.StatusCode < 200) && r.Request.URL.String() == url {
			result = Result{
				Usernane: username,
				URL:      data.URL,
				URLProbe: data.URLProbe,
				Proxied:  options.withTor,
				Exist:    true,
				Link:     url,
				Site:     site,
			}
		} else {
			result = Result{
				Usernane: username,
				URL:      data.URL,
				URLProbe: data.URLProbe,
				Proxied:  options.withTor,
				Site:     site,
			}
		}
	default:
		result = Result{
			Usernane: username,
			Proxied:  options.withTor,
			Exist:    false,
			Err:      true,
			ErrMsg:   "Unsupported error type `" + data.ErrorType + "`",
			Site:     site,
		}
	}

	r.Body.Close()

	return result
}

// Check content of

// WriteResult writes investigation result to stdout and file
func WriteResult(result Result) {
	if options.noColor {
		if result.Exist {
			logger.Printf("[%s] %s: %s\n", ("+"), result.Site, result.Link)
		} else {
			if result.Err {
				logger.Printf("[%s] %s: ERROR: %s", ("!"), result.Site, (result.ErrMsg))
			} else if options.verbose {
				logger.Printf("[%s] %s: %s", ("-"), result.Site, ("Not Found!"))
			}
		}
	} else {
		if result.Exist {
			logger.Printf("[%s] %s: %s\n", color.HiGreenString("+"), color.HiWhiteString(result.Site), result.Link)
		} else {
			if result.Err {
				logger.Printf("[%s] %s: %s: %s", color.HiRedString("!"), result.Site, color.HiMagentaString("ERROR"), color.HiRedString(result.ErrMsg))
			} else if options.verbose {
				logger.Printf("[%s] %s: %s", color.HiRedString("-"), result.Site, color.HiYellowString("Not Found!"))
			}
		}
	}

	return
}

func showBanner() {
	banner := `
                                              ."""-.
                                             /      \
 ____  _               _            _        |  _..--'-.
/ ___|| |__   ___ _ __| | ___   ___| |__    >.` + "`" + `__.-""\;"` + "`" + `
\___ \| '_ \ / _ \ '__| |/ _ \ / __| |/ /   / /(     ^\
 ___) | | | |  __/ |  | | (_) | (__|   <    '-` + "`" + `)     =|-.
|____/|_| |_|\___|_|  |_|\___/ \___|_|\_\    /` + "`" + `--.'--'   \ .-.
                                           .'` + "`" + `-._ ` + "`" + `.\    | J /`

	fmt.Printf("%v\n\n", banner)
}
