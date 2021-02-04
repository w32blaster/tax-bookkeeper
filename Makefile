.PHONY: test
test:
	go vet ./...
	go test -race -short ./...

build: test
	CGO_ENABLED=0 go build \
        -a -installsuffix cgo \
        -ldflags "-s -w" \
        -o tax-bookkeeper \
        cmd/bookkeeper/main.go