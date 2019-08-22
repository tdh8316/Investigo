package service

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/kamilsk/breaker"
	"github.com/kamilsk/retry"
	"github.com/kamilsk/retry/strategy"
	"golang.org/x/net/proxy"

	"github.com/lucmski/Investigo/config"
	"github.com/lucmski/Investigo/model"
)

// Specify Tor proxy ip and port
// var torProxy string = "socks5://127.0.0.1:9050" // 9150 w/ Tor Browser
// var UseTor bool = true

// Request makes HTTP request
func Request(target string, options config.Options) (*http.Response, model.RequestError) {
	// Add tor proxy

	var response *http.Response
	request, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", config.UserAgent)

	client := &http.Client{}
	// client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
	//	return errors.New("Redirect")
	// }

	if options.WithTor {

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
