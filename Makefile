PROJECT := sourcemap

BOLD := \033[1m
RESET := \033[0m

.DEFAULT_GOAL := help

.PHONY: refresh # Update the dependencies
refresh:
	@echo "$(BOLD)Refreshing dependencies...$(RESET)"
	@go mod tidy

.PHONY: build/ui # Build the UI
build/ui:
	@echo "${BOLD}Building UI...${RESET}"
	@cd frontend && pnpm run build

.PHONY: build/core # Build the project
build/core: refresh build/ui
	@echo "${BOLD}Building ${PROJECT}...${RESET}"
	@go build -o bin/${PROJECT} -ldflags '-s -w'

.PHONY: help # Display the help message
help:
	@echo "${BOLD}Available targets:${RESET}"
	@cat Makefile | grep '.PHONY: [a-z]' | sed 's/.PHONY: / /g' | sed 's/ #* / - /g'
