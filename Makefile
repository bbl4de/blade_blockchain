build:
	go build -o ./bin/projects

run: build
	./bin/projects

test:
	go test -v ./...