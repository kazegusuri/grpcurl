all: generate

PKG=github.com/kazegusuri/grpcurl
OUTPUT_DIR=_output

generate:
	mkdir -p $(OUTPUT_DIR)
	protoc -I. --go_out=plugins=grpc:$(OUTPUT_DIR) testdata/*.proto
	protoc -I. --go_out=plugins=grpc:$(OUTPUT_DIR) testdata/v2/*.proto
	cp $(OUTPUT_DIR)/$(PKG)/testdata/*.go testdata/
	cp $(OUTPUT_DIR)/$(PKG)/testdata/v2/*.go testdata/v2/

clean:
	rm -f testdata/*.pb.go

.PHONY: test
test:
	go test -v github.com/kazegusuri/grpcurl
