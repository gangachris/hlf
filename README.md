# hlf-cli [![Go Report Card](https://goreportcard.com/badge/github.com/gangachris/hlf-cli)](https://goreportcard.com/report/github.com/gangachris/hlf-cli) [![Build Status](https://travis-ci.org/gangachris/hlf-cli.svg?branch=master)](https://travis-ci.org/gangachris/hlf-cli)

An attempt to build a tiny cli to help setting up a hyperledger fabric environment quickly.

## Progress

- [x] Downloading platform binaries

- [x] Downloading Docker Images

- [x] Download Fabric Samples

- [ ] Spinning Up an example network

- [ ] Spinning up a network based on a configs (custom configtx, cryptoconfig)

- [ ] Hyperledger Composer (maybe)

- [ ] Deploy (maybe)

### Setup

Clone the repository and make sure you have make installed

```
make install
hlf
```

#### Download Prerequisites
**NOTE** samples not yet implemented
```
hlf download // this will all images, binaries and samples
hlf download [images, samples, binaries] // specify what to download e.g hlf download images
```
