SHELL := /bin/bash
.DEFAULT_GOAL := help

APP_DIR := VitaLineClinic
DOCKER ?= docker

.PHONY: help build test docker-up docker-down docker-ps docker-logs api-health api-info

help build test api-health api-info:
	@$(MAKE) -C $(APP_DIR) $@

docker-up:
	@$(DOCKER) compose up -d --build

docker-down:
	@$(DOCKER) compose down

docker-ps:
	@$(DOCKER) compose ps

docker-logs:
	@$(DOCKER) compose logs -f api

%:
	@$(MAKE) -C $(APP_DIR) $@
