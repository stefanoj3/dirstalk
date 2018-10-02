# Dirstalk
Dirstalk is a multi threaded application designed to brute force directories and files names on web/application servers.

The idea is to create a tool with at least the same functionalities as
[dirbuster](https://www.owasp.org/index.php/Category:OWASP_DirBuster_Project)
and then expand it further.

[Golang](https://github.com/golang/go) is the language of choice for the
project.


## Milestones

### Version 1.0
- [ ] Scan a given url
- [ ] Specify how many threads to use
- [ ] Specify http verbs to use
- [ ] Specify dictionary to use
- [ ] Specify http request timeout
- [ ] Specify verbosity of the log
- [ ] Socks5 support (with examples on how to use it with tor)
- [ ] Can generate dictionary starting from a folder containing files/subdirectories
- [ ] Can be compiled for multiple platform/architectures (min: OSX x64, Linux 386, Linux x64, Linux arm, Linux arm64)
- [ ] A CI is running the tests suite and display the code coverage

### Someday:
- [ ] Scan a website looking for links to bruteforce
- [ ] Expose a webserver that can be used to launch scans and check their status
- [ ] Display how much the scan has matched the dictionary (eg: how many entries had a match vs total of entries)
