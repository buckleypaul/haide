.PHONY: build test clean install

build:
	go build -o haide ./cmd/haide

test:
	go test -v ./...

clean:
	rm -f haide

install: build
	cp haide /usr/local/bin/

lint:
	golangci-lint run

fmt:
	go fmt ./...
