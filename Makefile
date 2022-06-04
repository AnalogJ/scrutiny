export CGO_ENABLED = 1

GO_WORKSPACE ?= /go/src/github.com/analogj/scrutiny

BINARY=\
  linux/amd64   \
  linux/arm-5   \
  linux/arm-6   \
  linux/arm-7   \
  linux/arm64   \

.ONESHELL: # Applies to every targets in the file! .ONESHELL instructs make to invoke a single instance of the shell and provide it with the entire recipe, regardless of how many lines it contains.
.PHONY: all $(BINARY)
all: $(BINARY) windows/amd64

$(BINARY): OS = $(word 1,$(subst /, ,$*))
$(BINARY): ARCH = $(word 2,$(subst /, ,$*))
$(BINARY): build/scrutiny-web-%:
	@echo "building web binary (OS = $(OS), ARCH = $(ARCH))"
	xgo -v --targets="$(OS)/$(ARCH)" -ldflags "-extldflags=-static -X main.goos=$(OS) -X main.goarch=$(ARCH)" -out scrutiny-web -tags "static netgo sqlite_omit_load_extension" ${GO_WORKSPACE}/webapp/backend/cmd/scrutiny/

	chmod +x "/build/scrutiny-web-$(OS)-$(ARCH)"
	file "/build/scrutiny-web-$(OS)-$(ARCH)" || true
	ldd "/build/scrutiny-web-$(OS)-$(ARCH)" || true

	@echo "building collector binary (OS = $(OS), ARCH = $(ARCH))"
	xgo -v --targets="$(OS)/$(ARCH)" -ldflags "-extldflags=-static -X main.goos=$(OS) -X main.goarch=$(ARCH)" -out scrutiny-collector-metrics -tags "static netgo" ${GO_WORKSPACE}/collector/cmd/collector-metrics/

	chmod +x "/build/scrutiny-collector-metrics-$(OS)-$(ARCH)"
	file "/build/scrutiny-collector-metrics-$(OS)-$(ARCH)" || true
	ldd "/build/scrutiny-collector-metrics-$(OS)-$(ARCH)" || true

windows/amd64: export OS = windows
windows/amd64: export ARCH = amd64
windows/amd64:
	@echo "building web binary (OS = $(OS), ARCH = $(ARCH))"
	xgo -v --targets="$(OS)/$(ARCH)" -ldflags "-extldflags=-static -X main.goos=$(OS) -X main.goarch=$(ARCH)" -out scrutiny-web -tags "static netgo sqlite_omit_load_extension" ${GO_WORKSPACE}/webapp/backend/cmd/scrutiny/

	@echo "building collector binary (OS = $(OS), ARCH = $(ARCH))"
	xgo -v --targets="$(OS)/$(ARCH)" -ldflags "-extldflags=-static -X main.goos=$(OS) -X main.goarch=$(ARCH)" -out scrutiny-collector-metrics -tags "static netgo" ${GO_WORKSPACE}/collector/cmd/collector-metrics/


docker-collector:
	@echo "building collector docker image"
	docker build --build-arg TARGETARCH=amd64 -f docker/Dockerfile.collector -t analogj/scrutiny-dev:collector .

docker-web:
	@echo "building web docker image"
	docker build --build-arg TARGETARCH=amd64 -f docker/Dockerfile.web -t analogj/scrutiny-dev:web .

docker-omnibus:
	@echo "building omnibus docker image"
	docker build --build-arg TARGETARCH=amd64 -f docker/Dockerfile -t analogj/scrutiny-dev:omnibus .

# reduce logging, disable angular-cli analytics for ci environment
frontend: export NPM_CONFIG_LOGLEVEL = warn
frontend: export NG_CLI_ANALYTICS = false
frontend:
	cd webapp/frontend
	npm install -g @angular/cli@9.1.4
	mkdir -p $(CURDIR)/dist
	npm install
	npm run build:prod -- --output-path=$(CURDIR)/dist

# clean:
# 	rm scrutiny-collector-metrics-* scrutiny-web-*
