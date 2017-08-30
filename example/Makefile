SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

example: $(SRC)
	env GOOS=linux GOARCH=amd64 go build

proto:
	protoc -I grpc/ grpc/secret.proto --go_out=plugins=grpc:grpc
