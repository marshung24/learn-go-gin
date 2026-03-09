# ===============================================
# Makefile - 常用開發指令封裝
# ===============================================

# 載入 .env 檔案（如果存在）
-include .env

# 預設值
DB_HOST ?= localhost
DB_PORT ?= 3306
DB_USER ?= app
DB_PASSWORD ?= app
DB_NAME ?= app

# MySQL DSN（golang-migrate 格式）
DB_DSN = mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)

# ===============================================
# Migration 指令
# ===============================================

# 執行所有 up migration
.PHONY: migrate-up
migrate-up:
	migrate -path migrations -database "$(DB_DSN)" up

# 回滾一版
.PHONY: migrate-down
migrate-down:
	migrate -path migrations -database "$(DB_DSN)" down 1

# 回滾所有 migration
.PHONY: migrate-down-all
migrate-down-all:
	migrate -path migrations -database "$(DB_DSN)" down -all

# 查看當前版本
.PHONY: migrate-version
migrate-version:
	migrate -path migrations -database "$(DB_DSN)" version

# 強制設定版本（用於修復 dirty 狀態）
.PHONY: migrate-force
migrate-force:
	@read -p "Enter version number: " version; \
	migrate -path migrations -database "$(DB_DSN)" force $$version

# 建立新的 migration 檔案
.PHONY: migrate-create
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# ===============================================
# 應用程式指令
# ===============================================

# 啟動應用程式
.PHONY: run
run:
	go run main.go

# 啟動應用程式（使用 air 熱重載）
.PHONY: dev
dev:
	air

# 執行測試
.PHONY: test
test:
	go test ./... -v

# 執行測試並產生覆蓋率報告
.PHONY: test-coverage
test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

# ===============================================
# Docker 指令
# ===============================================

# 啟動所有容器
.PHONY: docker-up
docker-up:
	docker compose up -d

# 停止所有容器
.PHONY: docker-down
docker-down:
	docker compose down

# 查看容器日誌
.PHONY: docker-logs
docker-logs:
	docker compose logs -f

# 重建並啟動容器
.PHONY: docker-rebuild
docker-rebuild:
	docker compose up -d --build

# ===============================================
# 工具指令
# ===============================================

# 整理依賴
.PHONY: tidy
tidy:
	go mod tidy

# 格式化程式碼
.PHONY: fmt
fmt:
	go fmt ./...

# 檢查程式碼
.PHONY: lint
lint:
	golangci-lint run

# 產生 Swagger 文件
.PHONY: swagger
swagger:
	swag init

# ===============================================
# 說明
# ===============================================

.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Migration:"
	@echo "  migrate-up        執行所有 up migration"
	@echo "  migrate-down      回滾一版"
	@echo "  migrate-down-all  回滾所有 migration"
	@echo "  migrate-version   查看當前版本"
	@echo "  migrate-force     強制設定版本"
	@echo "  migrate-create    建立新的 migration"
	@echo ""
	@echo "Application:"
	@echo "  run              啟動應用程式"
	@echo "  dev              使用 air 熱重載啟動"
	@echo "  test             執行測試"
	@echo "  test-coverage    執行測試並產生覆蓋率報告"
	@echo ""
	@echo "Docker:"
	@echo "  docker-up        啟動所有容器"
	@echo "  docker-down      停止所有容器"
	@echo "  docker-logs      查看容器日誌"
	@echo "  docker-rebuild   重建並啟動容器"
	@echo ""
	@echo "Tools:"
	@echo "  tidy             整理依賴"
	@echo "  fmt              格式化程式碼"
	@echo "  lint             檢查程式碼"
	@echo "  swagger          產生 Swagger 文件"
