all: install

test:
	GO111MODULE=on go test ./... -bench=. -v -race

install:
	GO111MODULE=on go install
