package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DBConfig はデータベース接続の設定情報を保持します
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// InitDB は新しいデータベース接続を初期化して返します
func InitDB() *gorm.DB {
	config := DBConfig{
		Host:     "127.0.0.1",
		Port:     "3306",
		User:     "root",
		Password: "",
		Name:     "task_management",
	}

	// DSN（Data Source Name）を作成
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.Name)

	// GORMを使ってデータベースに接続
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("データベース接続に失敗しました: %v", err)
	}

	// 接続を返す
	return db
}
