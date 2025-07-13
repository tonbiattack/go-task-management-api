# VSCode + Docker + Air + Delve で快適な Golang Debug ライフ

このガイドでは、参考記事「[GoLand + Docker + air + delve で快適な Golang Debug ライフ](https://zenn.dev/smpeotn/articles/8ee4b6b14970ff)」を参考に、VSCode環境で同様の快適なデバッグ環境を構築する方法を説明します。

## 🎯 目標

- **Air**: ファイル変更時の自動リロード
- **Delve**: リモートデバッグ
- **Docker**: 環境の統一化
- **VSCode**: 統合開発環境

これらを組み合わせて、**ファイル保存 → 自動ビルド → 自動再起動 → デバッグ継続**の流れを実現します。

## 🛠️ 技術スタック

```
VSCode ←→ Docker Container (Air + Delve + Go App)
             ↓
           MySQL Container
```

- **Air**: ホットリロードツール
- **Delve**: Goデバッガー
- **Docker**: コンテナ環境
- **VSCode**: エディタ + デバッガークライアント

## 📁 ファイル構成

```
├── .air.toml                      # Air設定ファイル
├── Dockerfile.hotreload           # ホットリロード用Dockerfile
├── docker-compose.hotreload.yml   # 統合環境設定
├── .vscode/
│   ├── launch.json               # デバッグ設定
│   └── tasks.json                # VSCodeタスク
└── tmp/                          # Air一時ファイル
```

## 🚀 セットアップ手順

### 1. 環境起動

```bash
# ホットリロード + デバッグ環境起動
docker-compose -f docker-compose.hotreload.yml up -d --build
```

### 2. 起動確認

```bash
# コンテナ状況確認
docker ps --filter name=task-management

# ログ確認（Airが動作しているか）
docker logs task-management-hotreload
```

正常な場合、以下のようなログが表示されます：
```
  __    _   ___  
 / /\  | | | |_) 
/_/--\ |_| |_| \_ , built with Go 

watching .
building...
API server listening at: [::]:2345
```

### 3. VSCodeでデバッグ接続

```bash
# VSCodeで F5 → "Remote Debug + Hot Reload"
```

### 4. 開発開始！

- ファイルを編集して保存
- Airが自動でリビルド・再起動
- ブレークポイントは保持される
- デバッグが継続される

## 📝 設定ファイル詳細

### .air.toml

```toml
[build]
  # Delveデバッガー付きで起動
  cmd = "dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient"
  
  # 監視対象ファイル
  include_ext = ["go", "tpl", "tmpl", "html"]
  
  # 除外ディレクトリ
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
```

### Dockerfile.hotreload

```dockerfile
FROM golang:1.20-alpine

# Air + Delve インストール
RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# ホットリロード起動
CMD ["air", "-c", ".air.toml"]
```

### launch.json（VSCode）

```json
{
  "name": "Remote Debug + Hot Reload",
  "type": "go",
  "request": "attach",
  "mode": "remote",
  "remotePath": "/go/src/app",
  "port": 2345,
  "host": "localhost"
}
```

## 🔄 ワークフロー

### 通常の開発フロー

```
1. ファイル編集・保存
   ↓
2. Air が変更検知
   ↓
3. 自動リビルド
   ↓
4. Delve付きで再起動
   ↓
5. VSCode が自動再接続
   ↓
6. デバッグ継続
```

### デバッグ操作

- **ブレークポイント設定**: エディタの行番号左クリック
- **変数確認**: Variables パネルで確認
- **ステップ実行**: F10（ステップオーバー）、F11（ステップイン）
- **続行**: F5

## 💡 便利な機能

### 1. 条件付きブレークポイント

```go
// ブレークポイントを右クリック → 条件設定
if userId == "specific-user" {
    // この条件の時だけ停止
}
```

### 2. ログポイント

```go
// ブレークポイントの代わりにログ出力
// 右クリック → "Add Logpoint"
fmt.Printf("userId: {userId}, count: {count}")
```

### 3. 監視式（Watch）

```bash
# Variables パネルで変数を監視
len(users)
req.Header.Get("Authorization")
```

## 🔧 カスタマイズ

### Air設定のカスタマイズ

```toml
# より高速なリビルド
[build]
  delay = 500  # デフォルト: 1000ms
  
# 特定ファイルのみ監視
[build]
  include_file = ["main.go", "handler.go"]
```

### Delve起動オプション

```toml
# より詳細なログ
cmd = "dlv debug --headless --listen=:2345 --log --api-version=2"

# 特定パッケージのみデバッグ
cmd = "dlv debug ./cmd/api --headless --listen=:2345"
```

## 🐛 トラブルシューティング

### 1. Airが反応しない

```bash
# コンテナのログ確認
docker logs task-management-hotreload

# Air設定確認
docker exec -it task-management-hotreload cat .air.toml
```

### 2. デバッガー接続エラー

```bash
# Delveポート確認
docker port task-management-hotreload 2345

# Delveプロセス確認
docker exec -it task-management-hotreload ps aux | grep dlv
```

### 3. ホットリロードが遅い

```bash
# ボリュームマウント確認
docker-compose -f docker-compose.hotreload.yml config

# .air.toml の delay 設定を調整
delay = 500  # ミリ秒
```

## 📊 パフォーマンス比較

| 方式 | ビルド時間 | 再起動時間 | デバッグ接続 |
|------|------------|------------|-------------|
| ローカル実行 | 高速 | 高速 | 即座 |
| 通常Docker | 普通 | 遅い | 手動 |
| **Air + Docker** | **普通** | **高速** | **自動** |

## 🎉 メリット

### 開発効率
- **保存 → 即デバッグ**: 手動操作不要
- **ブレークポイント保持**: 再起動後も有効
- **本番環境模擬**: Dockerコンテナで実行

### チーム開発
- **環境統一**: 全員同じDocker環境
- **設定共有**: .air.toml等で設定統一
- **依存関係**: docker-composeで管理

## 🚨 注意点

### リソース使用量
- CPUとメモリを多く使用
- ファイル監視でI/O負荷
- 開発時のみ使用推奨

### ファイル同期
- Windowsでは同期が遅い場合がある
- 大きなファイルは除外推奨

## 🔗 参考資料

- [Air - ホットリロードツール](https://github.com/cosmtrek/air)
- [Delve - Goデバッガー](https://github.com/go-delve/delve)
- [VSCode Go拡張](https://code.visualstudio.com/docs/languages/go)
- [元記事: GoLand版](https://zenn.dev/smpeotn/articles/8ee4b6b14970ff)

## 🏁 まとめ

VSCode + Docker + Air + Delve の組み合わせにより、GoLandに匹敵する快適なデバッグ環境を実現できます。

**キーポイント:**
1. **Air**: ファイル変更の自動検知・リビルド
2. **Delve**: 強力なリモートデバッグ機能
3. **Docker**: 環境統一と依存関係管理
4. **VSCode**: 直感的なデバッグインターface

この環境で、Golang開発がより快適で効率的になります！