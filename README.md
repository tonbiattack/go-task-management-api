# Go タスク管理 API

- Go タスク管理 API は、基本的なタスク管理機能を提供する RESTful API です。
- Go 言語の勉強として作成しました。
- この API を使用することで、タスクの作成、取得、更新、削除が可能です。
- 期限日を加えようとした時の既存の処理のリファクタのめんどくささと、json の parse などの定形処理がめんどくさかったので、フレームワーク Gin と OR マッパーの GORM ORM を導入しました。

## 2024/01/13 追加分

- **ワークフロー機能の追加**: ワークフロー機能を導入しました。この機能はタスク管理プロセスを柔軟に制御し、複数のステップを含むタスクの進行を追跡するのに役立ちます。
- **Postman ワークスペースの追加**: API のテストとドキュメンテーションのための Postman ワークスペースを追加しました。

# 使用技術

- Go
- Gin
- GORM ORM
- MySQL
- HTTP サーバー (標準ライブラリ)
- Gorilla Mux（ルーティング）
- Air

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

# サーバーの起動

```
go run main.go
```

# 注意事項

- この API はテストを含んでいません。
