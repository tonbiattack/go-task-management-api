#!/bin/bash

# 本番環境デプロイスクリプト
set -e

echo "=== Task Management API 本番デプロイ開始 ==="

# 設定ファイル確認
if [ ! -f ".env.prod" ]; then
    echo "エラー: .env.prod ファイルが見つかりません"
    echo ".env.prod.example をコピーして設定してください"
    exit 1
fi

# 環境変数読み込み
source .env.prod

# Docker Secretsの作成確認
echo "Docker Secrets確認中..."
REQUIRED_SECRETS=("task_management_db_password" "task_management_mysql_root_password" "task_management_mysql_password")

for secret in "${REQUIRED_SECRETS[@]}"; do
    if ! docker secret ls | grep -q "$secret"; then
        echo "エラー: Docker Secret '$secret' が存在しません"
        echo "以下のコマンドで作成してください:"
        echo "echo 'your-password' | docker secret create $secret -"
        exit 1
    fi
done

# データディレクトリ作成
echo "データディレクトリ作成中..."
sudo mkdir -p ${DATA_PATH:-./data}/mysql
sudo chown -R 999:999 ${DATA_PATH:-./data}/mysql

# バックアップディレクトリ作成
mkdir -p ./backup

# 既存のコンテナ停止
echo "既存のコンテナ停止中..."
docker-compose -f docker-compose.prod.yml down --remove-orphans

# イメージビルド
echo "本番用イメージビルド中..."
docker-compose -f docker-compose.prod.yml build --no-cache app

# データベース起動と初期化待機
echo "データベース起動中..."
docker-compose -f docker-compose.prod.yml up -d mysql
echo "データベース初期化待機中（60秒）..."
sleep 60

# アプリケーション起動
echo "アプリケーション起動中..."
docker-compose -f docker-compose.prod.yml up -d app

# Nginx起動（設定ファイルが存在する場合）
if [ -f "./nginx/nginx.conf" ]; then
    echo "Nginx起動中..."
    docker-compose -f docker-compose.prod.yml up -d nginx
fi

# ヘルスチェック
echo "ヘルスチェック実行中..."
sleep 30

# アプリケーションの起動確認
if curl -f -s "http://localhost:${APP_PORT:-8080}/health" > /dev/null 2>&1; then
    echo "✅ アプリケーションが正常に起動しました"
else
    echo "❌ アプリケーションの起動に失敗しました"
    echo "ログを確認してください:"
    echo "docker-compose -f docker-compose.prod.yml logs app"
    exit 1
fi

# サービス状況表示
echo ""
echo "=== デプロイ完了 ==="
echo "サービス状況:"
docker-compose -f docker-compose.prod.yml ps

echo ""
echo "アクセス情報:"
echo "- API: http://${DOMAIN:-localhost}:${APP_PORT:-8080}"
echo "- MySQL: localhost:${MYSQL_PORT:-3307}"

if [ -f "./nginx/nginx.conf" ]; then
    echo "- Web: http://${DOMAIN:-localhost}:${NGINX_HTTP_PORT:-80}"
    echo "- Web (HTTPS): https://${DOMAIN:-localhost}:${NGINX_HTTPS_PORT:-443}"
fi

echo ""
echo "監視コマンド:"
echo "- ログ確認: docker-compose -f docker-compose.prod.yml logs -f"
echo "- 状況確認: docker-compose -f docker-compose.prod.yml ps"
echo "- 停止: docker-compose -f docker-compose.prod.yml down"

echo ""
echo "=== デプロイ完了 ==="