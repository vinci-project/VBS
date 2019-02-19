# VBS Storage Node

## Install

```
go get -u https://github.com/vinci/VBS/node
```
You will have the executive at `$GOPATH/bin/node.exe`

## Build

After getting the package, in your `$GOPATH/src/github.com/vinci/VBS/node`

Run all dependency fetch with

`go get -d ./...`

Or install separately

```
go get github.com/urfave/cli
go get github.com/hashicorp/go-version
go get github.com/shirou/gopsutil
go get github.com/go-redis/redis
go get github.com/surol/speedtest-cli/speedtest
```
Then Run the build with

`go build`

## Usage

Run exe with `walletAddress` and `storagePath`

```
node.exe --walletAddress=xxx-xxx-xxxx-xxxx --storagePath=/etc/
```

### Hardware Check

Add `--hardware` flag to run hardware check

```
node.exe --walletAddress=xxx-xxx-xxxx-xxxx --storagePath=/etc/ --hardware
```

Sample Log

```
Node ID: 9cb795a9-1ce0-426e-a5d9-a03a00d0e8df
-------------------------------------------------------
Os Version: 10.0.17134 Build 17134
-------------------------------------------------------
Go Version: go1.11.2
-------------------------------------------------------
Redis Version: No Redis Found
-------------------------------------------------------
CPU type:Intel(R) Core(TM) i7-7700 CPU @ 3.60GHz, Mhz:3601
-------------------------------------------------------
Memory Size:17116762112, free memory size:10313175040
-------------------------------------------------------
Disk Model For Storage Path:NTFS
-------------------------------------------------------
Disk Available For Storage Path:65229668352
-------------------------------------------------------
Checking Hardware Requirements....
-------------------------------------------------------
Current OS 10.0.17134
Expected OS Windows 7+
OS Version Match true
-------------------------------------------------------
Current Go Version 1.11.2
Expected Go Version 1.11
Go Version Match true
-------------------------------------------------------
No Redis Found
-------------------------------------------------------
Current CPU Type 3.6Ghz
Expected CPU Type 2Ghz+
CPU Type Match true
-------------------------------------------------------
Current Free Memory 10313175040
Expected Free Memory 2GB+
Free Memory 2GB+ Match true
-------------------------------------------------------
Current Harddisk Available 65229668352
Expected Harddisk Available 1TB+
Harddisk Available Match false
-------------------------------------------------------
Storage path 'c:/go/' has permission
==> Done deleting Test file
-------------------------------------------------------
Doing Speedtest....
2018/12/21 02:45:58 Ping: 16 ms
Current Bandwith 45.54 M
Expected Bandwith 100M
Bandwith Available Match true
```

## Test

Run Test

`go test`

## Help

`node.exe --help`
