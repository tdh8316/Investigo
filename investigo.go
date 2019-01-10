package main

import (
	"os"
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
	color "github.com/fatih/color"
)


var sns = map[string]string {
    "Github": "https://github.com/?",
    "WordPress": "https://?.wordpress.com",
    "NAVER": "https://blog.naver.com/?",
    "DAUM Blog": "http://blog.daum.net/?",
    "Tistory": "https://?.tistory.com/",
    "Egloos": "http://?.egloos.com/",
    "Pinterest": "https://www.pinterest.com/?",
    "Instagram": "https://www.instagram.com/?",
    "Twitter": "https://twitter.com/?",
    "Steam": "https://steamcommunity.com/id/?",
    "YouTube": "https://www.youtube.com/user/?",
    "Reddit": "https://www.reddit.com/user/?",
    "Medium": "https://medium.com/@?",
    "Blogger": "https://?.blogspot.com/",
    "GitLab": "https://gitlab.com/?",
    "Google Plus": "https://plus.google.com/+?",
    "About.me": "https://about.me/?",
    "SlideShare": "https://slideshare.net/?",
    "BitBucket": "https://bitbucket.org/?",
    "Quora": "https://www.quora.com/profile/?",
    "SourceForge": "https://sourceforge.net/u/?",
    "Wix": "https://?.wix.com",
    "SoundCloud": "https://soundcloud.com/?",
    "Facebook": "https://www.facebook.com/?",
    "Disqus": "https://disqus.com/?",
    "DeviantART": "https://www.deviantart.com/?",
    "Spotify": "https://open.spotify.com/user/?",
    "Patreon": "https://www.patreon.com/?",
    "DailyMotion": "https://www.dailymotion.com/?",
    "Slack": "https://?.slack.com",
    "Kaggle": "https://www.kaggle.com/?",
    "Itch.io": "https://?.itch.io/",

    "Zhihu": "https://www.zhihu.com/people/?",
    "Gitee": "https://gitee.com/?",
}


func getPageSource(response *http.Response) string {
    bodyBytes, err := ioutil.ReadAll(response.Body)
    if err != nil {
        panic(err)
    }
    return string(bodyBytes)
}


func httpRequest(url string) (
        response *http.Response, respondedURL string) {
    request, _ := http.NewRequest("GET", url, nil)
    request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        panic(err)
    }
    respondedURL = response.Request.URL.String()
    
    return
}


func isUserExist(snsName string, username string) bool {
    url := sns[snsName]
    response, respondedURL := httpRequest(strings.Replace(url, "?", username, 1))
    snsName = strings.ToLower(snsName)

    //TODO: Turn if into case
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
    } else if snsName == "egloos" {
        if !strings.Contains(
            getPageSource(response),
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


func contains(array []string, str string) (bool, int) {
    for index, item := range array {
       if item == str {
          return true, index
       }
    }
    return false, 0
 }


func main() {
    // TODO: Add argument: --site
    args := os.Args[1:]
    disableColor, _ := contains(args, "--no-color")
    disableQuiet, _ := contains(args, "--verbose")
    specificSite, siteIndex := contains(args, "--site")
    specifiedSite := args[siteIndex + 1]

    for _, username := range args {
        if isOpt, _ := contains([]string{"--no-color", "--verbose", specifiedSite, "--site"}, username); isOpt {
            continue
        }
        fmt.Fprintf(color.Output, "%s %s on:\n", color.HiMagentaString("Searching username"), username)
        if specificSite {
            if isUserExist(specifiedSite, username) {
                if disableColor {
                    fmt.Printf(
                        "[+] %s: %s\n", specifiedSite, strings.Replace(sns[specifiedSite], "?", username, 1))
                } else {
                    fmt.Fprintf(color.Output,
                        "[%s] %s: %s\n",
                        color.HiGreenString("+"), color.HiWhiteString(specifiedSite),
                        color.WhiteString(strings.Replace(sns[specifiedSite], "?", username, 1)))
                }
            } else {
                if disableColor {
                    fmt.Printf(
                        "[-] %s: Not found!\n", specifiedSite)
                } else {
                    fmt.Fprintf(color.Output,
                        "[%s] %s: %s\n",
                        color.HiRedString("-"), color.HiWhiteString(specifiedSite),
                        color.HiYellowString("Not found!"))
                }
            }
            break
        }

        fileName := "./" + username + ".txt"
        if _, err := os.Stat(fileName); !os.IsNotExist(err) {
            if err = os.Remove(fileName); err != nil {
                panic(err)
            }
        }
        resFile, err := os.OpenFile(fileName, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0600)
        if err != nil {
            panic(err)
        }
        defer resFile.Close()

        for site := range sns {
            if isUserExist(site, username) {
                if disableColor {
                    fmt.Printf(
                        "[+] %s: %s\n", site, strings.Replace(sns[site], "?", username, 1))
                } else {
                    fmt.Fprintf(color.Output,
                        "[%s] %s: %s\n",
                        color.HiGreenString("+"), color.HiWhiteString(site),
                        color.WhiteString(strings.Replace(sns[site], "?", username, 1)))
                }
                if _, err = resFile.WriteString(strings.Replace(sns[site], "?", username, 1) + "\n");
                err != nil {
                    panic(err)
                }
            } else {
                if !disableQuiet {
                    continue
                }
                
                if disableColor {
                    fmt.Printf(
                        "[-] %s: Not found!\n", site)
                } else {
                    fmt.Fprintf(color.Output,
                        "[%s] %s: %s\n",
                        color.HiRedString("-"), color.HiWhiteString(site),
                        color.HiYellowString("Not found!"))
                }
            }
        }
    }
}
