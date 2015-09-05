ifaddr: *.go
	go build -ldflags "-X main.Version=$(shell cat VERSION) -X main.CommitHash=$(shell git rev-parse HEAD)"

clean:
	rm -f ifaddr

.PHONY: clean
