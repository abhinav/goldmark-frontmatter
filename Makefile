# "go install"-ed binaries will be placed here during development.
export GOBIN ?= $(shell pwd)/bin
export PATH := $(GOBIN):$(PATH)

GO_FILES = $(shell find . \
	   -path '*/.*' -prune -o \
	   '(' -type f -a -name '*.go' ')' -print)

REVIVE = bin/revive
STATICCHECK = bin/staticcheck

TOOLS = $(REVIVE) $(STATICCHECK)

.PHONY: all
all: build lint test

.PHONY: build
build:
	go build ./...

.PHONY: tools
tools: $(TOOLS)

.PHONY: test
test:
	go test -v -race ./...

.PHONY: cover
cover:
	go test -v -race -coverprofile=cover.out -coverpkg=./... ./...
	go tool cover -html=cover.out -o cover.html

.PHONY: lint
lint: gofmt revive staticcheck

.PHONY: gofmt
gofmt:
	$(eval FMT_LOG := $(shell mktemp -t gofmt.XXXXX))
	@gofmt -e -s -l $(GO_FILES) > $(FMT_LOG) || true
	@[ ! -s "$(FMT_LOG)" ] || \
		(echo "gofmt failed. Please reformat the following files:" | \
		cat - $(FMT_LOG) && false)

.PHONY: revive
revive: $(REVIVE)
	revive -set_exit_status ./...

.PHONY: staticcheck
staticcheck: $(STATICCHECK)
	staticcheck ./...

$(REVIVE): tools/go.mod
	cd tools && go install github.com/mgechev/revive

$(STATICCHECK): tools/go.mod
	cd tools && go install honnef.co/go/tools/cmd/staticcheck
