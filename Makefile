PROJECT := sourcemap

BOLD := \033[1m
RESET := \033[0m

.DEFAULT_GOAL := help

.PHONY: refresh # Update the dependencies
refresh:
	@echo "$(BOLD)Refreshing dependencies...$(RESET)"
	@go mod tidy

.PHONY: build # Build the project
build:
	@echo "${BOLD}Building ${PROJECT}...${RESET}"
	@go build -o bin/${PROJECT} -ldflags '-s -w'

.PHONY: help # Display the help message
help:
	@echo "${BOLD}Available targets:${RESET}"
	@cat Makefile | grep '.PHONY: [a-z]' | sed 's/.PHONY: / /g' | sed 's/ #* / - /g'
