# Suitable for cross compiling for ARM

GOCMD=go
GOBUILD=$(GOCMD) build -v
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get -v

BINARY=adc-cli
DIST=dist

simulator: init get
	$(GOBUILD) -o $(DIST)/simulator ./cmd/simulator

amd64: init get
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(DIST)/$@/$(BINARY) ./cmd/adc-cli

arm32v5: init get
	GOOS=linux GOARCH=arm GOARM=5 $(GOBUILD) -o $(DIST)/$@/$(BINARY) ./cmd/adc-cli

arm32v6: init get
	GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD) -o $(DIST)/$@/$(BINARY) ./cmd/adc-cli

# Raspberry pi zero W uses ARM v6
pi0w: arm32v6

clean:
	rm -fr $(DIST)

init:
	mkdir -p $(DIST)

get:
	$(GOGET) ./...

all: arm32v5 arm32v6 amd64 simulator
