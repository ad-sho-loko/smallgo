smallgo: main.go
	go build -o smallgo

test: smallgo
	go test
	./test.sh

# run tests on my Mac
t: smallgo
	make clean            # for build on ubuntu os.
	./docker_test.sh

fmt:
	go fmt ./...

clean:
	rm -f *.s smallgo tmp.s tmp

.PHONY: t test clean fmt