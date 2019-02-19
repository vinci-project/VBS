# VBS Client

The Go Daemon of VBS Client. The Daemon is intended to be used with VBS Client GUI at https://github.com/vinci/VBS/client-app. The executables are shipped with VBS Client GUI.

## Install

```
go get -u https://github.com/vinci/VBS/client
```
You will have the executive at `$GOPATH/bin/client.exe`

## Build

After getting the package, in your `$GOPATH/src/github.com/vinci/VBS/node`

Run all dependency fetch with

`go get -d ./...`

Or install separately

```
go get github.com/gin-gonic/gin
```
Then Run the build with

`go build`

## Usage

Simply double click to run or open with command line opens the go dammon server that listens on port `10080`
