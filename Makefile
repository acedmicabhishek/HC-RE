.PHONY: build run clean test

build:
	go build -o bin/hc-re cmd/main.go

run: build
	./bin/hc-re -rps 10000 -duration 30s -model model1 -cpu 1000

model1: build
	./bin/hc-re -rps 100000 -duration 10s -model model1 -cpu 100

model2: build
	./bin/hc-re -rps 100000 -duration 10s -model model2 -workers 100 -cpu 100

clean:
	rm -rf bin/

test:
	go test -v ./...

fmt:
	go fmt ./...

vet:
	go vet ./...
