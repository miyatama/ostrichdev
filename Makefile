gocmd=go
gobuild=$(gocmd) build
goclean=$(gocmd) clean
gotest=$(gocmd) test
goget=$(gocmd) get
binary_name=ostrichdev
binary_unix=$(binary_name)_unix

all: test build
build:
    $(gobuild) -o $(binary_name) -v
test:
    $(gotest) -v 
    $(gotest) -v ./ostrich/
clean:
    $(goclean)
    rm -f $(binary_name)
    rm -f $(binary_unix)
run:
    $(gobuild) -o $(binary_name) -v ./...
    ./$(binary_name)
deps:
    $(goget) github.com/hashicorp/logutils
    $(goget) github.com/gin-gonic/gin
