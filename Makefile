# Go parameters
GOCMD=go
GOBUILD      := $(GOCMD) build
GOCLEAN      := $(GOCMD) clean
GOTEST       := $(GOCMD) test
GOGET        := $(GOCMD) get
BINARY_NAME  := lockgit
BINARY_UNIX  := $(BINARY_NAME)_unix
BINARY_DARW  := $(BINARY_NAME)_darwin


all: test build
coverage: test cover

build:
	$(GOBUILD) -o build/$(BINARY_NAME) -v

test:
	mkdir -p build
	$(GOTEST) -v ./... -coverpkg='github.com/jswidler/lockgit/src/...' -coverprofile=build/c.out | tee build/go-test.out

cover:
	go tool cover -html=build/c.out -o build/coverage.html
	go-junit-report <build/go-test.out > build/go-test-report.xml

clean:
	$(GOCLEAN)
	rm -rf build

run:
	$(GOBUILD) -o build/$(BINARY_NAME) -v
	./build/$(BINARY_NAME)

deps:
	$(GOGET) github.com/jstemmer/go-junit-report
	$(GOGET) github.com/mitchellh/go-homedir
	$(GOGET) github.com/olekukonko/tablewriter
	$(GOGET) github.com/pkg/errors
	$(GOGET) github.com/spf13/cobra
	$(GOGET) github.com/spf13/viper


# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o build/$(BINARY_UNIX) -v
build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o build/$(BINARY_DARW) -v
