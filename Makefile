export CGO_ENABLED = 1


BINARY=\
  linux/amd64   \
  linux/arm-5   \
  linux/arm-6   \
  linux/arm-7   \
  linux/arm64   \


.PHONY: all $(BINARY)
all: $(BINARY) windows/amd64 freebsd/amd64 freebsd/386

$(BINARY): OS = $(word 1,$(subst /, ,$*))
$(BINARY): ARCH = $(word 2,$(subst /, ,$*))
$(BINARY): build/scrutiny-web-%:
	@echo "building web binary (OS = $(OS), ARCH = $(ARCH))"
	xgo -v --targets="$(OS)/$(ARCH)" -ldflags "-extldflags=-static -X main.goos=$(OS) -X main.goarch=$(ARCH)" -out scrutiny-web -tags "static netgo sqlite_omit_load_extension" ./webapp/backend/cmd/scrutiny/

	chmod +x "/build/scrutiny-web-$(OS)-$(ARCH)"
	file "/build/scrutiny-web-$(OS)-$(ARCH)" || true
	ldd "/build/scrutiny-web-$(OS)-$(ARCH)" || true

	@echo "building collector binary (OS = $(OS), ARCH = $(ARCH))"
	xgo -v --targets="$(OS)/$(ARCH)" -ldflags "-extldflags=-static -X main.goos=$(OS) -X main.goarch=$(ARCH)" -out scrutiny-collector-metrics -tags "static netgo" ./collector/cmd/collector-metrics/

	chmod +x "/build/scrutiny-collector-metrics-$(OS)-$(ARCH)"
	file "/build/scrutiny-collector-metrics-$(OS)-$(ARCH)" || true
	ldd "/build/scrutiny-collector-metrics-$(OS)-$(ARCH)" || true

windows/amd64: OS = windows
windows/amd64: ARCH = amd64
windows/amd64:
	@echo "building web binary (OS = $(OS), ARCH = $(ARCH))"
	xgo -v --targets="$(OS)/$(ARCH)" -ldflags "-extldflags=-static -X main.goos=$(OS) -X main.goarch=$(ARCH)" -out scrutiny-web -tags "static netgo sqlite_omit_load_extension" ./webapp/backend/cmd/scrutiny/

	@echo "building collector binary (OS = $(OS), ARCH = $(ARCH))"
	xgo -v --targets="$(OS)/$(ARCH)" -ldflags "-extldflags=-static -X main.goos=$(OS) -X main.goarch=$(ARCH)" -out scrutiny-collector-metrics -tags "static netgo" ./collector/cmd/collector-metrics/

freebsd/amd64: OS = freebsd
freebsd/amd64: ARCH = amd64
freebsd/amd64: GOOS = freebsd
freebsd/amd64: GOARCH = amd64
freebsd/amd64: CGO_ENABLED = 1
freebsd/amd64:
	mkdir -p /build

	@echo "building web binary (OS = $(OS), ARCH = $(ARCH))"
	go build -ldflags "-extldflags=-static -X main.goos=$(OS) -X main.goarch=$(ARCH)" -o /build/scrutiny-web-$(OS)-$(ARCH) -tags "static netgo sqlite_omit_load_extension" webapp/backend/cmd/scrutiny/scrutiny.go

	chmod +x "/build/scrutiny-web-$(OS)-$(ARCH)"
	file "/build/scrutiny-web-$(OS)-$(ARCH)" || true
	ldd "/build/scrutiny-web-$(OS)-$(ARCH)" || true

	@echo "building collector binary (OS = $(OS), ARCH = $(ARCH))"
	go build -ldflags "-extldflags=-static -X main.goos=$(OS) -X main.goarch=$(ARCH)" -o /build/scrutiny-collector-metrics-$(OS)-$(ARCH) -tags "static netgo" collector/cmd/collector-metrics/collector-metrics.go

	chmod +x "/build/scrutiny-collector-metrics-$(OS)-$(ARCH)"
	file "/build/scrutiny-collector-metrics-$(OS)-$(ARCH)" || true
	ldd "/build/scrutiny-collector-metrics-$(OS)-$(ARCH)" || true

freebsd/amd64: OS = freebsd
freebsd/amd64: ARCH = 386
freebsd/amd64: GOOS = freebsd
freebsd/amd64: GOARCH = 386
freebsd/amd64: CGO_ENABLED = 1
freebsd/amd64:
	mkdir -p /build
    env
	@echo "building web binary (OS = $(OS), ARCH = $(ARCH))"
	go build -ldflags "-extldflags=-static -X main.goos=$(OS) -X main.goarch=$(ARCH)" -o /build/scrutiny-web-$(OS)-$(ARCH) -tags "static netgo sqlite_omit_load_extension" webapp/backend/cmd/scrutiny/scrutiny.go

	chmod +x "/build/scrutiny-web-$(OS)-$(ARCH)"
	file "/build/scrutiny-web-$(OS)-$(ARCH)" || true
	ldd "/build/scrutiny-web-$(OS)-$(ARCH)" || true

	@echo "building collector binary (OS = $(OS), ARCH = $(ARCH))"
	go build -ldflags "-extldflags=-static -X main.goos=$(OS) -X main.goarch=$(ARCH)" -o /build/scrutiny-collector-metrics-$(OS)-$(ARCH) -tags "static netgo" collector/cmd/collector-metrics/collector-metrics.go

	chmod +x "/build/scrutiny-collector-metrics-$(OS)-$(ARCH)"
	file "/build/scrutiny-collector-metrics-$(OS)-$(ARCH)" || true
	ldd "/build/scrutiny-collector-metrics-$(OS)-$(ARCH)" || true



# clean:
# 	rm scrutiny-collector-metrics-* scrutiny-web-*
