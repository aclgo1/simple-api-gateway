
# Variáveis para facilitar a manutenção
PROTO_DIR := ./proto-service/orders/proto
OUT_DIR   := ./proto-service/orders/proto
# Busca automaticamente todos os arquivos .proto na pasta
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

REPOS = grpc-admin grpc-jwt-login grpc-mail grpc-orders grpc-product grpc-balance
BASE_URL = https://github.com/aclgo1
PARENT_DIR = ..

# Garante que o PATH inclua os binários do Go
export PATH := $(PATH):$(shell go env GOPATH)/bin

.PHONY: all clone-all pull-all up down proto api run

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
	@echo "Clonando repositórios para $(PARENT_DIR)..."
	@$(foreach repo, $(REPOS), \
		if [ ! -d "$(PARENT_DIR)/$(repo)" ]; then \
			git clone $(BASE_URL)/$(repo).git $(PARENT_DIR)/$(repo); \
		else \
			echo "$(repo) já existe, pulando clone."; \
		fi; \
	)

pull-all:
	@echo "Atualizando repositórios em $(PARENT_DIR)..."
	@$(foreach repo, $(REPOS), \
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