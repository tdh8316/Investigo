package downloaders

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/tidwall/gjson"

	"github.com/tdh8316/Investigo/internal/httpx"
)

func DownloadInstagram(ctx context.Context, client *http.Client, profileURL, outDir string, logger *log.Logger) error {
	if client == nil {
		return errors.New("instagram downloader: nil http client")
	}
	if logger == nil {
		logger = log.New(io.Discard, "", 0)
	}
	profileURL = strings.TrimSpace(profileURL)
	if profileURL == "" {
		return errors.New("instagram downloader: empty profileURL")
	}
	if outDir == "" {
		return errors.New("instagram downloader: empty outDir")
	}

	profileURL = strings.TrimRight(profileURL, "/") + "/"

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("instagram downloader: create outDir: %w", err)
	}

	// Fetch metadata JSON (best effort).
	metaBody, metaURL, err := fetchInstagramMetaJSON(ctx, client, profileURL)
	if err != nil {
		return err
	}

	// Optional: save metadata for debugging.
	_ = os.WriteFile(filepath.Join(outDir, "metadata.json"), metaBody, 0o600)

	uris := collectInstagramMediaURIs(metaBody)

	if len(uris) == 0 {
		return fmt.Errorf("instagram downloader: no downloadable media URLs found (metadata source: %s)", metaURL)
	}

	// Download media with bounded concurrency.
	const maxConcurrentDownloads = 6
	sem := make(chan struct{}, maxConcurrentDownloads)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	for i, mediaURL := range uris {
		wg.Add(1)

		go func() {
			defer wg.Done()

			select {
			case <-ctx.Done():
				mu.Lock()
				errs = append(errs, ctx.Err())
				mu.Unlock()
				return
			case sem <- struct{}{}:
			}
			defer func() { <-sem }()

			if err := downloadOne(ctx, client, profileURL, outDir, i, mediaURL); err != nil {
				logger.Printf("[instagram] download failed: %v", err)
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if len(errs) > 0 {
		// Return a compact summary error
		return fmt.Errorf("instagram downloader: %d download(s) failed (see log/output)", len(errs))
	}

	logger.Printf("[instagram] downloaded %d file(s) into %s", len(uris), outDir)
	return nil
}

func fetchInstagramMetaJSON(ctx context.Context, client *http.Client, profileURL string) ([]byte, string, error) {
	u, err := url.Parse(profileURL)
	if err != nil {
		return nil, "", fmt.Errorf("instagram downloader: parse profileURL: %w", err)
	}

	// Try a couple known variants (these have been unstable historically).
	candidates := []string{
		withQuery(u, "__a=1&__d=dis"),
		withQuery(u, "__a=1"),
	}

	var lastErr error
	for _, metaURL := range candidates {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, metaURL, nil)
		if err != nil {
			lastErr = err
			continue
		}
		req.Header.Set("User-Agent", httpx.DefaultUserAgent)
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("Referer", profileURL)

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 20<<20)) // 20 MiB safety cap
		_ = resp.Body.Close()

		if readErr != nil {
			lastErr = readErr
			continue
		}
		if resp.StatusCode != http.StatusOK {
			// Common: 401/403/429 or HTML challenge.
			snippet := string(body)
			if len(snippet) > 2000 {
				snippet = snippet[:2000]
			}
			lastErr = fmt.Errorf("instagram metadata: %s (url=%s) body=%q", resp.Status, metaURL, snippet)
			continue
		}

		// Basic sanity check: should look like JSON with graphql data.
		if !gjson.GetBytes(body, "graphql.user").Exists() {
			// Some responses may be different JSON or an HTML page even with 200.
			snippet := string(body)
			if len(snippet) > 1000 {
				snippet = snippet[:1000]
			}
			lastErr = fmt.Errorf("instagram metadata: unexpected payload (url=%s) snippet=%q", metaURL, snippet)
			continue
		}

		return body, metaURL, nil
	}

	if lastErr == nil {
		lastErr = errors.New("instagram downloader: failed to fetch metadata (unknown error)")
	}
	return nil, "", lastErr
}

func withQuery(base *url.URL, rawQuery string) string {
	u := *base
	u.RawQuery = rawQuery
	return u.String()
}

func collectInstagramMediaURIs(meta []byte) []string {
	var uris []string
	seen := make(map[string]struct{}, 128)

	add := func(s string) {
		s = strings.TrimSpace(s)
		if s == "" {
			return
		}
		if _, ok := seen[s]; ok {
			return
		}
		seen[s] = struct{}{}
		uris = append(uris, s)
	}

	// Profile picture (try HD then fallback).
	add(gjson.GetBytes(meta, "graphql.user.profile_pic_url_hd").String())
	add(gjson.GetBytes(meta, "graphql.user.profile_pic_url").String())

	addFromNode := func(node gjson.Result) {
		if node.Get("is_video").Bool() {
			if v := node.Get("video_url").String(); v != "" {
				add(v)
				return
			}
		}
		if v := node.Get("display_url").String(); v != "" {
			add(v)
		}
	}

	// Timeline media edges.
	for _, edge := range gjson.GetBytes(meta, "graphql.user.edge_owner_to_timeline_media.edges").Array() {
		node := edge.Get("node")
		addFromNode(node)

		// Sidecar children (if any)
		for _, subEdge := range node.Get("edge_sidecar_to_children.edges").Array() {
			subNode := subEdge.Get("node")
			addFromNode(subNode)
		}
	}

	return uris
}

func downloadOne(ctx context.Context, client *http.Client, referer, outDir string, index int, mediaURL string) error {
	u, err := url.Parse(mediaURL)
	if err != nil {
		return fmt.Errorf("parse media url: %w", err)
	}

	ext := strings.ToLower(path.Ext(u.Path))
	if ext == "" || len(ext) > 10 {
		ext = ".bin"
	}

	filename := fmt.Sprintf("%03d%s", index, ext)
	finalPath := filepath.Join(outDir, filename)
	tmpPath := finalPath + ".part"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, mediaURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", httpx.DefaultUserAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Referer", referer)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("GET %s: %w", mediaURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("GET %s: %s body=%q", mediaURL, resp.Status, string(b))
	}

	f, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("create file %s: %w", tmpPath, err)
	}

	_, copyErr := io.Copy(f, resp.Body)
	closeErr := f.Close()
	if copyErr != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("write %s: %w", tmpPath, copyErr)
	}
	if closeErr != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("close %s: %w", tmpPath, closeErr)
	}

	if err := os.Rename(tmpPath, finalPath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("rename %s -> %s: %w", tmpPath, finalPath, err)
	}

	return nil
}
