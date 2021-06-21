.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')

RELEASE_VERSION = v1.0.0

APP 			= fusion-gin-admin
SERVER_BIN  	= ./${APP}
RELEASE_ROOT 	= release
RELEASE_SERVER 	= release/${APP}
GIT_COUNT 		= $(shell git rev-list --all --count)
GIT_HASH        = $(shell git rev-parse --short HEAD)
RELEASE_TAG     = $(RELEASE_VERSION).$(GIT_COUNT).$(GIT_HASH)

all: start

build:
	@go build -ldflags "-w -s -X main.VERSION=$(RELEASE_TAG)" -o $(SERVER_BIN) .

start:
	@go run -ldflags "-X main.VERSION=$(RELEASE_TAG)" ./main.go web -c ./config_file/config.toml -m ./config_file/model.conf --menu ./config_file/menu.yaml

swagger:
	@swag init --parseDependency --generalInfo ./main.go --output ./swagger

wire:
	@wire gen ./app

test:
	cd ./test && go test -v

clean:
	rm -rf data release $(SERVER_BIN) test/data test/data

pack: build
	rm -rf $(RELEASE_ROOT) && mkdir -p $(RELEASE_SERVER)
	cp -r $(SERVER_BIN) configs $(RELEASE_SERVER)
	cd $(RELEASE_ROOT) && tar -cvf $(APP).tar ${APP} && rm -rf ${APP}
