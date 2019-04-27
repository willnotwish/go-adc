# Suitable for cross compiling for ARM

GOCMD=go
GOBUILD=$(GOCMD) build -v
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get -v

BINARY=go-adc
DIST=dist

# Raspberry pi zero W uses ARM v6
pi0w: arm32v6

pi3B: arm32v7

simulator: init get
	$(GOBUILD) -o $(DIST)/simulator ./cmd/simulator

amd64: init get
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(DIST)/$@/$(BINARY) ./cmd/main

arm32v5: init get
	GOOS=linux GOARCH=arm GOARM=5 $(GOBUILD) -o $(DIST)/$@/$(BINARY) ./cmd/main

arm32v6: init get
	GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD) -o $(DIST)/$@/$(BINARY) ./cmd/main

arm32v7: init get
	GOOS=linux GOARCH=arm GOARM=7 $(GOBUILD) -o $(DIST)/$@/$(BINARY) ./cmd/main

clean:
	rm -fr $(DIST)

init:
	mkdir -p $(DIST)

get:
	$(GOGET) ./...

all: arm32v5 arm32v6 amd64 simulator
