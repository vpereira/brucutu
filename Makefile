.PHONY: build deps test integration

build: 
	go build -v -o build/brucutu
deps:
	go mod tidy
test: 
	go test -v ./...

# it runs in the docker-compose environment, runner container
integration: build
	./scripts/test_invalid_parameters.sh
	./scripts/test_valid_use_cases.sh
