package counter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	ignore "github.com/sabhiram/go-gitignore"
)

// FileTypeSummary represents the summary of line counts
type FileTypeSummary struct {
	TotalLines int
	Counts     map[string]int // Renamed from TypeCounts to Counts
}

func Count(dirPath string, ignorePath string, separateCount bool) (FileTypeSummary, error) {
	summary := FileTypeSummary{
		TotalLines: 0,
		Counts:     make(map[string]int),
	}
	ignore, err := ignore.CompileIgnoreFile(ignorePath)

	if err != nil {
		fmt.Println("Error reading ignore file:", err)
		return summary, err
	}

	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if ignore.MatchesPath(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() {
			return nil
		}

		lineCount, err := countFileLines(path)
		if err != nil {
			return err
		}
		summary.TotalLines += lineCount

		if separateCount {
			// Get file extension or name if no extension
			fileType := filepath.Ext(path)
			if fileType == "" {
				fileType = filepath.Base(path)
			}
			summary.Counts[fileType] += lineCount
		}

		return nil
	})

	return summary, err
}

func countFileLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("\nError opening file:", err)
		return 0, err
	}
	defer file.Close()

	// Buffer for reading the file in chunks
	buf := make([]byte, 1024)
	lineCount := 0
	prevChar := byte(0)

	// Read the file and count newlines
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading file:", err)
			return 0, err
		}
		// Count newline characters in the current chunk
		for i, b := range buf[:n] {
			if b == '\n' && prevChar != '\r' {
				lineCount++
			} else if b == '\r' {
				lineCount++
				// Check if the next character is '\n' to avoid double counting
				if i+1 < n && buf[i+1] == '\n' {
					prevChar = b
					continue
				}
			}
			prevChar = b
		}
		// Break if end of file
		if err == io.EOF {
			break
		}
	}

	return lineCount, nil
}
