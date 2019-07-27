# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=acons
BINARY_LINUX_AMD64=$(BINARY_NAME)_linux_amd64
BINARY_LINUX_ARM=$(BINARY_NAME)_linux_arm
BINARY_WIN=$(BINARY_NAME)_win
BINARY_MAC=$(BINARY_NAME)_mac
ALL_BINARIES=$(BINARY_NAME) $(BINARY_MAC) $(BINARY_WIN) $(BINARY_LINUX_AMD64) $(BINARY_LINUX_ARM)
    
all: test build
build: *.go
	$(GOBUILD) -o $(BINARY_NAME) -v

cross-compile: build-mac build-win build-linux_arm build-linux_amd64

test: *.go
	$(GOTEST) -v ./...

clean: 
	$(GOCLEAN)
	rm -f $(ALL_BINARIES)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

deps:
	$(GOGET) github.com/markbates/goth
	$(GOGET) github.com/markbates/pop
    
# Cross compilation
build-mac: *.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_MAC) -v

build-win: *.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WIN) -v

build-linux_arm: *.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm $(GOBUILD) -o $(BINARY_LINUX_ARM) -v

build-linux_amd64: *.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_LINUX_AMD64) -v
