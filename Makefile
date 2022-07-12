.ONESHELL: # Applies to every targets in the file! .ONESHELL instructs make to invoke a single instance of the shell and provide it with the entire recipe, regardless of how many lines it contains.
.SHELLFLAGS = -ec

########################################################################################################################
# Global Env Settings
########################################################################################################################

GO_WORKSPACE ?= /go/src/github.com/analogj/scrutiny

COLLECTOR_BINARY_NAME = scrutiny-collector-metrics
WEB_BINARY_NAME = scrutiny-web
LD_FLAGS =

STATIC_TAGS =
# enable multiarch docker image builds
DOCKER_TARGETARCH_BUILD_ARG =
ifdef TARGETARCH
DOCKER_TARGETARCH_BUILD_ARG := $(DOCKER_TARGETARCH_BUILD_ARG) --build-arg TARGETARCH=$(TARGETARCH)
endif

# enable to build static binaries.
ifdef STATIC
export CGO_ENABLED = 0
LD_FLAGS := $(LD_FLAGS) -extldflags=-static
STATIC_TAGS := $(STATIC_TAGS) -tags "static netgo"
endif
ifdef GOOS
COLLECTOR_BINARY_NAME := $(COLLECTOR_BINARY_NAME)-$(GOOS)
WEB_BINARY_NAME := $(WEB_BINARY_NAME)-$(GOOS)
LD_FLAGS := $(LD_FLAGS) -X main.goos=$(GOOS)
endif
ifdef GOARCH
COLLECTOR_BINARY_NAME := $(COLLECTOR_BINARY_NAME)-$(GOARCH)
WEB_BINARY_NAME := $(WEB_BINARY_NAME)-$(GOARCH)
LD_FLAGS := $(LD_FLAGS) -X main.goarch=$(GOARCH)
endif
ifdef GOARM
COLLECTOR_BINARY_NAME := $(COLLECTOR_BINARY_NAME)-$(GOARM)
WEB_BINARY_NAME := $(WEB_BINARY_NAME)-$(GOARM)
endif
ifeq ($(OS),Windows_NT)
COLLECTOR_BINARY_NAME := $(COLLECTOR_BINARY_NAME).exe
WEB_BINARY_NAME := $(WEB_BINARY_NAME).exe
endif

########################################################################################################################
# Binary
########################################################################################################################
.PHONY: all
all: binary-all

.PHONY: binary-all
binary-all: binary-collector binary-web
	@echo "built binary-collector and binary-web targets"


.PHONY: binary-clean
binary-clean:
	go clean

.PHONY: binary-dep
binary-dep:
	go mod vendor

.PHONY: binary-test
binary-test: binary-dep
	go test -v $(STATIC_TAGS) ./...

.PHONY: binary-test-coverage
binary-test-coverage: binary-dep
	go test -coverprofile=coverage.txt -covermode=atomic -v $(STATIC_TAGS) ./...

.PHONY: binary-collector
binary-collector: binary-dep
	go build -ldflags "$(LD_FLAGS)" -o $(COLLECTOR_BINARY_NAME) $(STATIC_TAGS) ./collector/cmd/collector-metrics/
ifneq ($(OS),Windows_NT)
	chmod +x $(COLLECTOR_BINARY_NAME)
	file $(COLLECTOR_BINARY_NAME) || true
	ldd $(COLLECTOR_BINARY_NAME) || true
	./$(COLLECTOR_BINARY_NAME) || true
endif

.PHONY: binary-web
binary-web: binary-dep
	go build -ldflags "$(LD_FLAGS)" -o $(WEB_BINARY_NAME) $(STATIC_TAGS) ./webapp/backend/cmd/scrutiny/
ifneq ($(OS),Windows_NT)
	chmod +x $(WEB_BINARY_NAME)
	file $(WEB_BINARY_NAME) || true
	ldd $(WEB_BINARY_NAME) || true
	./$(WEB_BINARY_NAME) || true
endif

########################################################################################################################
# Binary
########################################################################################################################

.PHONY: binary-frontend
# reduce logging, disable angular-cli analytics for ci environment
binary-frontend: export NPM_CONFIG_LOGLEVEL = warn
binary-frontend: export NG_CLI_ANALYTICS = false
binary-frontend:
	cd webapp/frontend
	npm install -g @angular/cli@9.1.4
	mkdir -p $(CURDIR)/dist
	npm ci
	npm run build:prod -- --output-path=$(CURDIR)/dist

.PHONY: binary-frontend-test-coverage
# reduce logging, disable angular-cli analytics for ci environment
binary-frontend-test-coverage:
	cd webapp/frontend
	npm ci
	npx ng test --watch=false --browsers=ChromeHeadless --code-coverage

########################################################################################################################
# Docker
# NOTE: these docker make targets are only used for local development (not used by Github Actions/CI)
# NOTE: docker-web and docker-omnibus require `make binary-frontend` or frontend.tar.gz content in /dist before executing.
########################################################################################################################
.PHONY: docker-collector
docker-collector:
	@echo "building collector docker image"
	docker build $(DOCKER_TARGETARCH_BUILD_ARG) -f docker/Dockerfile.collector -t analogj/scrutiny-dev:collector .

.PHONY: docker-web
docker-web:
	@echo "building web docker image"
	docker build $(DOCKER_TARGETARCH_BUILD_ARG) -f docker/Dockerfile.web -t analogj/scrutiny-dev:web .

.PHONY: docker-omnibus
docker-omnibus:
	@echo "building omnibus docker image"
	docker build $(DOCKER_TARGETARCH_BUILD_ARG) -f docker/Dockerfile -t analogj/scrutiny-dev:omnibus .
