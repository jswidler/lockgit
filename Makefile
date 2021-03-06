# Go parameters
GOCMD=go
GOBUILD      := $(GOCMD) build
GOCLEAN      := $(GOCMD) clean
GOTEST       := $(GOCMD) test
GOGET        := $(GOCMD) get
BINARY_NAME  := lockgit

VERSION      ?= snapshot
GHRFLAGS     ?=

.PHONY: all test build clean

all: test build
coverage: test cover

build:
ifndef XCBUILD
	$(GOBUILD) -o build/$(BINARY_NAME) -v  -ldflags "-X github.com/jswidler/lockgit/pkg/build.Version=$(VERSION)"
else
	gox -output="build/$(VERSION)/{{.OS}}_{{.Arch}}/$(BINARY_NAME)" -os="darwin linux" -arch="386 amd64" -ldflags "-X github.com/jswidler/lockgit/pkg/build.Version=$(VERSION)"
	mkdir -p build/$(VERSION)
	zip -j build/$(VERSION)/$(BINARY_NAME)_$(VERSION)_darwin_386.zip build/$(VERSION)/darwin_386/$(BINARY_NAME)
	zip -j build/$(VERSION)/$(BINARY_NAME)_$(VERSION)_darwin_amd64.zip build/$(VERSION)/darwin_amd64/$(BINARY_NAME)
	tar -C build/$(VERSION)/linux_386 -czf build/$(VERSION)/$(BINARY_NAME)_$(VERSION)_linux_386.tar.gz $(BINARY_NAME)
	tar -C build/$(VERSION)/linux_amd64 -czf build/$(VERSION)/$(BINARY_NAME)_$(VERSION)_linux_amd64.tar.gz $(BINARY_NAME)
endif


release:
ifdef XCBUILD
	ghr  -u jswidler $(GHRFLAGS) v$(VERSION) build/$(VERSION)
endif

test:
	mkdir -p build
	test -z $$(gofmt -l .)
	bash -c "set -o pipefail && $(GOTEST) -v ./... -coverpkg='github.com/jswidler/lockgit/pkg/...' -coverprofile=build/c.out | tee build/go-test.out"

cover:
	go tool cover -html=build/c.out -o build/coverage.html
	go-junit-report <build/go-test.out > build/go-test-report.xml

clean:
	$(GOCLEAN)
	rm -rf build

run:
	$(GOBUILD) -o build/$(BINARY_NAME) -v
	./build/$(BINARY_NAME)