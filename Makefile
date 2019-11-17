build: deps
	go build -v -o build/brucutu
	
deps:
	go get -d -v ./...

test: deps
	go test -v

# it runs in the docker-compose environment, runner container
integration: build
	./build/brucutu -u ssh://ssh -a 2222 -l root -p superpassword || exit -1
	./build/brucutu -u ssh://ssh -a 2222 -L sample/users.txt -P sample/passwd.txt || exit -1
	./build/brucutu -u pop3://email -l foo -p thepassword || exit -1
	./build/brucutu -u pop3://email -L sample/users.txt -P sample/passwd.txt || exit -1

