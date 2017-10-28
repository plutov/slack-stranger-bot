all: install

test: deps
	go test ./... -bench=. -v -race

install: deps
	go install

deps:
	go get github.com/golang/dep/cmd/dep
	dep ensure
