.PHONY: test build

ifeq ("$(wildcard .env)","")
    $(shell cp env.sample .env)
endif

include .env

$(eval export $(grep -v '^#' .env | xargs -0))

GO_MODULE := loan-service
LDFLAGS   := -X "$(GO_MODULE)/config.Version=$(VERSION)"
VERSION   ?= $(shell git rev-parse --short HEAD)

################
# DEPENDENCIES #
################

# mockery
# https://github.com/vektra/mockery

MOCKERY := $(shell command -v mockery || echo "bin/mockery")
mockery: bin/mockery ## Installs mockery (mocks generation)

bin/mockery: VERSION := 2.32.4
bin/mockery: GITHUB  := vektra/mockery
bin/mockery: ARCHIVE := mockery_$(VERSION)_$(OSTYPE)_x86_64.tar.gz
bin/mockery: bin
	@printf "Installing mockery... "
	@curl -Ls $(call github_url) | tar -zOxf -  mockery > $@ && chmod +x $@
	@echo "done."

# golangci-lint
# https://github.com/golangci/golangci-lint

GOLANGCI := $(shell command -v golangci-lint || echo "bin/golangci-lint")
golangci-lint: bin/golangci-lint ## Installs golangci-lint (linter)

bin/golangci-lint: VERSION := 1.56.2
bin/golangci-lint: GITHUB  := golangci/golangci-lint
bin/golangci-lint: ARCHIVE := golangci-lint-$(VERSION)-$(OSTYPE)-amd64.tar.gz
bin/golangci-lint: bin
	@printf "Installing golangci-lint... "
	@curl -Ls $(shell echo $(call github_url) | tr A-Z a-z) | tar -zOxf - $(shell printf golangci-lint-$(VERSION)-$(OSTYPE)-amd64/golangci-lint | tr A-Z a-z ) > $@ && chmod +x $@
	@echo "done."

###########
# SCRIPTS #
###########

setup: $(MOCKERY) $(GOLANGCI) $(AIR)
setup:
	@docker build -f Dockerfile.database \
		--build-arg DB_USER=$(DB_USER) \
		--build-arg DB_PASSWORD=$(DB_PASSWORD) \
		--build-arg DB_NAME=$(DB_NAME) \
		-t loan-service-db-setup .
	@docker run --env-file .env --rm -it -v=/tmp/postgres/data:/var/lib/postgresql/data loan-service-db-setup
	@echo "Required tools are installed"

dep:
	@go mod tidy
	@go mod download

seed-db:
	@DB_HOST=$(LOCAL_DB_HOST) DB_PORT=$(LOCAL_DB_PORT) DB_NAME=$(DB_NAME) \
		DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) \
		go run database/seed/seed.go

build:
	@go mod tidy
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ./build/web ./app/web/

run: dep build
run:
	@docker-compose up --build

test:
	@go test ./... --short -cover

lint: $(GOLANGCI)
	@golangci-lint run

gen-mocks: $(MOCKERY)
	@mockery --all

