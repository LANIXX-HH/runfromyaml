# build the Go binary
build:
	go build
	cp runfromyaml ${HOME}/bin/

# run tests
test:
	go test -v ./...

# install dependencies
deps:
	go get -v ./...
	go mod download

# clean the project
clean:
	rm -rf runfromyaml
	rm -rf ${HOME}/bin/runfromyaml
