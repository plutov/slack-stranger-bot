all: install

test: deps
	go test ./bot -bench=. -v -race

install: deps
	go install

deps:
	dep ensure
