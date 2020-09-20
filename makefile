# Go パラメータ
GOCMD?=go
GOFMT:=goimports
GOBUILD:=$(GOCMD) build
GOCLEAN:=$(GOCMD) clean
GOTEST:=$(GOCMD) test
GOGET:=$(GOCMD) get
GOFILES:=$(shell find . -name "*.go")
BINARY_NAME:=kubectl-confirm

.PHONY: fmt
fmt:
	$(GOFMT) -w ${GOFILES}

.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v

.PHONY: clean
clean:
	$(GOCLEANN)
	rm -f ${BINARY_NAME}

.PHONY: install
install:
	make build
	mv ${BINARY_NAME} /usr/local/bin/${BINARY_NAME}
