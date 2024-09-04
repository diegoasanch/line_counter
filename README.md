# Line Counter

Line Counter is a command-line tool that counts the lines of code in a directory, providing summaries by file type or individual files.

## Features

- Count total lines of code in a directory
- Optionally count lines for each file type or individual files
- Ignore files and directories based on a gitignore-style file
- JSON output option
- Pretty-print option for both normal and JSON output
- Display execution time

## Installation

To install Line Counter, make sure you have Go installed on your system and the [`GOPATH` is set](https://go.dev/wiki/SettingGOPATH), then run:

```bash
go install github.com/diegoasanch/line_counter@latest
```

## Usage

```bash
line_counter [options] DIRECTORY
```

### Options

- `-s, --separate`: Get separate per-file count
- `-t, --time`: Show total execution time
- `-j, --json`: Print output in JSON format
- `-p, --pretty`: Enable pretty formatting (abbreviated numbers for normal output, indented JSON for JSON output)
- `-i, --ignore`: Path to the ignore file (default: ./IGNORE.txt)

#### Combining Options

You can combine multiple options into a single flag. For example: `-sjt` is equivalent to using `-s`, `-j`, and `-t` separately

This allows for more concise command-line usage when multiple options are needed.

### Examples

1. Basic usage:

```plain
$ line_counter /path/to/yourproject
total 1234
```

2. Count with file type breakdown and execution time:

```plain
$ line_counter -s -t /path/to/your/project
total 1234
-----------------------------------------
.go                             1000
.js                              150
.css                              84

runtime 0.05s
```

3. Pretty-printed output:

```plain
$ line_counter -s -p /path/to/your/project
total 1.2K
-----------------------------------------
.go                                    1k
.js                                   150
.css                                   84
```

4. JSON output:

```plain
$ line_counter -j /path/to/your/project
{"total_lines":1234}
```

5. Pretty-printed JSON output with separate counts and execution time:

```plain
$ line_counter -jpst /path/to/your/project
{
    "total_lines": 1234,
    "counts": {
    ".go": 1000,
    ".js": 150,
    ".css": 84
    },
    "runtime": 0.05
}
```

## Ignoring Files

Create an `IGNORE.txt` file in the same directory as the `line_counter` executable. This file follows the `.gitignore` format. For example:

```plain
.log
node_modules/
build/
```

This will ignore all `.log` files, the `node_modules` directory, and the `build` directory.

