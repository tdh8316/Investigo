package app

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"

	"github.com/tdh8316/Investigo/internal/cli"
	"github.com/tdh8316/Investigo/internal/data"
	"github.com/tdh8316/Investigo/internal/downloaders"
	"github.com/tdh8316/Investigo/internal/httpx"
	"github.com/tdh8316/Investigo/internal/output"
	"github.com/tdh8316/Investigo/internal/scan"
)

func Run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	fmt.Fprintln(stdout, "Investigo - Investigate Users Across Social Networks.")

	opts, usernames, err := cli.Parse(args, stdout, stderr)
	if err != nil {
		if errors.Is(err, cli.ErrHelp) {
			return 0
		}
		fmt.Fprintln(stderr, err.Error())
		return 2
	}

	color.NoColor = opts.NoColor

	httpClient, err := httpx.NewClient(httpx.ClientConfig{
		Timeout:     opts.Timeout,
		WithTor:     opts.WithTor,
		TorProxyURL: httpx.DefaultTorProxyURL,
	})
	if err != nil {
		fmt.Fprintf(stderr, "failed to initialize HTTP client: %v\n", err)
		return 1
	}

	// Load + optionally update database.
	sites, err := loadDatabase(ctx, httpClient, opts, stdout)
	if err != nil {
		fmt.Fprintf(stderr, "database error: %v\n", err)
		return 1
	}

	// Optional: filter sites.
	if len(opts.Sites) > 0 {
		sites = filterSites(sites, opts.Sites, stdout, opts.NoColor)
	}

	// -d/--download with no usernames: show available downloaders.
	if opts.Download && len(usernames) == 0 && !opts.Test {
		printDownloaders(stdout, opts.NoColor)
		return 0
	}

	// Back-compat behavior: if no usernames provided, prompt.
	if len(usernames) == 0 && !opts.Test {
		usernames = promptUsernames(stdout, os.Stdin)
		if len(usernames) == 0 {
			fmt.Fprintln(stderr, "no usernames provided")
			return 2
		}
	}

	// Build scanner once (reuses regex cache + client).
	scanner := scan.NewScanner(httpClient, scan.Config{
		UserAgent:    httpx.DefaultUserAgent,
		WithTor:      opts.WithTor,
		Download:     opts.Download,
		Concurrency:  opts.Concurrency,
		MaxBodyBytes: 2 << 20, // 2 MiB max body read for message checks
	}, downloaders.Downloaders)

	if opts.Test {
		return runTest(ctx, stdout, opts.NoColor, scanner, sites)
	}

	for _, username := range usernames {
		username = strings.TrimSpace(username)
		if username == "" {
			continue
		}

		// Header (stdout).
		if opts.NoColor {
			fmt.Fprintf(stdout, "\nInvestigating %s on:\n", username)
		} else {
			fmt.Fprintf(color.Output, "\nInvestigating %s on:\n", color.HiGreenString(username))
		}

		userDir := filepath.Join(opts.ResultsDir, username)
		if err := os.MkdirAll(userDir, 0o755); err != nil {
			fmt.Fprintf(stderr, "failed to create results dir %q: %v\n", userDir, err)
			return 1
		}

		// Buffer for out.txt (per-username; fixes original “global builder” bug).
		var buf strings.Builder

		printer := output.NewPrinter(stdout, opts.NoColor, opts.Verbose, &buf)

		downloadDir := filepath.Join(userDir, "downloads")
		if opts.Download {
			_ = os.MkdirAll(downloadDir, 0o755)
		}

		// Stream results as they complete.
		if err := scanner.ScanUsername(ctx, username, sites, downloadDir, printer.Logger(), printer.Result); err != nil &&
			!errors.Is(err, context.Canceled) {
			fmt.Fprintf(stderr, "scan error for %q: %v\n", username, err)
			// Keep going to next username or exit; old behavior was “best effort”.
		}

		if !opts.NoOutput {
			outPath := filepath.Join(userDir, "out.txt")
			if err := os.WriteFile(outPath, []byte(buf.String()), 0o600); err != nil {
				fmt.Fprintf(stderr, "failed to write %q: %v\n", outPath, err)
				return 1
			}
		}
	}

	return 0
}

func loadDatabase(ctx context.Context, client httpx.Doer, opts cli.Options, stdout io.Writer) (map[string]data.SiteData, error) {
	_, statErr := os.Stat(opts.DataFile)
	fileExists := statErr == nil

	if opts.UpdateBeforeRun || !fileExists {
		if opts.NoColor {
			fmt.Fprintf(stdout, "[!] Update database: Downloading...")
		} else {
			fmt.Fprintf(color.Output, "[%s] Update database: %s",
				color.HiBlueString("!"),
				color.HiYellowString("Downloading..."),
			)
		}

		if err := data.UpdateFromRemote(ctx, client, httpx.DefaultUserAgent, opts.DataFile); err != nil {
			if fileExists {
				// Fall back to existing database.
				if opts.NoColor {
					fmt.Fprintf(stdout, "[!] Failed to update database: %v (using existing)\n", err)
				} else {
					fmt.Fprintf(color.Output, "[%s] Failed to update database: %s (using existing)\n",
						color.HiRedString("!"),
						color.HiRedString(err.Error()),
					)
				}
			} else {
				return nil, fmt.Errorf("failed to update database and no existing database found: %w", err)
			}
		} else {
			if opts.NoColor {
				fmt.Fprintln(stdout, "[Done]")
			} else {
				fmt.Fprintf(color.Output, "[%s]\n", color.GreenString("Done"))
			}
		}
	}

	return data.LoadSites(opts.DataFile)
}

func filterSites(all map[string]data.SiteData, selected []string, stdout io.Writer, noColor bool) map[string]data.SiteData {
	if len(selected) == 0 {
		return all
	}

	// Build case-insensitive lookup.
	lut := make(map[string]string, len(all))
	for name := range all {
		lut[strings.ToLower(name)] = name
	}

	out := make(map[string]data.SiteData, len(selected))
	var unknown []string

	for _, s := range selected {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		key := strings.ToLower(s)
		if actual, ok := lut[key]; ok {
			out[actual] = all[actual]
		} else {
			unknown = append(unknown, s)
		}
	}

	if len(unknown) > 0 {
		msg := "Unknown sites ignored: " + strings.Join(unknown, ", ")
		if noColor {
			fmt.Fprintf(stdout, "[!] %s\n", msg)
		} else {
			fmt.Fprintf(color.Output, "[%s] %s\n", color.HiRedString("!"), color.HiYellowString(msg))
		}
	}

	if len(out) == 0 {
		msg := "No matching sites found; using full database."
		if noColor {
			fmt.Fprintf(stdout, "[!] %s\n", msg)
		} else {
			fmt.Fprintf(color.Output, "[%s] %s\n", color.HiRedString("!"), color.HiYellowString(msg))
		}
		return all
	}

	if noColor {
		fmt.Fprintf(stdout, "[i] Using %d site(s)\n", len(out))
	} else {
		fmt.Fprintf(color.Output, "[%s] Using %d site(s)\n", color.HiBlueString("i"), len(out))
	}
	return out
}

func promptUsernames(stdout io.Writer, stdin io.Reader) []string {
	fmt.Fprint(stdout, "Enter usernames to investigate separated by a space: ")
	r := bufio.NewReader(stdin)
	line, _ := r.ReadString('\n')
	line = strings.TrimSpace(line)
	return strings.Fields(line)
}

func printDownloaders(stdout io.Writer, noColor bool) {
	fmt.Fprintln(stdout, "List of sites that can download userdata:")

	keys := make([]string, 0, len(downloaders.Downloaders))
	for k := range downloaders.Downloaders {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if noColor {
			fmt.Fprintf(stdout, "[+] %s\n", k)
		} else {
			fmt.Fprintf(color.Output, "[%s] %s\n", color.HiGreenString("+"), color.HiWhiteString(k))
		}
	}
}

func runTest(ctx context.Context, stdout io.Writer, noColor bool, scanner *scan.Scanner, sites map[string]data.SiteData) int {
	if noColor {
		fmt.Fprintln(stdout, "[i] Checking site validity...")
	} else {
		fmt.Fprintf(color.Output, "[%s] Checking site validity...\n", color.HiBlueString("i"))
	}

	failCount, _ := scanner.ValidateSites(ctx, sites, func(f scan.ValidationFailure) {
		// Match the old output style as closely as possible.
		if f.Used.Err != nil || f.Unused.Err != nil {
			var msgParts []string
			if f.Used.Err != nil {
				msgParts = append(msgParts, "["+f.Used.Err.Error()+"]")
			}
			if f.Unused.Err != nil {
				msgParts = append(msgParts, "["+f.Unused.Err.Error()+"]")
			}
			errMsg := strings.Join(msgParts, "")

			if noColor {
				fmt.Fprintf(stdout, "[-] %s: Failed with error %s\n", f.Site, errMsg)
			} else {
				fmt.Fprintf(color.Output, "[-] %s: %s %s\n",
					f.Site,
					color.YellowString("Failed with error"),
					errMsg,
				)
			}
			return
		}

		if noColor {
			fmt.Fprintf(stdout,
				"[-] %s: Not working (%s: expected true, result is %t | %s: expected false, result is %t)\n",
				f.Site,
				f.UsedUsername, f.Used.Exists,
				f.UnusedUsername, f.Unused.Exists,
			)
		} else {
			fmt.Fprintf(color.Output,
				"[-] %s: %s (%s: expected true, result is %t | %s: expected false, result is %t)\n",
				f.Site,
				color.RedString("Not working"),
				f.UsedUsername, f.Used.Exists,
				f.UnusedUsername, f.Unused.Exists,
			)
		}
	})

	if noColor {
		fmt.Fprintln(stdout, "[Done]")
	} else {
		fmt.Fprintf(color.Output, "[%s]\n", color.GreenString("Done"))
	}

	fmt.Fprintf(stdout,
		"\nThese %d sites are not compatible with the Sherlock database.\n"+
			"Please check https://github.com/tdh8316/Investigo/#to-fix-incompatible-sites\n",
		failCount,
	)
	return 0
}
