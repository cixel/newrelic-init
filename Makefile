.PHONY: test cover install

test:
	go test ./... -covermode=count -coverprofile=coverage.out

cover:
	go tool cover -html=coverage.out

install:
	go install .
