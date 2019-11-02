smallgo: main.go
	go build -o smallgo *.go

test: smallgo
	./test.sh

test_all : smallgo
	go test ./...
	./test.sh

# run tests on my Mac
t: smallgo
	make clean # for build on ubuntu os
	./docker_test.sh
	make clean # remove binaries if tests succeed

fmt:
	go fmt ./...

clean:
	go clean
	rm -f *.s smallgo tmp.s tmp main

.PHONY: t test clean fmt