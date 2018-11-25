# Dirstalk
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/03cd98af06aa472cbd5587eec706d544)](https://app.codacy.com/app/esse-gab2/dirstalk?utm_source=github.com&utm_medium=referral&utm_content=stefanoj3/dirstalk&utm_campaign=Badge_Grade_Dashboard)
[![Build Status](https://travis-ci.com/stefanoj3/dirstalk.svg?branch=master)](https://travis-ci.com/stefanoj3/dirstalk)
[![codecov](https://codecov.io/gh/stefanoj3/dirstalk/branch/master/graph/badge.svg)](https://codecov.io/gh/stefanoj3/dirstalk)

Dirstalk is a multi threaded application designed to brute force
directories and files names on web/application servers.

The idea is to create a tool with at least the same functionalities as
[dirbuster](https://www.owasp.org/index.php/Category:OWASP_DirBuster_Project)
and then expand it further.

[Golang](https://github.com/golang/go) is the language of choice for the
project.

## Milestones

### Version 1.0
- [x] Scan a given url
- [x] Specify how many threads to use
- [x] Specify http verbs to use
- [x] Specify dictionary to use
- [x] Specify http request timeout
- [x] Specify verbosity of the log
- [x] Specify scan depth
- [x] Socks5 support
- [ ] Can generate dictionary starting from a folder containing files/subdirectories
- [ ] Can be compiled for multiple platform/architectures (min: OSX x64, Linux 386, Linux x64, Linux arm, Linux arm64)
- [x] A CI is running the tests suite and display the code coverage
- [ ] Print results as a tree highlighting the statuses received
- [ ] Produce detailed documentation

### Someday
- [ ] Scan a website pages looking for links to bruteforce
- [ ] Expose a webserver that can be used to launch scans and check their status
- [ ] Display how much the scan has matched the dictionary (eg: how many entries had a match vs total of entries)
