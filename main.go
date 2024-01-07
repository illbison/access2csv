package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"
)

type args struct {
	File        string
	Output      string
	ShowVersion bool
}

func main() {
	args := ParseArgs()

	if args.ShowVersion {
		ShowVersion()
		os.Exit(0)
	}

	if err := args.CheckRequiredArgs(); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	// open the source log file for reading line by line
	srcFile, err := os.Open(args.File)
	if err != nil {
		fmt.Printf("Error: failed to open %s\n", args.File)
		os.Exit(1)
	}
	defer srcFile.Close()
	scanner := bufio.NewScanner(srcFile)

	// open the output file to write each line parsed
	outFile, err := os.Create(args.Output)
	if err != nil {
		fmt.Printf("Error: failed to create %s\n", args.Output)
		os.Exit(1)
	}
	defer outFile.Close()

	// create a csv writer and write headers
	csvWriter := csv.NewWriter(outFile)
	csvWriter.Write([]string{"Host", "Clientid", "Userid", "Timestamp", "Method", "Resource", "Protocol", "Status", "Size", "Referer", "User-agent"})
	csvWriter.Flush()

	pattern := regexp.MustCompile(`^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\s([^\s]+)\s([^\s]+)\s\[([^\]]+)\]\s"([A-Z]+)\s([^\s]+)\s([^"]+)"\s(\d{3})\s([^\s]+)\s"(.*)"\s"(.*)"$`)

	for i := 1; scanner.Scan(); i++ {
		// read line and skip if empty
		line := scanner.Text()
		if line == "" {
			continue
		}

		match := pattern.FindStringSubmatch(line)

		if len(match) == 0 {
			fmt.Printf("Error: malformed structure at line %d\n", i)
			continue
		}

		// modify the size field to show 0 instead of -
		if match[9] == "-" {
			match[9] = "0"
		}

		// change the timestamp format
		ts, err := time.Parse("02/Jan/2006:15:04:05 -0700", match[4])
		if err != nil {
			fmt.Printf("Error: failed to parse timestamp at line %d\n", i)
		} else {
			match[4] = ts.Format("02/01/2006 15:04:05 MST")
		}

		// write parsed log slice to the output file
		csvWriter.Write(match[1:])
		csvWriter.Flush()
	}
}

func ParseArgs() args {
	args := args{}

	flag.StringVar(&args.File, "f", "", "")
	flag.StringVar(&args.Output, "o", "", "")
	flag.BoolVar(&args.ShowVersion, "v", false, "")

	flag.Usage = func() {
		fmt.Println(`Usage: access2csv -f <PATH> -o <PATH>

    Convert Apache's (Combined Log Format) access.log to csv

Required:
  -f <PATH>    Path to access.log file
  -o <PATH>    Path to output file

Optional:
  -h           Show this message and exit
  -v           Show version and exit`)

	}

	flag.Parse()

	return args
}

func ShowVersion() {
	fmt.Println("access2csv 1.1.0")
}

func (args args) CheckRequiredArgs() error {
	if args.File == "" {
		return fmt.Errorf("-f is required")
	}
	if args.Output == "" {
		return fmt.Errorf("-o is required")
	}
	return nil
}
