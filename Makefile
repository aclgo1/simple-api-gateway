
# Variáveis para facilitar a manutenção
PROTO_DIR := ./proto-service/orders/proto
OUT_DIR   := ./proto-service/orders/proto
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

# Lista de repositórios públicos
REPOS = grpc-jwt-login grpc-mail grpc-orders grpc-product grpc-balance
# Repositório Privado
PRIVATE_REPO = concurrency-example

GATEWAY_REPO = simple-api-gateway

BASE_URL = github.com/aclgo1
PARENT_DIR = ..

# Use: make clone-all GITHUB_TOKEN=seu_token_aqui
GITHUB_TOKEN ?= default_token

# Garante que o PATH inclua os binários do Go
export PATH := $(PATH):$(shell go env GOPATH)/bin

.PHONY: all clone-all pull-all up down proto run keys

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
	
	@if [ ! -d "$(PARENT_DIR)/$(PRIVATE_REPO)" ]; then \
		git clone https://oauth2:$(GITHUB_TOKEN)@$(BASE_URL)/$(PRIVATE_REPO).git $(PARENT_DIR)/$(PRIVATE_REPO); \
	else \
		echo "$(PRIVATE_REPO) já existe, pulando clone."; \
	fi

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
	@echo "Atualizando $(GATEWAY_REPO) (current dir)..."
	@git pull
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
keys:
	@echo "Limpando e recriando pasta certs no API Gateway..."
	@rm -rf certs
	@mkdir -p certs
	@echo "Gerando chave privada RSA (.pem)..."
	@openssl genpkey -algorithm RSA -out certs/private_key.pem -pkeyopt rsa_keygen_bits:2048
	@echo "Gerando chave pública RSA (.pem)..."
	@openssl rsa -pubout -in certs/private_key.pem -out certs/public_key.pem
	@echo "Limpando e copiando chave PÚBLICA para os microserviços..."
	@$(foreach repo, $(REPOS) $(PRIVATE_REPO), \
		if [ -d "$(PARENT_DIR)/$(repo)" ]; then \
			echo "Atualizando certs em $(PARENT_DIR)/$(repo)..."; \
			rm -rf $(PARENT_DIR)/$(repo)/certs; \
			mkdir -p $(PARENT_DIR)/$(repo)/certs; \
			cp certs/public_key.pem $(PARENT_DIR)/$(repo)/certs/public_key.pem; \
		else \
			echo "Aviso: Diretório $(repo) não encontrado. Pulando envio."; \
		fi; \
	)
	@echo "Chaves limpas, recriadas e distribuídas com sucesso!"