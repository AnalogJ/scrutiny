export CGO_ENABLED = 1

GO_WORKSPACE ?= /go/src/github.com/analogj/scrutiny

BINARY=\
  linux/amd64   \
  linux/arm-5   \
  linux/arm-6   \
  linux/arm-7   \
  linux/arm64   \


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

# clean:
# 	rm scrutiny-collector-metrics-* scrutiny-web-*
