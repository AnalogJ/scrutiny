.ONESHELL: # Applies to every targets in the file! .ONESHELL instructs make to invoke a single instance of the shell and provide it with the entire recipe, regardless of how many lines it contains.
.SHELLFLAGS = -ec
export GOTOOLCHAIN=go1.25.5

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

.PHONY: lint
lint:
	GOTOOLCHAIN=go1.25.5 go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.8.0
	golangci-lint run ./...

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
# Debian package
########################################################################################################################

SCRUTINY_VERSION = $(shell sed -n 's/^const VERSION = "\(.*\)"/\1/p' webapp/backend/pkg/version/version.go)
# NIGHTLY builds sort after the last release and before the next one
DEB_VERSION = $(SCRUTINY_VERSION)$(if $(NIGHTLY),+git$(shell git log -1 --format=%cd --date=format:%Y%m%d).$(shell git rev-parse --short=7 HEAD))
DEB_GOARCH = $(or $(GOARCH),$(shell go env GOARCH))
DEB_ARCH = $(if $(filter arm,$(DEB_GOARCH)),$(if $(filter 7,$(GOARM)),armhf,$(if $(filter 5,$(GOARM)),armel,$(error GOARM must be 5 (armel) or 7 (armhf)))),$(DEB_GOARCH))
# PKG is set per package target
DEB_ROOT = build/deb/$(PKG)_$(DEB_ARCH)
# reproducible: static, no local paths or VCS state embedded
DEB_GO_BUILD = $(MAKE) STATIC=1 GOOS=linux GOARCH=$(DEB_GOARCH) GOFLAGS="-trimpath -buildvcs=false"

define deb-start
	rm -rf $(DEB_ROOT)
	mkdir -p $(DEB_ROOT)/usr/bin $(DEB_ROOT)/usr/share/man/man8 $(DEB_ROOT)/etc/scrutiny $(DEB_ROOT)/DEBIAN
endef

define deb-finish
	install -Dm644 LICENSE $(DEB_ROOT)/usr/share/doc/$(PKG)/copyright
	cp $(addprefix packaging/deb/$(PKG)/,conffiles postinst prerm postrm) $(DEB_ROOT)/DEBIAN/
	INSTALLED_SIZE=$$(find $(DEB_ROOT) -path $(DEB_ROOT)/DEBIAN -prune -o -type f -printf '%s\n' | awk '{ kb += int(($$1 + 1023) / 1024) } END { print kb }')
	sed -e 's/@VERSION@/$(DEB_VERSION)/' -e 's/@ARCH@/$(DEB_ARCH)/' -e "s/@INSTALLED_SIZE@/$$INSTALLED_SIZE/" \
		packaging/deb/$(PKG)/control > $(DEB_ROOT)/DEBIAN/control
	# keep the umask and checkout modes out of the package
	chmod -R u=rwX,go=rX $(DEB_ROOT)
	dpkg-deb --root-owner-group -Zxz --build $(DEB_ROOT) $(PKG)_$(DEB_ARCH).deb
endef

package-deb-collector package-deb-web: export SOURCE_DATE_EPOCH ?= $(shell git log -1 --format=%ct)

.PHONY: package-deb-collector
package-deb-collector: PKG = scrutiny-collector
package-deb-collector:
	$(deb-start)
	$(DEB_GO_BUILD) binary-collector COLLECTOR_BINARY_NAME=$(DEB_ROOT)/usr/bin/scrutiny-collector-metrics
	install -Dm644 -t $(DEB_ROOT)/usr/lib/systemd/system packaging/systemd/scrutiny-collector.service packaging/systemd/scrutiny-collector.timer
	sed -e 's/^#api:$$/api:/' -e "s|^#  endpoint: 'http://localhost:8080'$$|  endpoint: 'http://localhost:8080'|" \
		example.collector.yaml > $(DEB_ROOT)/etc/scrutiny/collector.yml
	grep -qx "  endpoint: 'http://localhost:8080'" $(DEB_ROOT)/etc/scrutiny/collector.yml
	gzip -9n < docs/man/scrutiny-collector-metrics.8 > $(DEB_ROOT)/usr/share/man/man8/scrutiny-collector-metrics.8.gz
	$(deb-finish)

.PHONY: package-deb-web
package-deb-web: PKG = scrutiny-web
package-deb-web: dist
	$(deb-start)
	$(DEB_GO_BUILD) binary-web WEB_BINARY_NAME=$(DEB_ROOT)/usr/bin/scrutiny-web
	install -Dm644 -t $(DEB_ROOT)/usr/lib/systemd/system packaging/systemd/scrutiny-web.service
	mkdir -p $(DEB_ROOT)/usr/share/scrutiny
	cp -r dist $(DEB_ROOT)/usr/share/scrutiny/web
	sed -e 's#location: /opt/scrutiny/config/scrutiny.db#location: /var/lib/scrutiny/scrutiny.db#' \
		-e 's#path: /opt/scrutiny/web#path: /usr/share/scrutiny/web#' \
		example.scrutiny.yaml > $(DEB_ROOT)/etc/scrutiny/scrutiny.yaml
	grep -q 'location: /var/lib/scrutiny/scrutiny.db' $(DEB_ROOT)/etc/scrutiny/scrutiny.yaml
	grep -q 'path: /usr/share/scrutiny/web' $(DEB_ROOT)/etc/scrutiny/scrutiny.yaml
	gzip -9n < docs/man/scrutiny-web.8 > $(DEB_ROOT)/usr/share/man/man8/scrutiny-web.8.gz
	$(deb-finish)

# reuses an existing frontend build, e.g. one extracted from CI
dist:
	$(MAKE) binary-frontend

########################################################################################################################
# Binary
########################################################################################################################

.PHONY: binary-frontend
# reduce logging, disable angular-cli analytics for ci environment
binary-frontend: export NPM_CONFIG_LOGLEVEL = warn
binary-frontend: export NG_CLI_ANALYTICS = false
binary-frontend:
	cd webapp/frontend
	npm install -g @angular/cli@v13-lts
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
########################################################################################################################
.PHONY: docker-collector
docker-collector:
	@echo "building collector docker image"
	docker build $(DOCKER_TARGETARCH_BUILD_ARG) -f docker/Dockerfile.collector -t ghcr.io/analogj/scrutiny-dev:collector .

.PHONY: docker-web
docker-web:
	@echo "building web docker image"
	docker build $(DOCKER_TARGETARCH_BUILD_ARG) -f docker/Dockerfile.web -t ghcr.io/analogj/scrutiny-dev:web .

.PHONY: docker-omnibus
docker-omnibus:
	@echo "building omnibus docker image"
	docker build $(DOCKER_TARGETARCH_BUILD_ARG) -f docker/Dockerfile -t ghcr.io/analogj/scrutiny-dev:omnibus .
