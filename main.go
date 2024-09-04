package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/diegoasanch/line_counter/counter"
	"github.com/diegoasanch/line_counter/formatter"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:                   "Line counter",
		Usage:                  "Count lines of code in a directory",
		UsageText:              "line_counter [options] DIRECTORY",
		UseShortOptionHandling: true,
		Flags:                  flags,
		ArgsUsage:              "DIRECTORY",
		Action:                 action,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

var flags = []cli.Flag{
	&cli.BoolFlag{
		Name:    "separate",
		Aliases: []string{"s"},
		Usage:   "Get separate per-file count",
		Value:   false,
	},
	&cli.BoolFlag{
		Name:    "time",
		Aliases: []string{"t"},
		Usage:   "Show total execution time",
		Value:   false,
	},
	&cli.BoolFlag{
		Name:    "json",
		Aliases: []string{"j"},
		Usage:   "Print output in JSON format",
		Value:   false,
	},
	&cli.BoolFlag{
		Name:    "pretty",
		Aliases: []string{"p"},
		Usage:   "Enable pretty formatting. For JSON output: indent the JSON. For normal output: use abbreviated numbers (e.g., 23K instead of 23000)",
		Value:   false,
	},
}

func action(c *cli.Context) error {
	if c.NArg() < 1 {
		cli.ShowAppHelp(c)
		return fmt.Errorf("\nerror: directory path is required")
	}
	startTime := time.Now()

	dirPath := c.Args().Get(0)

	separateCount := c.Bool("separate")
	showTime := c.Bool("time")
	jsonOutput := c.Bool("json")
	prettyPrint := c.Bool("pretty")

	ignorePath := filepath.Join("./", "IGNORE.txt")
	summary, err := counter.Count(dirPath, ignorePath, separateCount)
	if err != nil {
		return fmt.Errorf("error counting lines: %w", err)
	}

	executionTime := time.Since(startTime).Seconds()
	return formatter.PrintSummary(summary, separateCount, showTime, executionTime, prettyPrint, jsonOutput)
}
