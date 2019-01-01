package web

import (
	"fmt"
	"strings"
	"net/http"
)

// IsUserExist check is user that use `username` as a id exist
func IsUserExist(url string, username string, site string) {
	url = strings.Replace(url, "?", username, 1)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	// Set useragent
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	respondedURL := response.Request.URL.String()

	if (response.StatusCode == 404) || (respondedURL != url && !strings.Contains(respondedURL, username)) || (
		respondedURL != url && strings.ToLower(site) == "wordpress") {
		fmt.Printf("[-] %s: Not found\n", site)
	} else if response.StatusCode == 200 {
		fmt.Printf("[+] %s: %s\n", site, url)
	} else {
		fmt.Printf("[-] %s: UNKNOWN RESPONSE %d", site, response.StatusCode)
	}
}
