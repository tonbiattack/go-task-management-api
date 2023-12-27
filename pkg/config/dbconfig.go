package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// DB はアプリケーション全体で使用されるグローバルなデータベース接続です。
var DB *sql.DB

func init() {
	var err error

	dbUser := "root"
	dbPassword := ""
	dbHost := "127.0.0.1"
	dbPort := "3306"
	dbName := "task_management"

	// データソース名を作成します
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	// データベース接続を開きます
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("データベースを開く際のエラー: %v", err)
	}

	// データベース接続をチェックします
	err = DB.Ping()
	if err != nil {
		log.Fatalf("データベースのping時のエラー: %v", err)
	}
}

// GetDB はデータベース接続を返します
func GetDB() *sql.DB {
	return DB
}
