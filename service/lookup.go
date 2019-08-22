package service

import (
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/lucmski/Investigo/config"
	"github.com/lucmski/Investigo/model"
)

type Investigo struct {
	db *gorm.DB
	// opts *model.Options
}

// Investigo investigate if username exists on social media.
func Lookup(username string, site string, data model.SiteData, options config.Options) model.Result {
	var url, urlProbe string
	result := model.Result{
		Usernane: username,
		URL:      data.URL,
		URLProbe: data.URLProbe,
		Proxied:  options.WithTor,
		Exist:    false,
		Site:     site,
		Err:      true,
		ErrMsg:   "No return value",
	}

	// string to display
	url = strings.Replace(data.URL, "{}", username, 1)

	if data.URLProbe != "" {
		urlProbe = strings.Replace(data.URLProbe, "{}", username, 1)
	} else {
		urlProbe = url
	}

	r, err := Request(urlProbe, options)

	if err != nil {
		if r != nil {
			r.Body.Close()
		}
		return model.Result{
			Usernane: username,
			URL:      data.URL,
			URLProbe: data.URLProbe,
			Proxied:  options.WithTor,
			Exist:    false,
			Site:     site,
			Err:      true,
			ErrMsg:   err.Error(),
		}
	}

	// check error types
	switch data.ErrorType {
	case "status_code":
		if r.StatusCode <= 300 || r.StatusCode < 200 {
			result = model.Result{
				Usernane: username,
				URL:      data.URL,
				URLProbe: data.URLProbe,
				Proxied:  options.WithTor,
				Exist:    true,
				Link:     url,
				Site:     site,
			}
		} else {
			result = model.Result{
				Site:     site,
				Usernane: username,
			}
		}
	case "message":
		if !strings.Contains(ReadResponseBody(r), data.ErrorMsg) {
			result = model.Result{
				Usernane: username,
				URL:      data.URL,
				URLProbe: data.URLProbe,
				Proxied:  options.WithTor,
				Exist:    true,
				Link:     url,
				Site:     site,
			}
		} else {
			// check if 404
			result = model.Result{
				URL:      data.URL,
				URLProbe: data.URLProbe,
				Proxied:  options.WithTor,
				Usernane: username,
				Site:     site,
			}
		}
	case "response_url":
		// In the original Sherlock implementation,
		// the error type `response_url` works as `status_code`.
		if (r.StatusCode <= 300 || r.StatusCode < 200) && r.Request.URL.String() == url {
			result = model.Result{
				Usernane: username,
				URL:      data.URL,
				URLProbe: data.URLProbe,
				Proxied:  options.WithTor,
				Exist:    true,
				Link:     url,
				Site:     site,
			}
		} else {
			result = model.Result{
				Usernane: username,
				URL:      data.URL,
				URLProbe: data.URLProbe,
				Proxied:  options.WithTor,
				Site:     site,
			}
		}
	default:
		result = model.Result{
			Usernane: username,
			Proxied:  options.WithTor,
			Exist:    false,
			Err:      true,
			ErrMsg:   "Unsupported error type `" + data.ErrorType + "`",
			Site:     site,
		}
	}

	r.Body.Close()

	return result
}
