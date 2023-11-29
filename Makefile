##
# Makefile to help manage docker-compose services
# Built on list_targets-Makefile:
#
# Inspired from:
#     https://gist.github.com/iNamik/73fd1081fe299e3bc897d613179e4aee
#
.PHONY: help about args list targets services up down ps client-api admin-api gateway-api

# If you need sudo to execute docker, then update these aliases
#
DOCKER         := docker
DOCKER_COMPOSE := docker compose
MOCK_GEN_BIN   := mockgen

# Default docker-compose file
#
FILE_DEFAULT_NAME := docker-compose.yml

# Default container for docker actions
# NOTE: EDIT THIS TO AVOID WARNING/ERROR MESSAGES
#
DB_CONTAINER_DEFAULT := "maani-db"
STORE_CONTAINER_DEFAULT := "store"
RETREIVAL_CONTAINER_DEFAULT := "retreival"
SQL_MIGRATION_PATH := "pkg/database/ent/migrate/data_seed.sql"

# Default db config for docker actions
# NOTE: EDIT THIS TO AVOID WARNING/ERROR MESSAGES
#
DB_USERNAME := "postgres"
DB_NAME := "postgres"

# Shell command for 'shell' target
#
SHELL_CMD := /bin/sh
UNAME := $(shell uname)


ME  := $(realpath $(firstword $(MAKEFILE_LIST)))
PWD := $(dir $(ME))

file       ?= "$(PWD)/$(FILE_DEFAULT_NAME)"
service    ?=
services   ?= $(service)
sql_migration_path ?= "$(PWD)/$(SQL_MIGRATION_PATH)"


.DEFAULT_GOAL := help


##
# help
# Displays a (hopefully) useful help screen to the user
# NOTE: Keep 'help' as first target in case .DEFAULT_GOAL is not honored
#
help: about targets args ## This help screen
ifeq ($(DB_CONTAINER_DEFAULT),"")
	$(warning WARNING: DB_CONTAINER_DEFAULT is not set. Please edit makefile)
else ifeq ($(DB_USERNAME),"")
	$(warning WARNING: DB_USERNAME is not set. Please edit makefile)
else ifeq ($(DB_NAME),"")
	$(warning WARNING: DB_NAME is not set. Please edit makefile)
else ifeq ($(SQL_MIGRATION_PATH),"")
	$(warning WARNING: SQL_MIGRATION_PATH is not set. Please edit makefile)
else ifeq ($(STORE_CONTAINER_DEFAULT),"")
	$(warning WARNING: STORE_CONTAINER_DEFAULT is not set. Please edit makefile)
else ifeq ($(RETREIVAL_CONTAINER_DEFAULT),"")
	$(warning WARNING: RETREIVAL_CONTAINER_DEFAULT is not set. Please edit makefile)
endif

about:
	@echo
	@echo "Makefile to help manage Maani Server's"

args:
	@echo
	@echo "Target arguments:"
	@echo
	@echo "    " "file"      "\t" "Location of docker-compose file (default = './$(FILE_DEFAULT_NAME)')"
	@echo "    " "service"   "\t" "Target service for docker-compose actions (default = all services)"
	@echo "    " "services"  "\t" "Target services for docker-compose actions (default = all services)"


##
# list
# Displays a list of targets, using '##' comment as target description
#
# NOTE: ONLY targets with ## comments are shown
#
list: targets ## see 'targets'
targets:  ## Lists targets
	@echo
	@echo "Make targets:"
	@echo
	@cat $(ME) | \
	sed -n -E 's/^([^.][^: ]+)\s*:(([^=#]*##\s*(.*[^[:space:]])\s*)|[^=].*)$$/    \1	\4/p' | \
	sort -u | \
	expand -t15
	@echo

##
# services
#
services: ## Lists services
	@$(DOCKER_COMPOSE) -f "$(file)" ps --services

##
# up
#
up: ## Starts containers (in detached mode) [file|service|services]
	@$(DOCKER_COMPOSE) -f "$(file)" up --detach $(services)

##
# stop
#
stop: ## stop containers [file|service|services]
	@$(DOCKER_COMPOSE) -f "$(file)" stop


##
# start
#
start: ## start containers [file|service|services]
	@$(DOCKER_COMPOSE) -f "$(file)" start

##
# down
#
down: ## Removes containers (preserves images and volumes) [file]
	@$(DOCKER_COMPOSE) -f "$(file)" down

##
# ps
#
ps: ## Shows status of containers [file|service|services]
	@$(DOCKER_COMPOSE) -f "$(file)" ps $(services)

##
# mock
#
mock: ## generates mock models from interface
	@$(MOCK_GEN_BIN) -source=./pkg/database/database.go -destination=./pkg/database/mocks/database_mock.go

coverage: ## test coverage
	go test ./... -v -coverprofile .coverage.out  -gcflags=-l
ifdef text
	go tool cover -func .coverage.out
else
	go tool cover -html=.coverage.out
endif

test: ## run tests
	go test ./... -gcflags=-l -v

##
# add-schema
#
add-schema: ## Generates Ent ORM schema for project [name]
	go run -mod=mod entgo.io/ent/cmd/ent new --target "./pkg/database/ent/schema" "$(name)"

##
# generate-schema
#
generate-schema:
	go run -mod=mod entgo.io/ent/cmd/ent generate "./pkg/database/ent/schema" --feature sql/upsert --feature sql/lock

##
# generate-gateway-api
#
generate-gateway-api: ## Generate Swagger Documentation from our models for Maani Gateway-Side APIs
	@swagger generate spec -w "./docs/swagger/gateway/" -o "./docs/swagger/gateway.yaml" --scan-models

##
# gateway-api
#
gateway-api:  ## Generates Swagger Documentation from our models for Maani Gateway-Side APIs
	@swagger serve -F=swagger "./docs/swagger/gateway.yaml"


##
# go-build
#
go-build:
	go build -o bin/store cmd/store/main.go
	go build -o bin/retreival cmd/retreival/main.go

##
# godoc
#
godoc: ## Generates Documentation from comments for Maani packages
	$(go env GOPATH)/bin/godoc -http=localhost:6060

##
# seed-data
#
seed-data:
	$(DOCKER) cp $(SQL_MIGRATION_PATH) $(DB_CONTAINER_DEFAULT):/migration.sql
	$(DOCKER) exec -i $(DB_CONTAINER_DEFAULT) psql -U $(DB_USERNAME) -d $(DB_NAME) -f /migration.sql

##
# run
#
run:
	docker build -f Dockerfile -t maani:latest .
	@$(DOCKER_COMPOSE) -f "$(file)" up -d