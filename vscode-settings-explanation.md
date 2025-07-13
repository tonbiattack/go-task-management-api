# VSCode Settings, Tasks, Launch 設定説明

このGoプロジェクトのVSCode設定ファイルについて説明します。

## Settings (`.vscode/settings.json`)

Goプロジェクト用の基本設定：

### Go言語設定
- `go.toolsManagement.autoUpdate: true` - Goツールの自動更新
- `go.useLanguageServer: true` - 言語サーバーの使用
- `go.buildOnSave: "package"` - 保存時にパッケージをビルド
- `go.lintOnSave: "package"` - 保存時にlint実行
- `go.vetOnSave: "package"` - 保存時にvet実行
- `go.formatTool: "goimports"` - フォーマッターにgoimportsを使用
- `go.lintTool: "golint"` - リンターにgolintを使用

### エディタ設定
- `editor.formatOnSave: true` - 保存時に自動フォーマット
- `editor.codeActionsOnSave.source.organizeImports: "explicit"` - 保存時にimport整理
- `files.eol: "\n"` - 改行文字をLFに統一
- `files.insertFinalNewline: true` - ファイル末尾に改行を挿入
- `files.trimTrailingWhitespace: true` - 行末の空白を削除

### デバッグ設定
- `go.delveConfig` - Delveデバッガーの詳細設定
  - 変数の表示制限、文字列長制限などを設定

## Tasks (`.vscode/tasks.json`)

自動化タスクの定義：

### Docker関連タスク
- `start-mysql-only` - MySQLコンテナのみ起動
- `start-debug-container` - デバッグ用コンテナ起動
- `stop-debug-container` - デバッグ用コンテナ停止
- `rebuild-debug-container` - デバッグ用コンテナ再構築
- `start-hotreload-container` - ホットリロード用コンテナ起動
- `stop-hotreload-container` - ホットリロード用コンテナ停止
- `rebuild-hotreload-container` - ホットリロード用コンテナ再構築
- `stop-docker` - 全Dockerコンテナ停止

### Go関連タスク
- `build-go` - Goアプリケーションのビルド（main.exe生成）
- `test-go` - Goテストの実行（`go test -v ./...`）

## Launch Configurations (`.vscode/launch.json`)

デバッグ実行設定：

### 1. Launch Go App (Local)
- ローカル環境でのアプリケーション実行
- MySQL接続先：localhost:3306
- 開発環境設定

### 2. Launch Go App (Docker DB)
- Dockerデータベース使用でのアプリケーション実行
- MySQL接続先：localhost:3307
- 事前タスク：`start-mysql-only`を実行

### 3. Remote Debug Docker Container
- Dockerコンテナ内のアプリケーションにリモートデバッグ接続
- ポート：2345でdelveデバッガーに接続
- パスマッピング：`${workspaceFolder}` → `/app`
- 事前タスク：`start-debug-container`を実行

### 4. Remote Debug (Manual Container)
- 手動で起動したコンテナへのリモートデバッグ接続
- 事前タスクなし（手動でコンテナ起動が必要）

### 5. Remote Debug + Hot Reload
- ホットリロード機能付きリモートデバッグ
- パスマッピング：`${workspaceFolder}` → `/go/src/app`
- 事前タスク：`start-hotreload-container`を実行

### 6. Debug Tests
- テストのデバッグ実行
- テスト専用データベース使用（task_management_test）
- `-test.v`フラグでverbose出力

## 開発ワークフロー

1. **ローカル開発**：Launch Go App (Local)を使用
2. **Docker環境でのテスト**：Launch Go App (Docker DB)を使用
3. **本格的なデバッグ**：Remote Debug Docker Containerを使用
4. **ホットリロード開発**：Remote Debug + Hot Reloadを使用
5. **テストデバッグ**：Debug Testsを使用

各設定はDockerコンテナとの連携を重視しており、開発環境の柔軟性を提供しています。