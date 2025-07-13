# Go タスク管理 API

Go タスク管理 API は、基本的なタスク管理機能を提供する RESTful API です。Go 言語の学習プロジェクトとして作成され、ワークフロー機能を含む高度なタスク管理機能を提供します。

> **⚠️ 注意**: このプロジェクトは学習・練習目的で作成されています。本番環境用の設定は理論的な実装であり、実際の本番環境での検証は行っていません。

## 更新履歴

### 2025/07/13 追加分
- **Docker化対応**: 開発・本番環境に対応したコンテナ化
- **マルチステージビルド**: セキュリティとパフォーマンスを最適化
- **環境変数対応**: 12 Factor App原則に基づく設定管理
- **本番環境対応**: Docker Secrets、Nginx、監視機能（練習・学習用実装）

### 2024/01/13 追加分
- **ワークフロー機能の追加**: ワークフロー機能を導入しました。この機能はタスク管理プロセスを柔軟に制御し、複数のステップを含むタスクの進行を追跡するのに役立ちます。
- **Postman ワークスペースの追加**: API のテストとドキュメンテーションのための Postman ワークスペースを追加しました。

# 使用技術

## バックエンド
- **Go 1.20**: メインプログラミング言語
- **Gin**: Webフレームワーク
- **GORM ORM**: O/Rマッパー
- **MySQL 8.0**: データベース

## インフラ・運用
- **Docker**: コンテナ化
- **Docker Compose**: オーケストレーション
- **Nginx**: リバースプロキシ（本番環境）
- **Alpine Linux**: 軽量コンテナベース

## 開発ツール
- **Air**: ホットリロード（開発環境）
- **Postman**: API テスト

# 機能

- 新しいタスクの作成
- すべてのタスクの一覧表示
- 特定のタスクの取得
- 特定のタスクの更新
- 特定のタスクの削除
- ワークフローの定義と管理
- ワークフローにおけるタスクのステータスの追跡

# エンドポイント

## タスク関係

- POST /task：新しいタスクを作成します。
- GET /tasks：すべてのタスクを一覧表示します。
- GET /task/{id}：特定のタスクを取得します。
- PUT /task/{id}：特定のタスクを更新します。
- DELETE /task/{id}：特定のタスクを削除します。

## ワークフロー関係

- POST /workflows：新しいワークフローを作成します。
- GET /workflows/{id}：特定のワークフローを取得します。
- PUT /task/advance/{task_id}：特定のタスクのワークフローステータスを次のステップに進めます。
- POST /workflow/{workflow_id}/start/{task_id}：特定のワークフローを特定のタスクに関連付けて開始します。
- POST /workflow-steps：ワークフローのステップを作成します。

# クイックスタート

## 開発環境での起動

### 1. Dockerを使用した起動（推奨）
```bash
# シンプル版
docker-compose up -d

# または詳細版
docker-compose -f docker-compose.dev.yml up -d
```

### 2. ローカル環境での起動
```bash
# 依存関係のインストール
go mod download

# サーバー起動
go run main.go
```

### 3. ホットリロード（開発用）
```bash
# Airを使用したホットリロード
air
```

### 4. VSCodeデバッガー使用

#### ローカルデバッグ
```bash
# MySQLのみDockerで起動
docker-compose up -d mysql

# VSCodeでF5キー押下、または
# Debug > Start Debugging > "Launch Go App (Docker DB)"
```

#### Docker リモートデバッグ
```bash
# 自動起動でリモートデバッグ（推奨）
# VSCodeで Debug > Start Debugging > "Remote Debug Docker Container"

# または手動でコンテナ起動
docker-compose -f docker-compose.debug.yml up -d --build
# その後、Debug > Start Debugging > "Remote Debug (Manual Container)"
```

#### 🚀 ホットリロード + リモートデバッグ（最強環境）
```bash
# Air + Delve 統合環境（GoLand風）
# VSCodeで Debug > Start Debugging > "Remote Debug + Hot Reload"

# または手動起動
docker-compose -f docker-compose.hotreload.yml up -d --build
```

**特徴:**
- ファイル保存で自動リビルド・再起動
- ブレークポイントが保持される
- GoLandのような快適な開発体験

詳細な操作方法は以下を参照:
- [debug-guide.md](./debug-guide.md) - 基本的なリモートデバッグ
- [VSCode-Docker-Air-Delve-Guide.md](./VSCode-Docker-Air-Delve-Guide.md) - ホットリロード統合環境

## 本番環境での起動（学習・練習用）

> **⚠️ 重要**: 以下の本番環境設定は学習目的で作成されており、実際の本番環境での使用前には十分な検証とセキュリティ監査が必要です。

### 1. 環境設定
```bash
# 環境変数ファイルをコピー
cp .env.prod.example .env.prod
# .env.prodを編集

# Docker Secretsを作成
echo "secure-password" | docker secret create task_management_db_password -
echo "root-password" | docker secret create task_management_mysql_root_password -
echo "user-password" | docker secret create task_management_mysql_password -
```

### 2. デプロイ実行
```bash
# 自動デプロイスクリプト実行
./deploy-prod.sh

# または手動実行
docker-compose -f docker-compose.prod.yml up -d
```

## アクセス情報

- **API**: http://localhost:8080
- **MySQL**: localhost:3306 (開発) / localhost:3307 (本番)
- **ヘルスチェック**: http://localhost:8080/health

# ファイル構成

```
├── main.go                    # メインアプリケーション
├── pkg/                       # パッケージディレクトリ
│   ├── config/
│   │   ├── dbconfig.go       # データベース設定
│   │   └── db/               # SQLスキーマとデータ
├── Dockerfile                # 開発用Dockerfile
├── Dockerfile.prod           # 本番用Dockerfile
├── docker-compose.yml        # 開発用Compose設定
├── docker-compose.dev.yml    # 開発用詳細Compose設定
├── docker-compose.prod.yml   # 本番用Compose設定
├── .env.prod.example         # 本番環境変数テンプレート
├── deploy-prod.sh           # 本番デプロイスクリプト
├── nginx/
│   └── nginx.conf           # Nginx設定（本番用）
└── DEVELOPMENT_SUMMARY.md   # 開発履歴とDocker詳細解説
```

# Docker環境について

このプロジェクトは開発環境と本番環境に対応したDocker設定を提供します：

## 開発環境の特徴
- シンプルな設定
- 平文パスワード使用
- ホットリロード対応
- デバッグモード有効

## 本番環境の特徴（練習・学習用実装）
- セキュリティ強化（非rootユーザー実行）
- Docker Secretsによる機密情報管理
- Nginxリバースプロキシ
- ヘルスチェック・監視機能
- リソース制限とリスタートポリシー

> **注意**: これらの設定は学習目的であり、実際の本番環境では追加のセキュリティ対策が必要です。

詳細な解説は [DEVELOPMENT_SUMMARY.md](./DEVELOPMENT_SUMMARY.md) を参照してください。

# 注意事項

## 学習プロジェクトについて
- このプロジェクトはGo言語とDocker技術の学習目的で作成されています
- 本番環境用設定は理論的な実装であり、実運用での検証は未実施です
- 実際の本番環境での使用は推奨しません

## 技術的な制限
- テストコードは未実装です
- 本番環境用設定は学習・練習レベルの実装です
- 実際の本番環境ではSSL/TLS設定を適切に行ってください
- Docker Secretsの作成が本番環境では必須です
- セキュリティ監査やパフォーマンステストは未実施です

## 実用化に向けて必要な作業
- 包括的なテストコードの実装
- セキュリティ監査とペネトレーションテスト
- パフォーマンステストと負荷テスト
- ログ監視とアラート設定
- バックアップ・リカバリ戦略の策定
- CI/CDパイプラインの構築
