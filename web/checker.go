package web

import (
    "strings"
    "net/http"
    "io/ioutil"
)

func getPageSource(url string, response *http.Response) string {
    bodyBytes, err := ioutil.ReadAll(response.Body)
    if err != nil { panic(err) }
    return string(bodyBytes)
}

// IsUserExist check is user that use `username` as a id exist
func IsUserExist(url string, username string, site string) bool {
    site = strings.ToLower(site)
    request, err := http.NewRequest("GET", url, nil)
    if err != nil { panic(err) }
    request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
    client := &http.Client{}
    response, err := client.Do(request)
    respondedURL := response.Request.URL.String()
    if err != nil { panic(err) }
    defer response.Body.Close()

    if site == "wordpress" {
        if respondedURL == url { return true }
        return false
    }

    if site == "steam" {
        bodyString := getPageSource(url, response)
        if !strings.Contains(bodyString,
            "The specified profile could not be found.") { 
                return true }
        return false
    }

    if site == "pinterest" {
        if url == respondedURL { return true }
        if strings.Contains(respondedURL, username) { return true }
        return false
    }

    if site == "gitlab" {
        if url == respondedURL { return true }
        return false
    }

    if response.StatusCode == 200 { return true }

    return false
}
