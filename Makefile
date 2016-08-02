.PHONY=bin
bin: test
	go build -o bin/stdjson main/main.go

.PHONY=get-deps
get-deps:
	glide update

.PHONY=test
test:
	go list -f '{{if len .TestGoFiles }}"go test -v -coverprofile={{.Dir}}/.coverprofile {{.ImportPath}}"{{end}}' ./... | grep -v vendor | xargs -n 1 sh -c

fmt:
	go list ./... | grep -v vendor | xargs -n 1 go fmt
