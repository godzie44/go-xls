.PHONY: lint test

lint:
	CGO_LDFLAGS=-lxlsreader golangci-lint run ;

test:
	CGO_LDFLAGS=-lxlsreader go test ./... -race -count=1 ;
