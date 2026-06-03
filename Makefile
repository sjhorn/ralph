.PHONY: build clean install

build:
	go build -o ralph ./cmd/ralph/

clean:
	rm -f ralph

install:
	go install ./cmd/ralph/
