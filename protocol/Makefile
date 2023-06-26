GOPATH=/go
PROTO_PATH=proto
REPO=/workspace/protocol
GO_PROTOCOL=go_protocol
TS_PROTOCOL=ts_protocol
LIVEKIT_VERSION=@v1.5.7
PROTOC_GEN_VERSION=@v1.0.1
PB_REL="https://github.com/protocolbuffers/protobuf/releases"

echo:
	@echo ${PWD}/${TS_PROTOCOL}
generate_go: 
	@protoc \
	-I ${GOPATH}/pkg/mod \
	-I ${GOPATH}/pkg/mod/github.com/livekit/protocol${LIVEKIT_VERSION} \
	-I ${GOPATH}/pkg/mod/github.com/envoyproxy/protoc-gen-validate${PROTOC_GEN_VERSION} \
	--proto_path=${PROTO_PATH} \
	${PROTO_PATH}/*.proto \
	--go_out=paths=source_relative:${GO_PROTOCOL} \
	--validate_out="lang=go,paths=source_relative:${GO_PROTOCOL}"

generate_ts:
	@protoc \
	-I ${GOPATH}/pkg/mod/github.com/livekit/protocol${LIVEKIT_VERSION} \
	--proto_path=${PROTO_PATH} \
	${PROTO_PATH}/*.proto \
	--ts_out=paths=source_relative:${PWD}/${TS_PROTOCOL} \
	 
generate: generate_go generate_ts

install:  manual_install
	sudo apt install -y protobuf-compiler 
	protoc --version
	sudo apt-get -y install protoc-gen-go

ts_install: 
	pnpm setup
	pnpm install -g protoc-gen-ts

manual_install: 
	curl -LO ${PB_REL}/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip
	unzip protoc-3.15.8-linux-x86_64.zip -d ${HOME}/.local
	export PATH="${PATH}:${HOME}/.local/bin"
	rm ./protoc-3.15.8-linux-x86_64.zip
	protoc --version
	go install github.com/envoyproxy/protoc-gen-validate
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	export PATH="${PATH}:${GOPATH}/bin"
	
version:
	protoc --version
	protoc-gen-validate --version
	protoc-gen-go --version

clear:
	rm -rf ${GO_PROTOCOL}/*
	rm -rf ${TS_PROTOCOL}/*
	mkdir -p ${GO_PROTOCOL}
	mkdir -p ${TS_PROTOCOL}