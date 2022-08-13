build:
	go build -o ./bin/gocc ./cmd/gocc/main.go

test: build
	./test.sh

clean:
	rm -f ./bin/gocc ./bin/tmp*

.PHONY: test clean