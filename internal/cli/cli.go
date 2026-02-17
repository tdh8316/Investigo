package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/fatih/color"
)

var ErrHelp = errors.New("help requested")

type Options struct {
	NoColor         bool
	NoOutput        bool
	Verbose         bool
	UpdateBeforeRun bool
	Test            bool
	WithTor         bool
	Download        bool

	DataFile    string
	Sites       []string
	Timeout     time.Duration
	Concurrency int
	ResultsDir  string
}

const usageText = `
usage:
  investigo [flags] USERNAME [USERNAMES...]
  investigo --test

positional arguments:
  USERNAMES             one or more usernames to investigate

flags:
  -h, --help            show this help message and exit
  --no-color            disable colored stdout output
  --no-output           disable file output
  --update              update database before run from Sherlock repository
  -t, --tor             use tor proxy
  -v, --verbose         verbose output
  -d, --download        download the contents of site if available
  --test                validate sites using username_claimed/unclaimed pairs

options:
  --database PATH       use custom database (default: data.json)
  --sites S1,S2,...     specific sites to investigate separated by comma (default: all sites)
  --timeout SECONDS     HTTP request timeout (default: 60)
  --concurrency N       max concurrent requests (default: 32)
  --results DIR         output directory (default: results)
`

func Parse(args []string, stdout, stderr io.Writer) (Options, []string, error) {
	var opts Options
	var (
		help     bool
		sitesCSV string
		timeoutS int
	)

	fs := flag.NewFlagSet("investigo", flag.ContinueOnError)
	fs.SetOutput(stderr)

	fs.Usage = func() {
		_, _ = fmt.Fprint(stdout, usageText)
	}

	// Help
	fs.BoolVar(&help, "h", false, "show help")
	fs.BoolVar(&help, "help", false, "show help")

	// Behavior flags
	fs.BoolVar(&opts.NoColor, "no-color", false, "disable colored output")
	fs.BoolVar(&opts.NoOutput, "no-output", false, "disable file output")
	fs.BoolVar(&opts.UpdateBeforeRun, "update", false, "update database before run")
	fs.BoolVar(&opts.Test, "test", false, "validate site database")
	fs.BoolVar(&opts.Verbose, "v", false, "verbose output")
	fs.BoolVar(&opts.Verbose, "verbose", false, "verbose output")
	fs.BoolVar(&opts.WithTor, "t", false, "use tor proxy")
	fs.BoolVar(&opts.WithTor, "tor", false, "use tor proxy")
	fs.BoolVar(&opts.Download, "d", false, "download contents if downloader exists")
	fs.BoolVar(&opts.Download, "download", false, "download contents if downloader exists")

	// Options
	fs.StringVar(&opts.DataFile, "database", "data.json", "custom database path")
	fs.StringVar(&sitesCSV, "sites", "", "comma-separated site list")
	fs.StringVar(&sitesCSV, "site", "", "comma-separated site list (compat)") // compat with old flag
	fs.IntVar(&timeoutS, "timeout", 60, "request timeout in seconds")
	fs.IntVar(&opts.Concurrency, "concurrency", 32, "max concurrent requests")
	fs.StringVar(&opts.ResultsDir, "results", "results", "results output directory")

	if err := fs.Parse(args); err != nil {
		return Options{}, nil, err
	}
	if help {
		fs.Usage()
		return Options{}, nil, ErrHelp
	}

	if timeoutS <= 0 {
		// Don't allow zero or negative timeouts; reset to default.
		timeoutS = 60
		if opts.NoColor {
			fmt.Fprintf(stdout, "[!] Invalid timeout value; using default of 60 seconds.\n")
		} else {
			fmt.Fprintf(color.Output, "[%s] Invalid timeout value; using default of %s.\n",
				color.HiRedString("!"),
				color.HiYellowString("60 seconds"),
			)
		}
	}
	opts.Timeout = time.Duration(timeoutS) * time.Second

	if opts.Concurrency <= 0 {
		opts.Concurrency = 32
	}

	if sitesCSV != "" {
		raw := strings.Split(sitesCSV, ",")
		opts.Sites = make([]string, 0, len(raw))
		for _, s := range raw {
			s = strings.TrimSpace(s)
			if s != "" {
				opts.Sites = append(opts.Sites, s)
			}
		}
		// Old behavior: when specifying sites, force verbose so you see misses/errors.
		opts.Verbose = true
	}

	usernames := fs.Args()
	return opts, usernames, nil
}
