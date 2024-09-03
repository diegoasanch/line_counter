package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/diegoasanch/line_counter/counter"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "line_counter",
		Usage: "Count lines of code in a directory",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dir",
				Aliases:  []string{"d"},
				Usage:    "Directory to count lines in",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			dirPath := c.String("dir")
			if dirPath == "" {
				return fmt.Errorf("directory path is required")
			}

			ignorePath := filepath.Join("./", "IGNORE.txt")
			summary, err := counter.Count(dirPath, ignorePath)
			if err != nil {
				return fmt.Errorf("[ERROR] %w", err)
			}

			printSummary(summary)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func printSummary(summary counter.FileTypeSummary) {
	fmt.Printf("%-30s %10s\n", "Total lines:", formatNumber(summary.TotalLines))

	fmt.Println(strings.Repeat("-", 41))
	// Create a slice of key-value pairs for sorting
	type kv struct {
		Key   string
		Value int
	}
	var sortedCounts []kv
	for k, v := range summary.TypeCounts {
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

func formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	return fmt.Sprintf("%dK", n/1000)
}
