# build go service
build:
	go build -o ./bin/api ./cmd/api

# run go service
run: clean build
	./bin/api

# unit test go service
test:
	go test ./...

# clean service binary 
clean:
	go clean
	rm -rf ./bin

