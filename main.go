package main

import (
	"encoding/json"
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
			jsonOutput := c.Bool("json")
			prettyPrint := c.Bool("pretty")

			ignorePath := filepath.Join("./", "IGNORE.txt")
			summary, err := counter.Count(dirPath, ignorePath, separateCount)
			if err != nil {
				return fmt.Errorf("[ERROR] %w", err)
			}

			executionTime := time.Since(startTime).Seconds()
			if jsonOutput {
				return printJsonSummary(summary, separateCount, showTime, executionTime, prettyPrint)
			} else {
				return printSummary(summary, separateCount, showTime, executionTime, prettyPrint)
			}
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printJsonSummary(summary counter.FileTypeSummary, separateCount bool, showTime bool, executionTime float64, prettyPrint bool) error {
	type JsonOutput struct {
		TotalLines int            `json:"total_lines"`
		Counts     map[string]int `json:"counts,omitempty"`
		Runtime    float64        `json:"runtime,omitempty"`
	}

	output := JsonOutput{
		TotalLines: summary.TotalLines,
	}
	if separateCount {
		output.Counts = summary.Counts
	}
	if showTime {
		output.Runtime = executionTime
	}

	var jsonData []byte
	var err error

	if prettyPrint {
		jsonData, err = json.MarshalIndent(output, "", "  ")
	} else {
		jsonData, err = json.Marshal(output)
	}

	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}
	fmt.Println(string(jsonData))
	return nil
}

func printSummary(summary counter.FileTypeSummary, separateCount bool, showTime bool, executionTime float64, prettyPrint bool) error {
	if prettyPrint {
		fmt.Printf("total %s\n", formatNumber(summary.TotalLines))
	} else {
		fmt.Printf("total %d\n", summary.TotalLines)
	}

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

		sort.Slice(sortedCounts, func(i, j int) bool {
			return sortedCounts[i].Value > sortedCounts[j].Value
		})

		for _, kv := range sortedCounts {
			if prettyPrint {
				fmt.Printf("%-30s %10s\n", kv.Key, formatNumber(kv.Value))
			} else {
				fmt.Printf("%-30s %10d\n", kv.Key, kv.Value)
			}
		}
	}

	if showTime {
		fmt.Printf("\nruntime %.2fs\n", executionTime)
	}

	return nil
}

func formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	return fmt.Sprintf("%dK", n/1000)
}
