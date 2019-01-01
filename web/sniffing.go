package web

import (
	"os"
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
		IsUserExist(webURL[site].(string), username, site)
	}
}
