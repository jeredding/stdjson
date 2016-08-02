## stdjson

[![Build Status](https://travis-ci.org/nkvoll/stdjson.svg?branch=master)](https://travis-ci.org/nkvoll/stdjson) [![Coverage Status](https://coveralls.io/repos/github/nkvoll/stdjson/badge.svg?branch=master)](https://coveralls.io/github/nkvoll/stdjson?branch=master)

A wrapper for a subprocess' standard out / standard err that converts the output to JSON.

**Note**: This project is not intended to be production grade at this time, but rather to serve as an small project for the author to get to know [Go](https://golang.org/).

### Features

- Capture groups using regex or grok patterns.
- Wrap stdout / stderr in different rewriters / different configurations.
- Multiline joining, based on continuation prefixes and a configurable timeout.
- Recursively match / dice up fields.

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

    $ stdjson -config examples/ls-rewriter.yaml -- ls -al
    {"line":"total ..."}
    ...
    {"group":"staff","links":1,"name":"Readme.md","perms":"-rw-r--r--","size":725,"time":"Aug  2 19:43","time":"...",user":"njal"}
    {"group":"staff","links":3,"name":"vendor","perms":"drwxr-xr-x","size":102,"time":"Aug  2 13:31","time":"...","user":"user"}
    ...
    $ stdjson -config examples/noop.yaml -- ls -al
    total ...
    ...
    -rw-r--r--   1 user  staff     725 Aug  2 19:43 Readme.md
    drwxr-xr-x   3 user  staff     102 Aug  2 13:31 vendor
    ...

### License

MIT
