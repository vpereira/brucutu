build: deps
	go build -v -o build/brucutu
	
deps:
	go get -d -v ./...

test: deps
	go test -v ./...

# it runs in the docker-compose environment, runner container
integration: build
	./scripts/test_invalid_parameters.sh
	./scripts/test_valid_use_cases.sh