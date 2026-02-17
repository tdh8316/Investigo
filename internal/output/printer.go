package output

import (
	"io"
	"log"
	"strings"

	"github.com/fatih/color"

	"github.com/tdh8316/Investigo/internal/scan"
)

type Printer struct {
	noColor bool
	verbose bool

	logger *log.Logger
	stream *log.Logger // optional (writes to buffer)
}

func NewPrinter(stdout io.Writer, noColor, verbose bool, buf *strings.Builder) *Printer {
	p := &Printer{
		noColor: noColor,
		verbose: verbose,
		logger:  log.New(stdout, "", 0),
	}
	if buf != nil {
		p.stream = log.New(buf, "", 0)
	}
	return p
}

func (p *Printer) Logger() *log.Logger {
	return p.logger
}

func (p *Printer) Result(result scan.Result) {
	// File output is always plain.
	if p.stream != nil {
		if result.Exists {
			p.stream.Printf("[%s] %s: %s", "+", result.Site, result.Link)
		} else if p.verbose {
			if result.Err != nil {
				p.stream.Printf("[%s] %s: ERROR: %s", "!", result.Site, result.Err.Error())
			} else {
				p.stream.Printf("[%s] %s: %s", "-", result.Site, "Not Found!")
			}
		}
	}

	// Stdout output (colored or not).
	if result.Exists {
		if p.noColor {
			p.logger.Printf("[%s] %s: %s", "+", result.Site, result.Link)
		} else {
			p.logger.Printf("[%s] %s: %s", color.HiGreenString("+"), color.HiWhiteString(result.Site), result.Link)
		}
		return
	}

	if p.verbose {
		if result.Err != nil {
			if p.noColor {
				p.logger.Printf("[%s] %s: ERROR: %s", "!", result.Site, result.Err.Error())
			} else {
				p.logger.Printf("[%s] %s: %s: %s",
					color.HiRedString("!"),
					result.Site,
					color.HiMagentaString("ERROR"),
					color.HiRedString(result.Err.Error()),
				)
			}
			return
		}

		if p.noColor {
			p.logger.Printf("[%s] %s: %s", "-", result.Site, "Not Found!")
		} else {
			p.logger.Printf("[%s] %s: %s", color.HiRedString("-"), result.Site, color.HiYellowString("Not Found!"))
		}
	}
}
