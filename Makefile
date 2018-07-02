# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=webcrawler
BINARY_UNIX=$(BINARY_NAME)_unix
    
all: test build
build: 
	$(GOBUILD) -o $(BINARY_NAME) -v main.go
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v main.go
	./$(BINARY_NAME)
deps:
	$(GOGET) golang.org/x/net/html
	$(GOGET) github.com/spf13/cobra
        
# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v main.go
docker-build:
	docker build . -t webcrawler 