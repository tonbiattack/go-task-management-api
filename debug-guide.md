# Docker リモートデバッグガイド

このガイドでは、Docker上で動作するGoアプリケーションをVSCodeでリモートデバッグする方法を説明します。

## 前提条件

- VSCodeにGo拡張機能がインストール済み
- Dockerがインストール済み
- Docker Composeが利用可能

## リモートデバッグの仕組み

```
VSCode (ホスト)  ←→ ポート2345 ←→ Docker Container (Delve)
                                        ↓
                                   Go Application
```

1. DockerコンテナでDelveデバッガーサーバーが起動
2. VSCodeがポート2345経由でDelveに接続
3. ブレークポイントやステップ実行が可能

## 使用方法

### 方法1: 自動起動でリモートデバッグ（推奨）

```bash
# 1. VSCodeでDebug > Start Debugging > "Remote Debug Docker Container"
# または F5キーを押してデバッグ設定を選択
```

この方法では：
- デバッグコンテナが自動で起動
- Delveデバッガーサーバーが開始
- VSCodeが自動でアタッチ

### 方法2: 手動でコンテナ起動

```bash
# 1. デバッグ用コンテナを手動で起動
docker-compose -f docker-compose.debug.yml up -d --build

# 2. コンテナが起動したらVSCodeでアタッチ
# Debug > Start Debugging > "Remote Debug (Manual Container)"
```

### 方法3: コマンドライン操作

```bash
# コンテナ起動
docker-compose -f docker-compose.debug.yml up -d --build

# デバッグ接続確認
docker logs task-management-api-debug

# コンテナ停止
docker-compose -f docker-compose.debug.yml down
```

## デバッグ操作

### ブレークポイントの設定
1. コードの行番号左側をクリック
2. 赤い丸（ブレークポイント）が表示される
3. F5でデバッグ開始

### デバッグ操作キー
- **F5**: デバッグ開始/続行
- **F10**: ステップオーバー（次の行）
- **F11**: ステップイン（関数内部）
- **Shift+F11**: ステップアウト（関数から抜ける）
- **Shift+F5**: デバッグ停止

### 変数の確認
- 左側のパネルの「変数」タブで確認
- コード上で変数にマウスオーバーで値表示
- ウォッチ式で特定の変数を監視

## トラブルシューティング

### 1. デバッガーに接続できない

```bash
# コンテナの状態確認
docker ps

# デバッグポートの確認
docker port task-management-api-debug

# ログの確認
docker logs task-management-api-debug
```

### 2. ブレークポイントが効かない

```bash
# コンテナを完全に再ビルド
docker-compose -f docker-compose.debug.yml down
docker-compose -f docker-compose.debug.yml up -d --build --force-recreate
```

### 3. ソースコードが同期されない

- ファイルパスの確認
- `substitutePath`設定の確認
- コンテナ内のソースコードと一致しているか確認

### 4. Delveデバッガーサーバーが起動しない

```bash
# コンテナ内でDelveの状態確認
docker exec -it task-management-api-debug ps aux

# 手動でDelve起動確認
docker exec -it task-management-api-debug dlv version
```

## 高度な設定

### カスタム起動引数
`docker-compose.debug.yml`でDelveの起動引数をカスタマイズ可能：

```yaml
CMD ["dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "--continue", "exec", "./main", "--", "--custom-arg"]
```

### 別ポートでデバッグ
ポート競合がある場合：

```yaml
ports:
  - "2346:2345"  # ホスト側ポートを変更
```

対応してVSCodeの`launch.json`も修正：

```json
"port": 2346
```

## 利点と制限

### 利点
- 本番環境に近い状態でデバッグ可能
- Dockerコンテナの環境変数や依存関係を利用
- 複数サービス間の連携もデバッグ可能

### 制限
- ローカルデバッグより若干動作が重い
- ネットワーク経由のため遅延がある
- コンテナの再起動が必要な場合がある

## 参考資料

- [Delve公式ドキュメント](https://github.com/go-delve/delve)
- [VSCode Go拡張機能](https://code.visualstudio.com/docs/languages/go)
- [Docker Compose リファレンス](https://docs.docker.com/compose/)