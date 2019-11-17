all: deps
	go build -o build/brucutu
deps:
	go get -d -v ./...

test: deps
	go test -v ./..