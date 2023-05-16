BINARY_NAME=htg

build:
		@go build -o ${BINARY_NAME} -v

make run:
		@go run main.go
