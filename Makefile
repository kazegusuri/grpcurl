all: generate

PKG=github.com/kazegusuri/grpcurl
OUTPUT_DIR=_output

generate:
	mkdir -p $(OUTPUT_DIR)
	protoc -I. --go_out=plugins=grpc:$(OUTPUT_DIR) internal/testdata/*.proto
	protoc -I. --go_out=plugins=grpc:$(OUTPUT_DIR) internal/testdata/v2/*.proto
	cp $(OUTPUT_DIR)/$(PKG)/internal/testdata/*.go internal/testdata/
	cp $(OUTPUT_DIR)/$(PKG)/internal/testdata/v2/*.go internal/testdata/v2/

clean:
	rm -f internal/testdata/*.pb.go

.PHONY: test
test:
	go test -v github.com/kazegusuri/grpcurl
