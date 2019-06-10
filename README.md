# Dirstalk
[![Build Status](https://travis-ci.com/stefanoj3/dirstalk.svg?branch=master)](https://travis-ci.com/stefanoj3/dirstalk)
[![codecov](https://codecov.io/gh/stefanoj3/dirstalk/branch/master/graph/badge.svg)](https://codecov.io/gh/stefanoj3/dirstalk)
[![Scrutinizer Code Quality](https://scrutinizer-ci.com/g/stefanoj3/dirstalk/badges/quality-score.png?b=master)](https://scrutinizer-ci.com/g/stefanoj3/dirstalk/?branch=master)
![Docker Pulls](https://img.shields.io/docker/pulls/stefanoj3/dirstalk.svg)
![GitHub](https://img.shields.io/github/license/stefanoj3/dirstalk.svg)

Dirstalk is a multi threaded application designed to brute force paths on web servers.

The tool contains functionalities similar to the ones offered by
[dirbuster](https://www.owasp.org/index.php/Category:OWASP_DirBuster_Project) 
and [dirb](https://tools.kali.org/web-applications/dirb).

## Contents
- [How to use it](#-how-to-use-it)
    - [Scan](#scan)
    - [Useful resources](#useful-resources)
    - [Dictionary generator](#dictionary-generator)
- [Download](#-download)
- [Development](#-development)
- [License](https://github.com/stefanoj3/dirstalk/blob/master/LICENSE.md)
- [Milestones](#-milestones)

## [↑](#contents) How to use it

### Scan

Perform a scan with the least amount of parameters (target and dictionary are the only mandatory ones):
```bash
dirstalk scan http://someaddress.url/ --dictionary mydictionary.txt
```

You can get the application to print all the optional parameters:
```bash
dirstalk scan -h
```

##### Example of a customized scan:
```bash
dirstalk scan http://someaddress.url/ \
--dictionary mydictionary.txt \
--http-methods GET,POST \
--http-timeout 10000 \
--scan-depth 10 \
--threads 10 \
--socks5 127.0.0.1:9150 \
--cookie name=value \
--use-cookie-jar \
--user-agent my_user_agent \
--header "Authorization: Bearer 123"

```


##### Explained:
- `--dictionary` to specify the dictionary file - can be a local file or a public remote url
- `--http-methods` to specify which HTTP methods to use for the scan (default `GET`) 
- `--http-timeout` request timeout in millisecond
- `--scan-depth` the maximum recursion depth
- `--threads` the number of threads performing concurrent requests
- `--socks5` SOCKS5 server to connect to (all the requests including the one to fetch the remote dictionary will go through it)
- `--cookie` cookie to add to each request; eg name=value (can be specified multiple times)
- `--use-cookie-jar` enables the use of a cookie jar: it will retain any cookie sent from the server and send them for the following requests
- `--user-agent` user agent to use for http requests
- `--header` header to add to each request; eg name=value (can be specified multiple times)

##### Useful resources
- [here](https://github.com/dustyfresh/dictionaries/tree/master/DirBuster-Lists) you can find dictionaries that can be used with dirstalk
- [tordock](https://github.com/stefanoj3/tordock) is a containerized Tor SOCKS5 that you can use easily with dirstalk 
(just `docker run -d -p 127.0.0.1:9150:9150 stefanoj3/tordock:latest` and then when launching a
 scan specify the following flag: `--socks5 127.0.0.1:9150`)

### Dictionary generator
Dirstalk can also produce it's own dictionaries, useful for example if you
want to check if a specific set of files is available on a given web server.

##### Example:
```bash
dirstalk dictionary.generate /path/to/local/files --out mydictionary.txt
```
The result will be printed to the stdout if no out flag is specified.

## [↑](#contents) Download
You can download a release from [here](https://github.com/stefanoj3/dirstalk/releases)
or you can use a docker image. (eg `docker run stefanoj3/dirstalk dirstalk <cmd>`)


## [↑](#contents) Development
All you need to do local development is to have [make](https://www.gnu.org/software/make/)
and [golang](https://golang.org/) available and the GOPATH correctly configured.

Then you can just:
```bash
go get github.com/stefanoj3/dirstalk         # (or your fork) to obtain the source code
cd $GOPATH/src/github.com/stefanoj3/dirstalk # to go inside the project folder
make dep                                     # to fetch all the required tools and dependencies
make test                                    # to run the test suite
make check                                   # to check for any code style issue
make fix                                     # to automatically fix the code style using goimports
make build                                   # to build an executable for your host OS (not tested under windows) 
```

[dep](https://github.com/golang/dep) is the tool of choice for dependency management.

If you have [gawk](https://www.gnu.org/software/gawk/) available in your PATH you can also just launch
```bash
make help
```
to obtain a description of every command available in the Makefile.

Wanna add a functionality? fix a bug? fork and create a PR.

## [↑](#contents) Milestones

### Version 1.0
- [x] Scan a given url
- [x] Specify how many threads to use
- [x] Specify http verbs to use
- [x] Specify dictionary to use (either local or remote file)
- [x] Specify http request timeout
- [x] Specify verbosity of the log
- [x] Specify scan depth
- [x] Socks5 support
- [x] Can generate dictionary starting from a folder containing files/subdirectories
- [x] Can be compiled for multiple platform/architectures (min: OSX x64, Linux 386, Linux x64, Linux arm, Linux arm64)
- [x] A CI is running the tests suite and display the code coverage
- [x] Print results as a tree and a detailed list of the results found
- [x] Has detailed documentation

### Someday
- [ ] Add support for rotating SOCKS5 proxies
- [ ] Scan a website pages looking for links to bruteforce
- [ ] Expose a webserver that can be used to launch scans and check their status
- [ ] Introduce metrics that can give a sense of how much of the dictionary was found on the remote server
