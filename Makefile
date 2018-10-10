all: install

test:
	go test ./... -bench=. -v -race

install:
	GO111MODULE=on go install
