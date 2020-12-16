# Forothree


Tools for bypassing 403 pages with build-in alteration and adding headers. Payload mostly generated based on multiple 403 bypass i found in twitter. Created using Go with Fasthttp module

## Features
  - Build in Concurrency
  - Default alteration for URL with a directory
  - Default wordlist containing headers to bypassing
  - Recursive alteration for URL with subdirectories (Experimental)


### Requirement

Forothree uses a [Fasthttp](https://github.com/valyala/fasthttp) to work properly

### Installation

```sh
$ go build forothree.go
$ ./forothree
```



### Usage

```sh
$ ./forothree -h
Usage of /tmp/go-build145687935/b001/exe/forothree:
  -b     disable header bypass
  -c     disable recursive bypass
  -e string
         set custom headers, ex head1:myhead,head2:yourhead (default "Connection:close")
  -hl
         show header location
  -l     show response length
  -m string
         set request method (default "GET")
  -o string
         specify output file name
  -r int
         set max number of retries (default 2)
  -s string
        -s specify status code, ex 200,404 (default "200,404,403,301,404")
  -t int
         specify request timeout in seconds (default 3)
  -u string
         url target
  -ul string
         url list target

```



### Example
```sh
$ ./forothree -u http://scanme.nmap.org/adminpage.php
domain : http://scanme.nmap.org/adminpage.php |code : 404 |
domain : http://scanme.nmap.org/.;adminpage.php |code : 404 |
domain : http://scanme.nmap.org/Adminpage.php |code : 404 |
domain : http://scanme.nmap.org/adminpage.php.json |code : 404 |
domain : http://scanme.nmap.org/adminpage.php# |code : 404 |
domain : http://scanme.nmap.org/adminpage.php/. |code : 404 |
domain : http://scanme.nmap.org/adminpage.php// |code : 404 |
domain : http://scanme.nmap.org/adminpage.php//dir@evil.com |code : 404 |
domain : http://scanme.nmap.org/adminpage.php/~ |code : 404 |
domain : http://scanme.nmap.org/adminpage.php%2500 |code : 404 |
domain : http://scanme.nmap.org/\..\.\adminpage.php |code : 404 |
domain : http://scanme.nmap.org/adminpage.php? |code : 404 |
domain : http://scanme.nmap.org/adminpage.php?? |code : 404 |
domain : http://scanme.nmap.org/adminpage.php%09 |code : 404 |
domain : http://scanme.nmap.org/adminpage.php/..;/ |code : 404 |
domain : http://scanme.nmap.org/adminpage.php..;/ |code : 404 |
domain : http://scanme.nmap.org/.;/adminpage.php |code : 404 |
domain : http://scanme.nmap.org/adminpage.php//google.com |code : 404 |
domain : http://scanme.nmap.org/..;/adminpage.php |code : 404 |
domain : http://scanme.nmap.org//adminpage.php |code : 404 |
domain : http://scanme.nmap.org/%2e/adminpage.php |code : 404 |
domain : http://scanme.nmap.org/adminpage.php%20 |code : 404 |
domain : http://scanme.nmap.org/./adminpage.php/./ |code : 404 |
domain : http://scanme.nmap.org/adminpage.php |code : 404 |xtra-header : X-Client-IP:127.0.0.1 |
domain : http://scanme.nmap.org/adminpage.php |code : 404 |xtra-header : CF-Connecting-IP:127.0.0.1 |
domain : http://scanme.nmap.org/adminpage.php |code : 404 |xtra-header : X-Cluster-Client-IP:127.0.0.1 |
domain : http://scanme.nmap.org/adminpage.php |code : 404 |xtra-header : X-Host:127.0.0.1 |
domain : http://scanme.nmap.org/adminpage.php |code : 404 |xtra-header : X-Originating-IP:127.0.0.1 |
domain : http://scanme.nmap.org/adminpage.php |code : 404 |xtra-header : X-Original-URL:/adminpage.php |
domain : http://scanme.nmap.org/adminpage.php |code : 404 |xtra-header : X-Real-IP:127.0.0.1 |
domain : http://scanme.nmap.org/adminpage.php |code : 404 |xtra-header : X-Forwarded-Host:127.0.0.1 |
domain : http://scanme.nmap.org/adminpage.php |code : 404 |xtra-header : Fastly-Client-IP:127.0.0.1 |
domain : http://scanme.nmap.org/adminpage.php |code : 404 |xtra-header : True-Client-IP:127.0.0.1 |
domain : http://scanme.nmap.org/adminpage.php |code : 404 |xtra-header : X-Client-IP:127.0.0.1 |
domain : http://scanme.nmap.org/adminpage.php |code : 404 |xtra-header : X-Custom-IP-Authorization:127.0.0.1 |
domain : http://scanme.nmap.org/adminpage.php |code : 404 |xtra-header : X-Remote-IP:127.0.0.1 |
domain : http://scanme.nmap.org/adminpage.php |code : 404 |xtra-header : X-Rewrite:/adminpage.php |
domain : http://scanme.nmap.org/adminpage.php |code : 404 |xtra-header : X-Forwarded-For:127.0.0.1 |
domain : http://scanme.nmap.org/adminpage.php |code : 404 |xtra-header : X-Forwarded-For:127.0.0.1 |
domain : http://scanme.nmap.org/adminpage.php |code : 404 |xtra-header : X-Remote-Addr:127.0.0.1 |
domain : http://scanme.nmap.org/adminpage.php |code : 404 |xtra-header : X-Forwarded-By:127.0.0.1 |
domain : http://scanme.nmap.org/adminpage.php |code : 404 |xtra-header : X-Forwarded-Proto:127.0.0.1 |
```
