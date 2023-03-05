.PHONY: build up down logs

# Variables
EXT=development.yml
ENV=development.env

COMPOSE=docker-compose -f docker-compose.yml -f $(EXT) --env-file

# Targets
build:
	$(COMPOSE) $(ENV) build

up:
	$(COMPOSE) $(ENV) up

down:
	$(COMPOSE) $(ENV) down

logs:
	$(COMPOSE) $(ENV) logs -f

# Set the default target to 'up'
.DEFAULT_GOAL := up