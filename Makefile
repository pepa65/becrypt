.PHONY: build
build: clean
	CGO_ENABLED=0 go build -o becrypt

.PHONY: clean
clean:
	rm -f becrypt

.PHONY: vet
vet:
	go vet ./...

default: build
