# oneapi-cli tool

`oneapi-cli` is a tool to help you get started with Intel<sup>®</sup> oneAPI

## Where to find Intel oneAPI.

This tool does not provide any of the tools that may be required to compile/run the samples `oneapi-cli` can extract for you.

Please visit https://software.intel.com/en-us/oneapi for details.

## Development Install 

Fetch using 
```bash
go get github.com/intel/oneapi-cli
``` 
Alternativly see the tags/releases for a binary build for your OS.

## Building
Go 1.13 should be used to build the CLI/TUI app.

```bash
git clone https://github.com/intel/oneapi-cli.git
cd oneapi-cli
go build
./oneapi-cli
```

There is also a `build.sh` which will embed version information within the build.