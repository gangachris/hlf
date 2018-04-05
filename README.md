# hlf-cli [![Go Report Card](https://goreportcard.com/badge/github.com/gangachris/hlf-cli)](https://goreportcard.com/report/github.com/gangachris/hlf-cli) [![Build Status](https://travis-ci.org/gangachris/hlf-cli.svg?branch=master)](https://travis-ci.org/gangachris/hlf-cli)

An attempt to build a tiny cli to help setting up a hyperledger fabric environment quickly.

## Progress

- [x] Downloading platform binaries

- [x] Downloading Docker Images

- [ ] Spinning Up an example network

- [ ] Spinning up a network based on a configs (custom configtx, cryptoconfig)

- [ ] Hyperledger Composer (maybe)

- [ ] Deploy (maybe)

### Download Prerequisites

This downloads all the prerequisites required to run a hyperledger fabric instance. Includes platform binaries and docker images.

Note that the platform binaries are downloaded to `$HOME/.hlf-cli` (TODO: should be configurable)

```
go run main.go download -h
```
