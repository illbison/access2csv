# access2csv

A tool to parse and convert Apache's (Combined Log Format) access.log to csv

## Installation

- Download a prebuilt binary from [release page](https://github.com/illbison/access2csv/releases/latest)

  _or_
- `git clone https://github.com/illbison/access2csv ; cd access2csv ; go get ; go build -ldflags="-s -w" .`

## Usage

```console
Convert Apache's (Combined Log Format) access.log to csv
access2csv 1.0.0
Usage: access2csv [-f <PATH>] [-o <PATH>]

Options:
  -f <PATH>              path to access.log file
  -o <PATH>              path to output file
  --help, -h             display this help and exit
  --version              display version and exit
```
