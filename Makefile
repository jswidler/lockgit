# Go parameters
GOCMD=go
GOBUILD      := $(GOCMD) build
GOCLEAN      := $(GOCMD) clean
GOTEST       := $(GOCMD) test
GOGET        := $(GOCMD) get
BINARY_NAME  := lockgit

VERSION      ?= snapshot
GHRFLAGS     ?=

all: test build
coverage: test cover

build:
ifndef XCBUILD
	$(GOBUILD) -o build/$(BINARY_NAME) -v
else
	gox -output="build/$(VERSION)/{{.OS}}_{{.Arch}}/$(BINARY_NAME)" -os="darwin linux" -arch="386 amd64"
	mkdir -p pkg/$(VERSION)
	zip -j pkg/$(VERSION)/$(BINARY_NAME)_$(VERSION)_darwin_386.zip build/$(VERSION)/darwin_386/$(BINARY_NAME)
	zip -j pkg/$(VERSION)/$(BINARY_NAME)_$(VERSION)_darwin_amd64.zip build/$(VERSION)/darwin_amd64/$(BINARY_NAME)
	tar -C build/$(VERSION)/linux_386 -czf pkg/$(VERSION)/$(BINARY_NAME)_$(VERSION)_linux_386.tar.gz $(BINARY_NAME)
	tar -C build/$(VERSION)/linux_amd64 -czf pkg/$(VERSION)/$(BINARY_NAME)_$(VERSION)_linux_amd64.tar.gz $(BINARY_NAME)
endif


release:
ifdef XCBUILD
	ghr  -u jswidler $(GHRFLAGS) v$(VERSION) pkg/$(VERSION)
endif

test:
	mkdir -p build
	$(GOTEST) -v ./... -coverpkg='github.com/jswidler/lockgit/src/...' -coverprofile=build/c.out | tee build/go-test.out

cover:
	go tool cover -html=build/c.out -o build/coverage.html
	go-junit-report <build/go-test.out > build/go-test-report.xml

clean:
	$(GOCLEAN)
	rm -rf build
	rm -rf pkg

run:
	$(GOBUILD) -o build/$(BINARY_NAME) -v
	./build/$(BINARY_NAME)

deps:
	$(GOGET) github.com/jstemmer/go-junit-report
	$(GOGET) github.com/mitchellh/gox
	$(GOGET) github.com/mitchellh/go-homedir
	$(GOGET) github.com/olekukonko/tablewriter
	$(GOGET) github.com/pkg/errors
	$(GOGET) github.com/spf13/cobra
	$(GOGET) github.com/spf13/viper