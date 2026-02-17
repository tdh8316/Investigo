package downloaders

import (
	"context"
	"log"
	"net/http"
)

// DownloaderFunc is a site-specific downloader hook.
// Called only when a username is found on that site.
// - profileURL: the found profile URL
// - outDir: directory dedicated to this site (e.g. results/<user>/downloads/instagram)
// - logger: use for user-visible logs (do NOT os.Exit / log.Fatal)
type DownloaderFunc func(ctx context.Context, client *http.Client, profileURL, outDir string, logger *log.Logger) error

// Downloaders is keyed by lowercase site name.
var Downloaders = map[string]DownloaderFunc{
	"instagram": DownloadInstagram,
}
