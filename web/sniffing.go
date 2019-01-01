package web

import (
	"os"
	"fmt"
	"strings"
	"io/ioutil"
	"encoding/json"
)

var webURL map[string]interface{}

func init() {
	dataFile, err := os.Open("data.json")
	if err != nil {
		panic(err)
	}
	defer dataFile.Close()

	byteValue, _ := ioutil.ReadAll(dataFile)
	json.Unmarshal([]byte(byteValue), &webURL)
}


// Sniffer search username across social media
func Sniffer(username string) {
	for site := range webURL {
		url := strings.Replace(webURL[site].(string), "?", username, 1)
		
		if IsUserExist(url, username, site) {
			fmt.Printf("[+] %s: %s\n", site, url)
		} else {
			fmt.Printf("[-] %s: Not found\n", site)
		}
	}
}
