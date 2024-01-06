package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	arg "github.com/alexflint/go-arg"
)

type args struct {
	File   string `arg:"-f,--" placeholder:"<PATH>" help:"path to access.log file"`
	Output string `arg:"-o,--" placeholder:"<PATH>" help:"path to output file"`
}

func (args) Description() string {
	return "Convert Apache's (Combined Log Format) access.log to csv"
}

func (args) Version() string {
	return "access2csv v1.0.1"
}

func (args args) CheckArgs() error {
	if args.File == "" {
		return fmt.Errorf("-f is required")
	}
	if args.Output == "" {
		return fmt.Errorf("-o is required")
	}
	return nil
}

func main() {
	args := args{}
	arg.MustParse(&args)

	if err := args.CheckArgs(); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	data, err := os.ReadFile(args.File)
	if err != nil {
		fmt.Printf("Error: failed to open %s\n", args.File)
		os.Exit(1)
	}

	content := strings.Split(string(data), "\n")

	pattern := regexp.MustCompile(`^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\s([^\s]+)\s([^\s]+)\s\[([^\]]+)\]\s"([A-Z]+)\s([^\s]+)\s([^"]+)"\s(\d{3})\s([^\s]+)\s"(.*)"\s"(.*)"$`)

	var parsed [][]string

	for index, value := range content {
		if value == "" {
			continue
		}

		match := pattern.FindStringSubmatch(value)

		if len(match) == 0 {
			fmt.Printf("Error: malformed structure at line %d\n", index+1)
			continue
		}

		// modify the size field to show 0 instead of -
		if match[9] == "-" {
			match[9] = "0"
		}

		// change the timestamp format
		ts, err := time.Parse("02/Jan/2006:15:04:05 -0700", match[4])
		if err != nil {
			fmt.Printf("Error: failed to parse timestamp at line %d\n", index+1)
		} else {
			match[4] = ts.Format("02/01/2006 15:04:05 MST")
		}

		parsed = append(parsed, match[1:])
	}

	output, err := os.Create(args.Output)
	if err != nil {
		fmt.Printf("Error: failed to create %s\n", args.Output)
		os.Exit(1)
	}
	defer output.Close()

	fields := []string{"Host", "Clientid", "Userid", "Timestamp", "Method", "Resource", "Protocol", "Status", "Size", "Referer", "User-agent"}

	writer := csv.NewWriter(output)

	if err = writer.Write(fields); err != nil {
		fmt.Printf("Error: failed to write headers to %s\n", args.Output)
		os.Exit(1)
	}

	if err = writer.WriteAll(parsed); err != nil {
		fmt.Printf("Error: failed to write data to %s\n", args.Output)
		os.Exit(1)
	}
}
