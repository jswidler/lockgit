# Go parameters
GOCMD=go
GOBUILD      := $(GOCMD) build
GOCLEAN      := $(GOCMD) clean
GOTEST       := $(GOCMD) test
GOGET        := $(GOCMD) get
BINARY_NAME  := lockgit
BINARY_UNIX  := $(BINARY_NAME)_unix

TEST_RESULTS ?= build

all: test build
build:
	$(GOBUILD) -o $(TEST_RESULTS)/$(BINARY_NAME) -v

test:
	mkdir -p $(TEST_RESULTS)
	$(GOTEST) -v ./... -coverpkg='github.com/jswidler/lockgit/src/...' -coverprofile=$(TEST_RESULTS)/c.out
	go tool cover -html=$(TEST_RESULTS)/c.out -o $(TEST_RESULTS)/coverage.html

clean:
	$(GOCLEAN)
	rm -rf $(TEST_RESULTS)

run:
	$(GOBUILD) -o $(TEST_RESULTS)/$(BINARY_NAME) -v
	./$(TEST_RESULTS)/$(BINARY_NAME)

deps:
	# $(GOGET) github.com/markbates/goth
	# $(GOGET) github.com/markbates/pop


# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(TEST_RESULTS)/$(BINARY_UNIX) -v
docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v
