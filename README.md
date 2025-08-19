# oneapi-cli tool
[![Go Report Card](https://goreportcard.com/badge/github.com/intel/oneapi-cli)](https://goreportcard.com/report/github.com/intel/oneapi-cli)

`oneapi-cli` is a tool to help you get started with Intel<sup>®</sup> oneAPI

## Where to find Intel<sup>®</sup>  oneAPI.

This tool does not provide any of the tools that may be required to compile/run the samples `oneapi-cli` can extract for you.

Please visit https://software.intel.com/en-us/oneapi for details.

## Development Install 

Fetch using 
```bash
go get github.com/intel/oneapi-cli
``` 
Alternatively see the tags/releases for a binary build for your OS.

## Building
Go 1.24.4 should be used to build the CLI/TUI app.

```bash
git clone https://github.com/intel/oneapi-cli.git
cd oneapi-cli
go build
./oneapi-cli
```

There is also a `build.sh` which will embed version information within the build.