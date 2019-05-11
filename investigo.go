package main

import (
	"encoding/json"
	"fmt"
	color "github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type error interface {
    Error() string
}

var sns = map[string]string{}
var snsCaseLower = map[string]string{}

var userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36"

func readPageSource(response *http.Response) string {
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	return string(bodyBytes)
}

func httpRequest(url string) (
	response *http.Response, respondedURL string, err error) {
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("User-Agent",
		userAgent)
	client := &http.Client{}
	response, clientError := client.Do(request)
	if clientError == nil {
		respondedURL = response.Request.URL.String()
		err = nil
	} else {
		respondedURL = ""
		err = clientError
	}

	return
}

// Check if the username exists on snsName
func isUserExist(snsName string, username string, caseLower bool) bool {
	url := sns[snsName]
	if caseLower {
		url = snsCaseLower[strings.ToLower(snsName)]
	}

	response, respondedURL, err := httpRequest(strings.Replace(url, "?", username, 1))
	if err != nil {
		fmt.Fprintf(color.Output, color.HiYellowString("Failed to make a connection to %s\n"), snsName)
		// fmt.Println(err)
		log, _ := os.OpenFile("http-request-exception.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		defer log.Close()
		log.WriteString(err.Error() + "\n")
		return false
	}

	snsName = strings.ToLower(snsName)

	switch snsName {
	case "wordpress":
		if respondedURL == url {
			return true
		}
		return false
	case "steam":
		if !strings.Contains(
			readPageSource(response),
			"The specified profile could not be found.") {
			return true
		}
		return false
	case "pinterest":
		if url == respondedURL || strings.Contains(respondedURL, username) {
			return true
		}
		return false
	case "gitlab":
		if url == respondedURL {
			return true
		}
		return false
	case "egloos":
		if !strings.Contains(
			readPageSource(response),
			"블로그가 존재하지 않습니다") {
			return true
		}
		return false
	}

	if response.StatusCode == 200 {
		return true
	}
	return false
}

// Array, element
func contains(array []string, str string) (bool, int) {
	for index, item := range array {
		if item == str {
			return true, index
		}
	}
	return false, 0 // Index
}

func loadSNSList() {
	jsonFile, err := os.Open("sites.json")

	if err != nil {
		udpateSNSList()
		jsonFile, _ = os.Open("sites.json")
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var snsInterface map[string]interface{}
	json.Unmarshal([]byte(byteValue), &snsInterface)
	// Json to map
	for k, v := range snsInterface {
		sns[k] = v.(string)
	}
}

func udpateSNSList() {
	fmt.Printf("Update investigo... ")
	response, _, err := httpRequest("https://raw.githubusercontent.com/tdh8316/Investigo/master/sites.json1")
	if err != nil || response.StatusCode == 404 {
		panic("Failed to connect to Investigo repository.")
	}
	jsonData := readPageSource(response)

	fileName := "sites.json"
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		if err = os.Remove(fileName); err != nil {
			panic(err)
		}
	}
	dataFile, _ := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	defer dataFile.Close()

	if _, err := dataFile.WriteString(jsonData); err != nil {
		fmt.Fprintf(color.Output, color.RedString("Failed to update data\n"))
	}
	fmt.Println("Done.")
}

func printHelp() {
	fmt.Println("Investigo - Investigate User Across Social Networks.")
	fmt.Println("\nUsage: go run investigo.go [-h] [--no-color] [--verbose] [--update] [--site SITE_NAME] USERNAMES\n" +
		"\n positional arguments:\n\tUSERNAMES\t   Usernames to investigate")

	fmt.Println(`
 optional arguments:
	--site             Specify sites to search. (e.g --site github)
	--disable-color
	--verbose          Print sites USERNAME is not found.
	--update           Update data automatically`)
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		printHelp()
	}
	disableColor, _ := contains(args, "--no-color")
	disableQuiet, _ := contains(args, "--verbose")
	updateData, _ := contains(args, "--update")
	printHelpAndExit, _ := contains(args, "--help")
	if !printHelpAndExit {
		printHelpAndExit, _ = contains(args, "-h")
	}
	specificSite, siteIndex := contains(args, "--site")
	specifiedSite := ""
	if specificSite {
		specifiedSite = args[siteIndex+1]
	}

	useCustomUserAgent, argIndex := contains(args, "--user-agent")
	if useCustomUserAgent {
		userAgent = args[argIndex+1]
	}

	if updateData {
		udpateSNSList()
	}

	if printHelpAndExit {
		printHelp()
		os.Exit(0)
	}

	loadSNSList()

	for _, username := range args {
		if isOpt, _ :=
			contains([]string{"--no-color", "--verbose", specifiedSite, "--site", "--update", "--user-agent"}, username); isOpt {
			continue
		}
		if disableColor {
			fmt.Printf("Searching username %s\n", username)
		} else {
			fmt.Fprintf(color.Output, "%s %s\n", color.HiMagentaString("Searching username"), username)
		}
		if specificSite {
			// Case ignore
			for k, v := range sns {
				snsCaseLower[strings.ToLower(k)] = v
			}

			specifiedSite = strings.ToLower(specifiedSite)

			if _, isExist := snsCaseLower[specifiedSite]; !isExist {
				if disableColor {
					fmt.Println("Unknown site: " + specifiedSite)
				} else {
					fmt.Fprintf(color.Output, "%s: %s", color.RedString("Unknown site"), specifiedSite)
				}
				break
			}

			if isUserExist(specifiedSite, username, true) {
				if disableColor {
					fmt.Printf(
						"[+] %s: %s\n", specifiedSite, strings.Replace(
							snsCaseLower[strings.ToLower(specifiedSite)], "?", username, 1))
				} else {
					fmt.Fprintf(color.Output,
						"[%s] %s: %s\n",
						color.HiGreenString("+"), color.HiWhiteString(specifiedSite),
						color.WhiteString(
							strings.Replace(snsCaseLower[strings.ToLower(specifiedSite)],
								"?", username, 1)))
				}
			} else {
				if disableColor {
					fmt.Printf(
						"[-] %s: Not found!\n", specifiedSite)
				} else {
					fmt.Fprintf(color.Output,
						"[%s] %s: %s\n", color.HiRedString("-"),
						color.HiWhiteString(specifiedSite),
						color.HiYellowString("Not found!"))
				}
			}
			break
		}

		fileName := username + ".txt"
		if _, err := os.Stat(fileName); !os.IsNotExist(err) {
			if err = os.Remove(fileName); err != nil {
				panic(err)
			}
		}
		resFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
		defer resFile.Close()

		for site := range sns {
			if isUserExist(site, username, false) {
				if disableColor {
					fmt.Printf(
						"[+] %s: %s\n", site, strings.Replace(sns[site], "?", username, 1))
				} else {
					fmt.Fprintf(color.Output,
						"[%s] %s: %s\n",
						color.HiGreenString("+"), color.HiWhiteString(site),
						color.WhiteString(strings.Replace(sns[site], "?", username, 1)))
				}

				if _, err = resFile.WriteString(site + ": " + strings.Replace(sns[site], "?", username, 1) + "\n"); err != nil {
					panic(err)
				}

			} else {
				if !disableQuiet {
					continue
				}

				if disableColor {
					fmt.Printf("[-] %s: Not found!\n", site)
				} else {
					fmt.Fprintf(color.Output,
						"[%s] %s: %s\n",
						color.HiRedString("-"), color.HiWhiteString(site),
						color.HiYellowString("Not found!"))
				}
			}
		}
		fmt.Println("\nYour search results have been saved to " + fileName)
	}
}
