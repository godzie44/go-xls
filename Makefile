.PHONY: lint test run-example

lint:
	CGO_LDFLAGS=-lxlsreader golangci-lint run ;

test:
	CGO_LDFLAGS=-lxlsreader go test ./... -race -count=1 ;

example:
	docker build -f example/Dockerfile -t xls-ex . && \
	docker run xls-ex > table.html ;
