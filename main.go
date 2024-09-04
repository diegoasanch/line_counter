package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/diegoasanch/line_counter/counter"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:                   "Line counter",
		Usage:                  "Count lines of code in a directory",
		UsageText:              "line_counter [options] DIRECTORY",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
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
		},
		ArgsUsage: "DIRECTORY",
		Action: func(c *cli.Context) error {
			startTime := time.Now()

			if c.NArg() < 1 {
				return fmt.Errorf("directory path is required")
			}
			dirPath := c.Args().Get(0)

			separateCount := c.Bool("separate")
			showTime := c.Bool("time")

			ignorePath := filepath.Join("./", "IGNORE.txt")
			summary, err := counter.Count(dirPath, ignorePath, separateCount)
			if err != nil {
				return fmt.Errorf("[ERROR] %w", err)
			}

			printSummary(summary, separateCount)

			if showTime {
				executionTime := time.Since(startTime).Seconds()
				fmt.Printf("\nExecution time: %.2f seconds\n", executionTime)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func printSummary(summary counter.FileTypeSummary, separateCount bool) {
	fmt.Printf("%-30s %10s\n", "Total lines:", formatNumber(summary.TotalLines))

	// Create a slice of key-value pairs for sorting
	if separateCount {
		fmt.Println(strings.Repeat("-", 41))
		type kv struct {
			Key   string
			Value int
		}
		var sortedCounts []kv
		for k, v := range summary.Counts {
			sortedCounts = append(sortedCounts, kv{k, v})
		}

		// Sort by value (line count) in descending order
		sort.Slice(sortedCounts, func(i, j int) bool {
			return sortedCounts[i].Value > sortedCounts[j].Value
		})

		// Print sorted results
		for _, kv := range sortedCounts {
			fmt.Printf("%-30s %10s\n", kv.Key, formatNumber(kv.Value))
		}
	}
}

func formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	return fmt.Sprintf("%dK", n/1000)
}
