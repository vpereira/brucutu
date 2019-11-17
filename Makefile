all: deps
	go build -v -o build/brucutu
deps:
	go get -d -v ./...

test: deps
	go test -v 