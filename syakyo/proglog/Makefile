CONFIG_PATH=${HOME}/.proglog/

.PHONY: generate
generate: $(TARGETS)
	@protoc api/v1/*.proto                      \
		--go_out=.                          \
		--go-grpc_out=.                     \
		--go_opt=paths=source_relative      \
		--go-grpc_opt=paths=source_relative \
		--proto_path=.

.PHONY: test
test: generate-cert ${CONFIG_PATH}/model.conf ${CONFIG_PATH}/policy.csv
	@go test -v -race ./...

.PHONY: generate-cert
generate-cert:
	@cfssl gencert -initca test/ca-csr.json 2>/dev/null | cfssljson -bare ca
	@cfssl gencert                      \
		-ca=ca.pem                  \
		-ca-key=ca-key.pem          \
		-config=test/ca-config.json \
		-profile=server             \
		test/server-csr.json 2>/dev/null | cfssljson -bare server
	@cfssl gencert                      \
		-ca=ca.pem                  \
		-ca-key=ca-key.pem          \
		-config=test/ca-config.json \
		-profile=client             \
		-cn=root                    \
		test/client-csr.json 2>/dev/null | cfssljson -bare root-client
	@cfssl gencert                      \
		-ca=ca.pem                  \
		-ca-key=ca-key.pem          \
		-config=test/ca-config.json \
		-profile=client             \
		-cn=nobody                  \
		test/client-csr.json 2>/dev/null | cfssljson -bare nobody-client
	@mv *.pem *.csr ${CONFIG_PATH}

${CONFIG_PATH}/model.conf: test/model.conf
	@cp $< $@

${CONFIG_PATH}/policy.csv: test/policy.csv
	@cp $< $@
