# Dirstalk
![](https://github.com/stefanoj3/dirstalk/workflows/CI/badge.svg)
[![codecov](https://codecov.io/gh/stefanoj3/dirstalk/branch/master/graph/badge.svg)](https://codecov.io/gh/stefanoj3/dirstalk)
[![Scrutinizer Code Quality](https://scrutinizer-ci.com/g/stefanoj3/dirstalk/badges/quality-score.png?b=master)](https://scrutinizer-ci.com/g/stefanoj3/dirstalk/?branch=master)
![Docker Pulls](https://img.shields.io/docker/pulls/stefanoj3/dirstalk.svg)
![GitHub](https://img.shields.io/github/license/stefanoj3/dirstalk.svg)

Dirstalk is a multi threaded application designed to brute force paths on web servers.

The tool contains functionalities similar to the ones offered by
[dirbuster](https://www.owasp.org/index.php/Category:OWASP_DirBuster_Project) 
and [dirb](https://tools.kali.org/web-applications/dirb).

Here you can see it in action:
[![asciicast](https://asciinema.org/a/ehvNAUetjWbNExQegA2KPaHuY.svg)](https://asciinema.org/a/ehvNAUetjWbNExQegA2KPaHuY)

## Contents
- [How to use it](#-how-to-use-it)
    - [Scan](#scan)
    - [Useful resources](#useful-resources)
    - [Dictionary generator](#dictionary-generator)
- [Download](#-download)
- [Development](#-development)
- [License](https://github.com/stefanoj3/dirstalk/blob/master/LICENSE.md)

## [↑](#contents) How to use it

The application is self-documenting, launching `dirstalk -h` will return all the available commands with a 
short description, you can get the help for each command by doing `distalk <command> -h`.

EG `dirstalk result.diff -h`

### Scan

To perform a scan you need to provide at least a dictionary and a URL:
```shell script
dirstalk scan http://someaddress.url/ --dictionary mydictionary.txt
```

As mentioned before, to see all the flags available for the scan command you can 
just call the command with the `-h` flag:
```shell script
dirstalk scan -h
```

##### Example of how you can customize a scan:
```shell script
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


##### Currently available flags:
```shell script
      --cookie stringArray             cookie to add to each request; eg name=value (can be specified multiple times)
  -d, --dictionary string              dictionary to use for the scan (path to local file or remote url)
      --header stringArray             header to add to each request; eg name=value (can be specified multiple times)
  -h, --help                           help for scan
      --http-cache-requests            cache requests to avoid performing the same request multiple times within the same scan (EG if the server reply with the same redirect location multiple times, dirstalk will follow it only once) (default true)
      --http-methods strings           comma separated list of http methods to use; eg: GET,POST,PUT (default [GET])
      --http-statuses-to-ignore ints   comma separated list of http statuses to ignore when showing and processing results; eg: 404,301 (default [404])
      --http-timeout int               timeout in milliseconds (default 5000)
      --out string                     path where to store result output
      --scan-depth int                 scan depth (default 3)
      --socks5 string                  socks5 host to use
  -t, --threads int                    amount of threads for concurrent requests (default 3)
      --use-cookie-jar                 enables the use of a cookie jar: it will retain any cookie sent from the server and send them for the following requests
      --user-agent string              user agent to use for http requests
```

##### Useful resources
- [here](https://github.com/dustyfresh/dictionaries/tree/master/DirBuster-Lists) you can find dictionaries that can be used with dirstalk
- [tordock](https://github.com/stefanoj3/tordock) is a containerized Tor SOCKS5 that you can use easily with dirstalk 
(just `docker run -d -p 127.0.0.1:9150:9150 stefanoj3/tordock:latest` and then when launching a
 scan specify the following flag: `--socks5 127.0.0.1:9150`)

### Dictionary generator
Dirstalk can also produce it's own dictionaries, useful for example if you
want to check if a specific set of files is available on a given web server.

##### Example:
```shell script
dirstalk dictionary.generate /path/to/local/files --out mydictionary.txt
```
The result will be printed to the stdout if no out flag is specified.

## [↑](#contents) Download
You can download a release from [here](https://github.com/stefanoj3/dirstalk/releases)
or you can use a docker image. (eg `docker run stefanoj3/dirstalk dirstalk <cmd>`)

If you are using an arch based linux distribution you can fetch it via AUR: https://aur.archlinux.org/packages/dirstalk/

Example:
```shell script
yay -S aur/dirstalk
```


## [↑](#contents) Development
All you need to do local development is to have [make](https://www.gnu.org/software/make/)
and [golang](https://golang.org/) available and the GOPATH correctly configured.

Then you can just clone the project, enter the folder and:
```shell script
make dep                                     # to fetch dependencies
make tests                                   # to run the test suite
make check                                   # to check for any code style issue
make fix                                     # to automatically fix the code style using goimports
make build                                   # to build an executable for your host OS (not tested under windows) 
```

```shell script
make help
```
will print a description of every command available in the Makefile.

Wanna add a functionality? fix a bug? fork and create a PR.

## [↑](#contents) Plans for the future
- Add support for rotating SOCKS5 proxies
- Scan a website pages looking for links to bruteforce
- Expose a webserver that can be used to launch scans and check their status
- Introduce metrics that can give a sense of how much of the dictionary was found on the remote server
