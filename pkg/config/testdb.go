//go:build !windows
// +build !windows

package config

import (
	"github.com/tonbiattack/go-task-management-api/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite"
)

// InitTestDB テスト用のインメモリSQLiteデータベースを初期化します
func InitTestDB() *gorm.DB {
	// modernc.org/sqlite を使用（Pure Go実装）
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("テストデータベースの接続に失敗しました: " + err.Error())
	}

	// テーブルの自動作成
	err = db.AutoMigrate(
		&model.Task{},
		&model.Workflow{},
		&model.WorkflowStep{},
		&model.TaskWorkflowStatus{},
	)
	if err != nil {
		panic("テストテーブルの作成に失敗しました: " + err.Error())
	}

	return db
}
