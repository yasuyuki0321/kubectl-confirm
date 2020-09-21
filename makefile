# Go パラメータ
GOCMD = go
GOFMT = goimports
GOBUILD = $(GOCMD) build
GOCLEAN	= $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOGENERATE = ${GOCMD} generate
# ターゲットパラメータ
GOFILES	= $(shell find . -name "*.go")
BINARY_NAME = kubectl-confirm

# タスク
.PHONY: fmt
fmt:
	$(GOFMT) -w ${GOFILES}

.PHONY: setup
setup:
	${GOGET} github.com/jessevdk/go-assets-builder

.PHONY: build
build:
	${GOGENERATE}
	$(GOBUILD) -o $(BINARY_NAME) main.go bindata.go
	rm -f bindata.go

.PHONY: clean
clean:
	$(GOCLEANN)
	rm -f ${BINARY_NAME}

.PHONY: install
install:
	make setup
	make build
	mv $(BINARY_NAME) ${GOPATH}/bin/
