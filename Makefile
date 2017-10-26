all: install

test: deps
	go test ./bot -bench=. -v -race

install: deps
	go install

deps:
	go get github.com/golang/dep/cmd/dep
	dep ensure
