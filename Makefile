
# Variáveis para facilitar a manutenção
PROTO_DIR := ./proto-service/orders/proto
OUT_DIR   := ./proto-service/orders/proto
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

# Lista de repositórios públicos
REPOS = grpc-admin grpc-jwt-login grpc-mail grpc-orders grpc-product grpc-balance
# Repositório Privado
PRIVATE_REPO = concurrency-example

BASE_URL = github.com/aclgo1
PARENT_DIR = ..

# Use: make clone-all GITHUB_TOKEN=seu_token_aqui
GITHUB_TOKEN ?= default_token

# Garante que o PATH inclua os binários do Go
export PATH := $(PATH):$(shell go env GOPATH)/bin

.PHONY: all clone-all pull-all up down proto run

all: clone-all up

run:
	go run cmd/simple-api-gateway/main.go

proto: $(PROTO_FILES)
	@echo "Gerando arquivos Go a partir dos protos..."
	protoc --proto_path=$(PROTO_DIR) \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
	@echo "Concluído."

clone-all:
	@echo "Clonando repositórios públicos para $(PARENT_DIR)..."
	@$(foreach repo, $(REPOS), \
		if [ ! -d "$(PARENT_DIR)/$(repo)" ]; then \
			git clone https://$(BASE_URL)/$(repo).git $(PARENT_DIR)/$(repo); \
		else \
			echo "$(repo) já existe, pulando clone."; \
		fi; \
	)
	@echo "Clonando repositório privado $(PRIVATE_REPO)..."
	@if [ ! -d "$(PARENT_DIR)/$(PRIVATE_REPO)" ]; then \
		git clone https://$(GITHUB_TOKEN)@$(BASE_URL)/$(PRIVATE_REPO).git $(PARENT_DIR)/$(PRIVATE_REPO); \
	else \
		echo "$(PRIVATE_REPO) já existe, pulando clone."; \
	fi

pull-all:
	@echo "Atualizando todos os repositórios em $(PARENT_DIR)..."
	@$(foreach repo, $(REPOS) $(PRIVATE_REPO), \
		if [ -d "$(PARENT_DIR)/$(repo)" ]; then \
			echo "Atualizando $(repo)..."; \
			cd $(PARENT_DIR)/$(repo) && git pull; \
		fi; \
	)

up:
	@echo "Iniciando Docker Compose..."
	docker-compose up --build -d

down:
	@echo "Parando serviços..."
	docker-compose down