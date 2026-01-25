# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=main
BINARY_UNIX=$(BINARY_NAME)_unix

# Test parameters
TEST_TIMEOUT=30s
COVERAGE_OUT=coverage.out

.PHONY: all build clean test test-verbose test-coverage test-integration deps help

all: deps test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./main.go

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(COVERAGE_OUT)

# 依存関係のダウンロード
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# 全テストの実行
test:
	$(GOTEST) -timeout $(TEST_TIMEOUT) -v ./...

# 詳細なテスト実行
test-verbose:
	$(GOTEST) -timeout $(TEST_TIMEOUT) -v -race ./...

# カバレッジ付きテスト実行
test-coverage:
	$(GOTEST) -timeout $(TEST_TIMEOUT) -race -coverprofile=$(COVERAGE_OUT) ./...
	$(GOCMD) tool cover -html=$(COVERAGE_OUT)

# 統合テストのみ実行
test-integration:
	$(GOTEST) -timeout $(TEST_TIMEOUT) -v ./tests/integration/...

# リポジトリテストのみ実行
test-repository:
	$(GOTEST) -timeout $(TEST_TIMEOUT) -v ./pkg/repository/...

# ハンドラーテストのみ実行
test-handler:
	$(GOTEST) -timeout $(TEST_TIMEOUT) -v ./pkg/handler/...

# テストのクリーンアップ
test-clean:
	$(GOCLEAN) -testcache

# プロジェクトのセットアップ（初回実行時）
setup: deps
	@echo "プロジェクトのセットアップが完了しました"

# ヘルプ
help:
	@echo "利用可能なコマンド:"
	@echo "  all            - 依存関係ダウンロード、テスト実行、ビルド"
	@echo "  build          - アプリケーションのビルド"
	@echo "  clean          - ビルド成果物とテストキャッシュの削除"
	@echo "  deps           - 依存関係のダウンロードと整理"
	@echo "  test           - 全テストの実行"
	@echo "  test-verbose   - 詳細なテスト実行"
	@echo "  test-coverage  - カバレッジ付きテスト実行"
	@echo "  test-integration - 統合テストのみ実行"
	@echo "  test-repository  - リポジトリテストのみ実行"
	@echo "  test-handler     - ハンドラーテストのみ実行"
	@echo "  test-clean       - テストキャッシュの削除"
	@echo "  setup          - プロジェクトのセットアップ"
	@echo "  help           - このヘルプの表示"
