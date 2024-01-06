# access2csv

A tool to parse and convert Apache's (Combined Log Format) access.log to csv

## Installation

- Download a prebuilt binary from [release page](https://github.com/illbison/access2csv/releases/latest)

  _or_
- `git clone https://github.com/illbison/access2csv ; cd access2csv ; go get ; go build -ldflags="-s -w" .`

## Usage

```console
Usage: access2csv -f <PATH> -o <PATH>

    Convert Apache's (Combined Log Format) access.log to csv

Required:
  -f <PATH>    Path to access.log file
  -o <PATH>    Path to output file

Optional:
  -h           Show this message and exit
  -v           Show version and exit
```
