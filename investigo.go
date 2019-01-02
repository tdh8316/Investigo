package main

import "os"
import "fmt"
import "strings"
import "net/http"
import "io/ioutil"


var sns = map[string]string {
    "Github": "https://github.com/?",
    "WordPress": "https://?.wordpress.com",
    "NAVER": "https://blog.naver.com/?",
    "DAUM Blog": "http://blog.daum.net/?",
    "Pinterest": "https://www.pinterest.com/?",
    "Instagram": "https://www.instagram.com/?",
    "Twitter": "https://twitter.com/?",
    "Steam": "https://steamcommunity.com/id/?",
    "YouTube": "https://www.youtube.com/user/?",
    "Reddit": "https://www.reddit.com/user/?",
    "Medium": "https://medium.com/@?",
    "Blogger": "https://?.blogspot.com/",
    "GitLab": "https://gitlab.com/?",
}


func getPageSource(response *http.Response) string {
    bodyBytes, err := ioutil.ReadAll(response.Body)
    if err != nil {
        panic(err)
    }
    return string(bodyBytes)
}


func request(url string) (*http.Response, string) {
    request, _ := http.NewRequest("GET", url, nil)
    request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        panic(err)
    }
    respondedURL := response.Request.URL.String()
    
    return response, respondedURL
}


func isUserExist(snsName string, username string) bool {
    url := sns[snsName]
    response, respondedURL := request(strings.Replace(url, "?", username, 1))
    snsName = strings.ToLower(snsName)

    if snsName == "wordpress" {
        if respondedURL == url {
            return true
        }
        return false
    } else if snsName == "steam" {
        if !strings.Contains(
            getPageSource(response),
            "The specified profile could not be found.") { 
                return true
            }
        return false
    } else if snsName == "pinterest" {
        if url == respondedURL || strings.Contains(respondedURL, username) {
            return true
        }
        return false
    } else if snsName == "gitlab" {
        if url == respondedURL {
            return true
        }
        return false
    }

    if response.StatusCode == 200 {
        return true
    }
    return false
}


func main() {
    for _, username := range os.Args[1:] {
        fmt.Printf("Searching username %s on\n", username)
        for site := range sns {
            if isUserExist(site, username) {
                fmt.Printf("[+] %s: %s\n", site, strings.Replace(sns[site], "?", username, 1))
            } else {
                fmt.Printf("[-] %s: Not found!\n", site)
            }
        }
    }
}
