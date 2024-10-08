package formatter

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/diegoasanch/line_counter/config"
	"github.com/diegoasanch/line_counter/counter"
)

func PrintSummary(summary counter.FileTypeSummary, cfg config.Config, executionTime float64) error {
	if cfg.JSONOutput {
		return PrintJsonSummary(summary, cfg, executionTime)
	} else {
		return PrintNormalSummary(summary, cfg, executionTime)
	}
}

func PrintJsonSummary(summary counter.FileTypeSummary, cfg config.Config, executionTime float64) error {
	type JsonOutput struct {
		TotalLines int            `json:"total_lines"`
		Counts     map[string]int `json:"counts,omitempty"`
		Runtime    float64        `json:"runtime,omitempty"`
	}

	output := JsonOutput{
		TotalLines: summary.TotalLines,
	}
	if cfg.SeparateCount {
		output.Counts = summary.Counts
	}
	if cfg.ShowTime {
		output.Runtime = executionTime
	}

	var jsonData []byte
	var err error

	if cfg.PrettyPrint {
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

func PrintNormalSummary(summary counter.FileTypeSummary, cfg config.Config, executionTime float64) error {
	if cfg.PrettyPrint {
		fmt.Printf("total %s\n", FormatNumber(summary.TotalLines))
	} else {
		fmt.Printf("total %d\n", summary.TotalLines)
	}

	if cfg.SeparateCount {
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
			if cfg.PrettyPrint {
				fmt.Printf("%-30s %10s\n", kv.Key, FormatNumber(kv.Value))
			} else {
				fmt.Printf("%-30s %10d\n", kv.Key, kv.Value)
			}
		}
	}

	if cfg.ShowTime {
		fmt.Printf("\nruntime %.2fs\n", executionTime)
	}

	return nil
}

func FormatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	return fmt.Sprintf("%dK", n/1000)
}
