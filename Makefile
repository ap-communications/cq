.PHONY: help
.DEFAULT_GOAL := help

OS_LIST=windows linux darwin
ARCH_LIST=amd64 386

test: ## Run all test
	@go test ./... -cover

build: ## Build binary file
	@go build -o cq src/cq.go

cross-compile: ## Build binaries for Windows, Linux, and macOS of x64 and x86
	@mkdir bin
	@for GOOS in ${OS_LIST}; do \
		for GOARCH in ${ARCH_LIST}; do \
			if [ $$GOOS = windows ]; then \
				GOOS=$$GOOS GOARCH=$$GOARCH go build -o bin/cq-$$GOOS-$$GOARCH.exe src/cq.go; \
			else \
				GOOS=$$GOOS GOARCH=$$GOARCH go build -o bin/cq-$$GOOS-$$GOARCH src/cq.go; \
			fi \
		done \
	done

clean:
	@-rm -rf bin/

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
