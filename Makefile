build:
	go build -o gh-get src/*.go

build-static: 
	CGO_ENABLED=0 go build -o build/linux/gh-get src/*.go
