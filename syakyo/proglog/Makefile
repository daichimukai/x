PROTO_FILES	=	$(wildcard api/v1/*.proto)
TARGETS		=	$(PROTO_FILES:.proto=.pb.go)

.PHONY: generate
generate: $(TARGETS)

%.pb.go: %.proto FORCE
	@protoc $<                             \
		--go_out=.                     \
		--go_opt=paths=source_relative \
		--proto_path=.

.PHONY: test
test:
	@go test -v -race ./...

.PHONY: FORCE
FORCE:
