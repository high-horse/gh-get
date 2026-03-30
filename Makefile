.PHONY: build build-static

build:
	go build -o build/local/gh-get src/*.go

build-static:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/linux/gh-get src/*.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/darwin/gh-get src/*.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o build/darwin/gh-get-arm64 src/*.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/windows/gh-get.exe src/*.go
