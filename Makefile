BINARY_NAME=discord-dl

build:
	go build -o bin/${BINARY_NAME}

run:
	go build -o bin/${BINARY_NAME}
	./bin/${BINARY_NAME}
clean:
	go clean
	rm bin/${BINARY_NAME}
	rm bigbrother.db
	rm -r media/*