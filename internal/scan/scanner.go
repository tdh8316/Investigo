package scan

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/dlclark/regexp2"

	"github.com/tdh8316/Investigo/internal/data"
	"github.com/tdh8316/Investigo/internal/downloaders"
	"github.com/tdh8316/Investigo/internal/httpx"
)

type Scanner struct {
	client      *http.Client
	cfg         Config
	downloaders map[string]downloaders.DownloaderFunc

	// Cache compiled regexCheck per site
	regexCache    sync.Map // siteName -> *regexp2.Regexp
	regexErrCache sync.Map // siteName -> error
}

func NewScanner(client *http.Client, cfg Config, dls map[string]downloaders.DownloaderFunc) *Scanner {
	if cfg.Concurrency <= 0 {
		cfg.Concurrency = 32
	}
	if cfg.MaxBodyBytes <= 0 {
		cfg.MaxBodyBytes = 2 << 20
	}
	if cfg.UserAgent == "" {
		cfg.UserAgent = httpx.DefaultUserAgent
	}
	if dls == nil {
		dls = map[string]downloaders.DownloaderFunc{}
	}

	return &Scanner{
		client:      client,
		cfg:         cfg,
		downloaders: dls,
	}
}

func (s *Scanner) ScanUsername(
	ctx context.Context,
	username string,
	sites map[string]data.SiteData,
	downloadDir string,
	logger *log.Logger,
	onResult func(Result),
) error {
	if onResult == nil {
		return fmt.Errorf("onResult callback is nil")
	}

	siteNames := make([]string, 0, len(sites))
	for name := range sites {
		siteNames = append(siteNames, name)
	}
	sort.Strings(siteNames)

	workers := min(s.cfg.Concurrency, len(siteNames))
	if workers == 0 {
		return nil
	}

	jobs := make(chan string) // Channel of site names to investigate.
	results := make(chan Result, workers)

	var wg sync.WaitGroup
	wg.Add(workers)
	for range workers {
		// Worker goroutines.
		go func() {
			defer wg.Done()
			for site := range jobs {
				results <- s.Investigo(ctx, username, site, sites[site], downloadDir, logger)
			}
		}()
	}

	// Wait for workers to finish and close results channel when done.
	go func() {
		defer close(results)
		wg.Wait()
	}()

	go func() {
		defer close(jobs)
		for _, site := range siteNames {
			select {
			case <-ctx.Done():
				return
			case jobs <- site:
			}
		}
	}()

	for res := range results {
		onResult(res) // Callback for each result.
	}

	return ctx.Err()
}

func (s *Scanner) ValidateSites(
	ctx context.Context,
	sites map[string]data.SiteData,
	onFailure func(ValidationFailure),
) (int, error) {
	if onFailure == nil {
		return 0, fmt.Errorf("onFailure callback is nil")
	}

	siteNames := make([]string, 0, len(sites))
	for name := range sites {
		siteNames = append(siteNames, name)
	}
	sort.Strings(siteNames)

	workers := min(s.cfg.Concurrency, len(siteNames))
	if workers == 0 {
		return 0, nil
	}

	jobs := make(chan string)
	failures := make(chan ValidationFailure, workers)

	var wg sync.WaitGroup
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			for site := range jobs {
				sd := sites[site]
				if sd.UsedUsername == "" || sd.UnusedUsername == "" {
					// If we can't validate, count as failure (same spirit as old test mode).
					f := ValidationFailure{
						Site:           site,
						UsedUsername:   sd.UsedUsername,
						UnusedUsername: sd.UnusedUsername,
						Used: Result{
							Username: sd.UsedUsername, Site: site, Proxied: s.cfg.WithTor,
							Err: fmt.Errorf("missing username_claimed/username_unclaimed in database"),
						},
						Unused: Result{
							Username: sd.UnusedUsername, Site: site, Proxied: s.cfg.WithTor,
							Err: fmt.Errorf("missing username_claimed/username_unclaimed in database"),
						},
					}
					failures <- f
					continue
				}

				used := s.Investigo(ctx, sd.UsedUsername, site, sd, "", nil)
				unused := s.Investigo(ctx, sd.UnusedUsername, site, sd, "", nil)

				if used.Exists && !unused.Exists {
					continue
				}

				failures <- ValidationFailure{
					Site:           site,
					UsedUsername:   sd.UsedUsername,
					UnusedUsername: sd.UnusedUsername,
					Used:           used,
					Unused:         unused,
				}
			}
		}()
	}

	go func() {
		defer close(failures)
		wg.Wait()
	}()

	go func() {
		defer close(jobs)
		for _, site := range siteNames {
			select {
			case <-ctx.Done():
				return
			case jobs <- site:
			}
		}
	}()

	count := 0
	for f := range failures {
		count++
		onFailure(f)
	}

	return count, ctx.Err()
}

func (s *Scanner) Investigo(
	ctx context.Context,
	username string,
	site string,
	sd data.SiteData,
	downloadDir string,
	logger *log.Logger,
) Result {
	res := Result{
		Username:      username,
		Site:          site,
		URLTemplate:   sd.URL,
		ProbeTemplate: sd.URLProbe,
		Proxied:       s.cfg.WithTor,
	}

	if sd.URL == "" {
		res.Err = fmt.Errorf("missing url in database")
		return res
	}

	profileURL := strings.ReplaceAll(sd.URL, "{}", username)
	probeURL := profileURL
	if sd.URLProbe != "" {
		probeURL = strings.ReplaceAll(sd.URLProbe, "{}", username)
	}

	// Optional username regexCheck (cached per site).
	if sd.RegexCheck != "" {
		re, err := s.getRegex(site, sd.RegexCheck)
		if err != nil {
			res.Err = fmt.Errorf("invalid regexCheck: %w", err)
			return res
		}
		ok, err := re.MatchString(username)
		if err != nil {
			res.Err = fmt.Errorf("regexCheck match error: %w", err)
			return res
		}
		if !ok {
			// Username not valid for this site => treat as not found (no error).
			return res
		}
	}

	req, err := httpx.NewRequest(ctx, http.MethodGet, probeURL, nil, s.cfg.UserAgent)
	if err != nil {
		res.Err = err
		return res
	}

	resp, err := s.client.Do(req)
	if err != nil {
		res.Err = err
		return res
	}
	defer resp.Body.Close()

	switch sd.ErrorType {
	case "status_code":
		if resp.StatusCode == http.StatusOK {
			res.Exists = true
			res.Link = profileURL
		}

	case "message":
		body, err := s.readBody(resp.Body)
		if err != nil {
			res.Err = err
			return res
		}

		notFound, err := containsErrorMessage(body, sd.ErrorMsg)
		if err != nil {
			res.Err = err
			return res
		}

		res.Exists = !notFound
		if res.Exists {
			res.Link = profileURL
		}

	case "response_url":
		finalURL := ""
		if resp.Request != nil && resp.Request.URL != nil {
			finalURL = resp.Request.URL.String()
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 400 && finalURL == profileURL {
			res.Exists = true
			res.Link = profileURL
		}

	default:
		res.Err = fmt.Errorf("unsupported error type %q", sd.ErrorType)
		return res
	}

	// Optional download hook.
	if res.Exists && s.cfg.Download {
		if dl, ok := s.downloaders[strings.ToLower(site)]; ok {
			if downloadDir == "" {
				return res
			}
			siteDir := filepath.Join(downloadDir, strings.ToLower(site))
			_ = os.MkdirAll(siteDir, 0o755)

			// Do not convert download errors into not found (keep [+] result).
			_ = dl(ctx, s.client, profileURL, siteDir, logger)
		}
	}

	return res
}

func (s *Scanner) getRegex(site, expr string) (*regexp2.Regexp, error) {
	if v, ok := s.regexCache.Load(site); ok {
		return v.(*regexp2.Regexp), nil
	}
	if v, ok := s.regexErrCache.Load(site); ok {
		return nil, v.(error)
	}

	re, err := regexp2.Compile(expr, 0)
	if err != nil {
		s.regexErrCache.Store(site, err)
		return nil, err
	}
	s.regexCache.Store(site, re)
	return re, nil
}

func (s *Scanner) readBody(r io.Reader) (string, error) {
	limited := io.LimitReader(r, s.cfg.MaxBodyBytes)
	b, err := io.ReadAll(limited)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func containsErrorMessage(body string, errorMsg any) (bool, error) {
	switch v := errorMsg.(type) {
	case nil:
		return false, fmt.Errorf("errorMsg is missing (nil) for errorType=message")
	case string:
		if v == "" {
			return false, nil
		}
		return strings.Contains(body, v), nil
	case []any:
		for _, it := range v {
			s, ok := it.(string)
			if !ok {
				continue
			}
			if s != "" && strings.Contains(body, s) {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, fmt.Errorf("unsupported errorMsg type %T for errorType=message", errorMsg)
	}
}
