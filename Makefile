BINARY  := yank
PKG     := github.com/korotinm/yank-run-cli
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
LDFLAGS := -s -w -X $(PKG)/cmd.Version=$(VERSION) -X $(PKG)/cmd.Commit=$(COMMIT)

.PHONY: build install vet fmt clean cross

build:
	@mkdir -p bin
	go build -trimpath -ldflags '$(LDFLAGS)' -o bin/$(BINARY) .

install:
	go install -trimpath -ldflags '$(LDFLAGS)' .

vet:
	go vet ./...

fmt:
	gofmt -s -w .

clean:
	rm -rf bin dist

# Cross-compile to all supported platforms.
PLATFORMS := \
	darwin/amd64 darwin/arm64 \
	linux/amd64  linux/arm64 \
	windows/amd64

cross: clean
	@mkdir -p dist
	@for p in $(PLATFORMS); do \
		os=$${p%/*}; arch=$${p#*/}; \
		ext=""; [ "$$os" = "windows" ] && ext=".exe"; \
		out="dist/$(BINARY)-$(VERSION)-$$os-$$arch$$ext"; \
		echo "→ $$out"; \
		GOOS=$$os GOARCH=$$arch go build -trimpath -ldflags '$(LDFLAGS)' -o $$out . || exit 1; \
	done
