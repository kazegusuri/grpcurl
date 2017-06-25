all: generate

PKG=github.com/kazegusuri/grpcurl
OUTPUT_DIR=_output

generate:
	mkdir -p $(OUTPUT_DIR)
	protoc  --go_out=plugins=grpc:$(OUTPUT_DIR) testdata/*.proto
	cp $(OUTPUT_DIR)/$(PKG)/testdata/*.go testdata/

clean:
	rm -f testdata/*.pb.go
