SHELL := /bin/bash
.DEFAULT_GOAL := help

APP_DIR := VitaLineClinic
DOCKER ?= docker

.PHONY: help build test docker-up docker-down docker-ps docker-logs docker-gui-logs api-health api-info gui-status gui-run

help build test api-health api-info gui-status gui-run:
	@$(MAKE) -C $(APP_DIR) $@

docker-up:
	@$(DOCKER) compose up -d --build

docker-down:
	@$(DOCKER) compose down

docker-ps:
	@$(DOCKER) compose ps

docker-logs:
	@$(DOCKER) compose logs -f api

docker-gui-logs:
	@$(DOCKER) compose logs -f frontend

%:
	@$(MAKE) -C $(APP_DIR) $@
