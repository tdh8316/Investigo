package web

import (
	"fmt"
	"strings"
	"net/http"
)

// IsUserExist check is user that use `username` as a id exist
func IsUserExist(url string, username string, site string) {
	url = strings.Replace(url, "?", username, 1)
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	if response.StatusCode == 404 || response.Request.URL.String() != url {
		fmt.Printf("[-] %s: Not found username `%s`\n", site, username)
	} else if response.StatusCode == 200 {
		fmt.Printf("[+] %s: %s\n", site, url)
	}
}