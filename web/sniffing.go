package web

import (
    //"os"
    "fmt"
    "strings"
    //"io/ioutil"
    "encoding/json"
)

const jsonData string = `
{
    "Github": "https://github.com/?",
    "WordPress": "https://?.wordpress.com",
    "NAVER": "https://blog.naver.com/?",
    "DAUM": "http://blog.daum.net/?",
    "Pinterest": "https://www.pinterest.com/?",
    "Instagram": "https://www.instagram.com/?",
    "Twitter": "https://twitter.com/?",
    "Steam": "https://steamcommunity.com/id/?",
    "YouTube": "https://www.youtube.com/user/?",
    "Reddit": "https://www.reddit.com/user/?",
    "Medium": "https://medium.com/@?",
    "Blogger": "https://?.blogspot.com/",
    "GitLab": "https://gitlab.com/?"
}
`

var webURL map[string]interface{}

func init() {
    /*dataFile, err := os.Open("data.json")
    if err != nil {
        panic(err)
    }
    defer dataFile.Close()

    byteValue, _ := ioutil.ReadAll(dataFile)*/
    json.Unmarshal([]byte(jsonData), &webURL)
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
