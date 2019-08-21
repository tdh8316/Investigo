package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/agrison/go-tablib"
	color "github.com/fatih/color"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	"github.com/k0kubun/pp"
	"github.com/kamilsk/breaker"
	"github.com/kamilsk/retry"
	"github.com/kamilsk/retry/strategy"
	_ "github.com/mattn/go-sqlite3"
	"github.com/qor/admin"
	"golang.org/x/net/proxy"
)

const (
	dataFileName  string = "data.json"
	userAgent     string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36"
	maxGoroutines int    = 64
)

var (
	DB    *gorm.DB
	Admin *admin.Admin
	DBook *tablib.Databook
)

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

	options.withExport, argIndex = HasElement(args, "-e", "--export")
	if options.withExport {
		args = append(args[:argIndex], args[argIndex+1:]...)
	}

	options.withAdmin, argIndex = HasElement(args, "-a", "--admin")
	if options.withAdmin {
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

	configor.Load(&Config, "config.yml", "data.yml")
	pp.Println("config: ", Config)

	// Loads site data from sherlock database and assign to a variable.
	initializeSiteData(options.checkForUpdate)

	if help, _ := HasElement(args, "-h", "--help"); help || len(args) < 1 {
		os.Exit(0)
	}

	// Loads extra site data
	initializeExtraSiteData()

	if options.withTor {
		pp.Println("Using tor...")
	}

	if options.withAdmin {
		DB, _ = gorm.Open("sqlite3", "investigo.db")
		DB.AutoMigrate(&Result{})
		// Initalize
		Admin = admin.New(&admin.AdminConfig{DB: DB})

		// Allow to use Admin to manage User, Product
		Admin.AddResource(&Result{})
		// Admin.AddResource(&SiteData{})
	}

	if options.withExport {
		DBook = tablib.NewDatabook()
	}

	for _, username := range args {
		var output *tablib.Dataset
		if options.withExport {
			output = tablib.NewDataset([]string{"Status", "Username", "Site", "Info"})
		}
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
				if options.withExport {
					if res.Exist || res.Err {
						output.AppendValues(res.Exist, username, site, res.ErrMsg)
					}
				}
				<-guard
			}(site)
		}
		waitGroup.Wait()
		DBook.AddSheet(username, output)
	}
	if options.withExport {
		// fmt.Println(DBook.YAML())
		for name := range DBook.Sheets() {
			ods := DBook.Sheet(name).Dataset().Tabular("markdown" /* tablib.TabularMarkdown */)
			fmt.Println(ods)
		}
	}
	if options.withAdmin {
		// initalize an HTTP request multiplexer
		mux := http.NewServeMux()

		// Mount admin interface to mux
		Admin.MountTo("/admin", mux)
		fmt.Println("Listening on: 9000")
		http.ListenAndServe(":9000", mux)
	}

	return
}

// Result of Investigo function
type Result struct {
	gorm.Model `yaml:-`
	Usernane   string
	Exist      bool
	Proxied    bool
	Site       string
	URL        string
	URLProbe   string
	Link       string
	Err        bool
	ErrMsg     string
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
	gorm.Model `yaml:-`
	ErrorType  string `json:"errorType"`
	ErrorMsg   string `json:"errorMsg"`
	URL        string `json:"url"`
	URLMain    string `json:"urlMain"`
	URLProbe   string `json:"urlProbe"`
	URLError   string `json:"errorUrl"`
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

		r, err := Request("https://github.com/sherlock-project/sherlock/blob/master/data.json")
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
}

// Specify Tor proxy ip and port
// var torProxy string = "socks5://127.0.0.1:9050" // 9150 w/ Tor Browser
// var UseTor bool = true

// Request makes HTTP request
func Request(target string) (*http.Response, RequestError) {
	// Add tor proxy

	// https://github.com/kamilsk/retry#retryretry
	// backoff strategy
	var response *http.Response
	/*
		var response *http.Response

		action := func(uint) error {
			var err error
			response, err = http.Get("https://github.com/kamilsk/retry")
			return err
		}

		if err := retry.Retry(breaker.BreakByTimeout(time.Minute), action, strategy.Limit(3)); err != nil {
			if err == retry.Interrupted {
				// timeout exceeded
			}
			// handle error
		}
		// work with response
	*/

	request, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", userAgent)
	client := &http.Client{}

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
		// fmt.Println("tor IP:", getIPAdress(request))
	}
	action := func(uint) error {
		var err error
		// pp.Printf("retry# %d \n", i)
		response, err = client.Do(request)
		return err
	}
	if err := retry.Retry(breaker.BreakByTimeout(time.Second*30), action, strategy.Limit(3)); err != nil {
		if err == retry.Interrupted {
			// timeout exceeded
			return nil, err
		}
		// handle error
		return nil, err
	}
	// pp.Println(request)
	return response, err // client.Do(request)
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
	// pp.Println(r)
	//pp.Println("forward: ", forward)

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
			if options.withAdmin {
				if err := DB.Create(&result).Error; err != nil {
					fmt.Println(err)
					return
				}
			}
		} else {
			if result.Err {
				logger.Printf("[%s] %s: %s: %s", color.HiRedString("!"), result.Site, color.HiMagentaString("ERROR"), color.HiRedString(result.ErrMsg))
				if options.withAdmin {
					if err := DB.Create(&result).Error; err != nil {
						fmt.Println(err)
						return
					}
				}
			} else if options.verbose {
				logger.Printf("[%s] %s: %s", color.HiRedString("-"), result.Site, color.HiYellowString("Not Found!"))
			}
		}
	}

	return
}

// test configor for extra rules
var Config = struct {
	APPName string `default:"investigo"`
	DB      struct {
		Name     string
		User     string `default:"root"`
		Password string `required:"true" env:"DBPassword"`
		Port     uint   `default:"3306"`
	}
	Contacts []struct {
		Name  string
		Email string `required:"true"`
	}
	SiteData []SiteData
}{}

//ipRange - a structure that holds the start and end of a range of ip addresses
type ipRange struct {
	start net.IP
	end   net.IP
}

// inRange - check to see if a given ip address is within a range given
func inRange(r ipRange, ipAddress net.IP) bool {
	// strcmp type byte comparison
	if bytes.Compare(ipAddress, r.start) >= 0 && bytes.Compare(ipAddress, r.end) < 0 {
		return true
	}
	return false
}

var privateRanges = []ipRange{
	{
		start: net.ParseIP("10.0.0.0"),
		end:   net.ParseIP("10.255.255.255"),
	},
	{
		start: net.ParseIP("100.64.0.0"),
		end:   net.ParseIP("100.127.255.255"),
	},
	{
		start: net.ParseIP("172.16.0.0"),
		end:   net.ParseIP("172.31.255.255"),
	},
	{
		start: net.ParseIP("192.0.0.0"),
		end:   net.ParseIP("192.0.0.255"),
	},
	{
		start: net.ParseIP("192.168.0.0"),
		end:   net.ParseIP("192.168.255.255"),
	},
	{
		start: net.ParseIP("198.18.0.0"),
		end:   net.ParseIP("198.19.255.255"),
	},
}

// isPrivateSubnet - check to see if this ip is in a private subnet
func isPrivateSubnet(ipAddress net.IP) bool {
	// my use case is only concerned with ipv4 atm
	if ipCheck := ipAddress.To4(); ipCheck != nil {
		// iterate over all our ranges
		for _, r := range privateRanges {
			// check if this ip is in a private range
			if inRange(r, ipAddress) {
				return true
			}
		}
	}
	return false
}

func getIPAdress(r *http.Request) string {
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := strings.Split(r.Header.Get(h), ",")
		// march from right to left until we get a public address
		// that will be the address right before our proxy.
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(addresses[i])
			// header can contain spaces too, strip those out.
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() || isPrivateSubnet(realIP) {
				// bad address, go to next
				continue
			}
			return ip
		}
	}
	return ""
}

func getMyIp() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				os.Stdout.WriteString(ipnet.IP.String() + "\n")
			}
		}
	}
}
