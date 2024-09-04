package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/diegoasanch/line_counter/config"
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
		Version:                "0.0.1",
		EnableBashCompletion:   true,
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
	},
	&cli.BoolFlag{
		Name:    "time",
		Aliases: []string{"t"},
		Usage:   "Show total execution time",
	},
	&cli.BoolFlag{
		Name:    "json",
		Aliases: []string{"j"},
		Usage:   "Print output in JSON format",
	},
	&cli.BoolFlag{
		Name:    "pretty",
		Aliases: []string{"p"},
		Usage:   "Enable pretty formatting. For JSON output: indent the JSON. For normal output: use abbreviated numbers (e.g., 23K instead of 23000)",
	},
}

func action(c *cli.Context) error {
	if c.NArg() < 1 {
		cli.ShowAppHelp(c)
		return fmt.Errorf("\nerror: directory path is required")
	}
	startTime := time.Now()

	cfg := config.Config{
		SeparateCount: c.Bool("separate"),
		ShowTime:      c.Bool("time"),
		JSONOutput:    c.Bool("json"),
		PrettyPrint:   c.Bool("pretty"),
	}

	dirPath := c.Args().Get(0)
	ignorePath := filepath.Join("./", "IGNORE.txt")

	summary, err := counter.Count(dirPath, ignorePath, cfg.SeparateCount)
	if err != nil {
		return fmt.Errorf("error counting lines: %w", err)
	}

	executionTime := time.Since(startTime).Seconds()
	return formatter.PrintSummary(summary, cfg, executionTime)
}
