.PHONY=bin
bin: test
	go build -o bin/stdjson main/main.go

.PHONY=get-deps
get-deps:
	glide update

.PHONY=test
test:
	go test -v github.com/nkvoll/stdjson/config github.com/nkvoll/stdjson/rewriter github.com/nkvoll/stdjson/testutil github.com/nkvoll/stdjson