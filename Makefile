PKG=github.com/kazegusuri/grpcurl
OUTPUT_DIR=_output

all: build

dep:
	go mod download

build:
	go build -o grpcurl .

generate:
	mkdir -p $(OUTPUT_DIR)
	protoc -I. -Ithird_party/googleapis --go_out=plugins=grpc:$(OUTPUT_DIR) internal/testdata/*.proto
	protoc -I. -Ithird_party/googleapis --go_out=plugins=grpc:$(OUTPUT_DIR) internal/testdata/v2/*.proto
	cp $(OUTPUT_DIR)/$(PKG)/internal/testdata/*.go internal/testdata/
	cp $(OUTPUT_DIR)/$(PKG)/internal/testdata/v2/*.go internal/testdata/v2/

clean:
	rm -f internal/testdata/*.pb.go

.PHONY: test
test:
	go test -v github.com/kazegusuri/grpcurl
