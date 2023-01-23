.PHONY: build
build: clean
	CGO_ENABLED=0 go build -o becrypt

.PHONY: clean
clean:
	rm -f becrypt

.PHONY: test
test:
	go test -race ./...

.PHONY: vet
vet:
	go vet ./...

default: build
