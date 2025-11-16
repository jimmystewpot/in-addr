TOOL := "jimmystewpot/in-addr"
INTERACTIVE := $(shell [ -t 0 ] && echo 1)
TEST_DIRS := ./...
REPORTS_DIR := ci
GOLANGCI_LINT_IMAGE := golangci/golangci-lint
GOLANGCI_LINT_VERSION := v2.6.2
GOLANGCI_LINT_CMD := docker run --rm -v ${PWD}:/app -w /app ${GOLANGCI_LINT_IMAGE}:${GOLANGCI_LINT_VERSION} golangci-lint

test-all: deps lint test

reportdir:
	if [ ! -d "$(REPORTS_DIR)" ]; then mkdir $(REPORTS_DIR);  fi

deps:
	@echo ""
	@echo "***** Installing dependencies for ${TOOL} *****"
	go clean --cache
	docker pull ${GOLANGCI_LINT_IMAGE}:${GOLANGCI_LINT_VERSION}

lint: deps reportdir
	@echo ""
	@echo "***** linting ${TOOL} with golangci-lint *****"
ifdef INTERACTIVE
	${GOLANGCI_LINT_CMD} run -v $(TEST_DIRS)
else
	${GOLANGCI_LINT_CMD} run --output.checkstyle.path stdout -v $(TEST_DIRS) 1> $(REPORTS_DIR)/checkstyle-lint.xml
endif
.PHONY: lint

test:
	@echo ""
	@echo "***** Testing ${TOOL} *****"
ifdef INTERACTIVE
	go test -a -v -race $(TEST_DIRS)
else
	go test -a -v -race -coverprofile=$(REPORTS_DIR)/coverage.txt -covermode=atomic -json $(TEST_DIRS) 1> $(REPORTS_DIR)/testreport.json
endif
	@echo ""