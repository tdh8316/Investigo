package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const SherlockDataURL = "https://raw.githubusercontent.com/sherlock-project/sherlock/refs/heads/master/sherlock_project/resources/data.json"

type SiteData struct {
	ErrorType string `json:"errorType"`
	ErrorMsg  any    `json:"errorMsg"`

	URL      string `json:"url"`
	URLMain  string `json:"urlMain"`
	URLProbe string `json:"urlProbe"`
	URLError string `json:"errorUrl"`

	UsedUsername   string `json:"username_claimed"`
	UnusedUsername string `json:"username_unclaimed"`
	RegexCheck     string `json:"regexCheck"`
}

// LoadSites loads sherlock-style data.json but safely ignores the top-level "$schema".
func LoadSites(filename string) (map[string]SiteData, error) {
	raw, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var entries map[string]json.RawMessage
	if err := json.Unmarshal(raw, &entries); err != nil {
		return nil, fmt.Errorf("parse json: %w", err)
	}

	out := make(map[string]SiteData, len(entries))
	for siteName, msg := range entries {
		// Skip the JSON Schema entry if present.
		if siteName == "$schema" {
			continue
		}

		var sd SiteData
		if err := json.Unmarshal(msg, &sd); err != nil {
			return nil, fmt.Errorf("site %q: %w", siteName, err)
		}
		out[siteName] = sd
	}

	return out, nil
}

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

func UpdateFromRemote(ctx context.Context, client Doer, userAgent string, destPath string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, SherlockDataURL, nil)
	if err != nil {
		return err
	}
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read a small snippet for diagnostics.
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("download failed: %s (%s)", resp.Status, string(snippet))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return err
	}

	tmp := destPath + ".tmp"
	if err := os.WriteFile(tmp, body, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, destPath)
}
