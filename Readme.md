## stdjson

[![Build Status](https://travis-ci.org/nkvoll/stdjson.svg?branch=master)](https://travis-ci.org/nkvoll/stdjson) [![Coverage Status](https://coveralls.io/repos/github/nkvoll/stdjson/badge.svg?branch=master)](https://coveralls.io/github/nkvoll/stdjson?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/nkvoll/stdjson)](https://goreportcard.com/report/github.com/nkvoll/stdjson) [![codebeat badge](https://codebeat.co/badges/afc287d8-e9fe-47f9-9b61-db14b022604f)](https://codebeat.co/projects/github-com-nkvoll-stdjson)

A wrapper for a subprocess' standard out / standard err that converts the output to JSON.

**Note**: This project is not intended to be production grade at this time, but rather to serve as an small project for the author to get to know [Go](https://golang.org/).

### Features

- Capture groups using regex or grok patterns.
- Wrap stdout / stderr in different rewriters / different configurations.
- Multiline joining, based on continuation prefixes and a configurable timeout.
- Recursively match / dice up fields.
- Adding default fields (arbitrary key/values) for every emitted output object.

### Building

Checkout to a proper tree:

    $ git checkout .. $GOPATH/github.com/nkvoll/stdjson
    $ cd $GOPATH/github.com/nkvoll/stdjson
    
Get vendored dependencies:
    
    $ glide install

Test and build:
    
    $ make test
    $ make bin
    
### Running

``` json
$ stdjson -config examples/ls-rewriter.yaml -- ls -al
{"line":"total ..."}
...
{"group":"staff","links":1,"name":"Readme.md","perms":"-rw-r--r--","size":725,"time":"Aug  2 19:43","time":"...","user":"njal"}
{"group":"staff","links":3,"name":"vendor","perms":"drwxr-xr-x","size":102,"time":"Aug  2 13:31","time":"...","user":"user"}
...
```
``` console
$ stdjson -config examples/noop.yaml -- ls -al
total ...
...
-rw-r--r--   1 user  staff     725 Aug  2 19:43 Readme.md
drwxr-xr-x   3 user  staff     102 Aug  2 13:31 vendor
...
```

### License

MIT
