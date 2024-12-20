package infra

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupDB initializes the database connection and creates tables if they don't exist
func SetupDB() *gorm.DB {
	// データソース名を環境変数から取得
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// データベースに接続
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// マイグレーション
	initSQL := `
	-- ユーザーテーブル
	CREATE TABLE IF NOT EXISTS users (
		id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		icon_url VARCHAR(255),
		google_id VARCHAR(50) UNIQUE,
		qiita_id VARCHAR(50) UNIQUE,
		Qiita_link BOOLEAN DEFAULT FALSE
	);

	-- 記事テーブル
	CREATE TABLE IF NOT EXISTS articles (
		id SERIAL PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		title VARCHAR(50) NOT NULL,
		main_md TEXT NOT NULL,
		slide_md TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		like_count INT DEFAULT 0,
		public BOOLEAN DEFAULT FALSE,
		qiita_article BOOLEAN DEFAULT FALSE
	);

	-- updated_atフィールドを自動更新するための関数を作成
	CREATE OR REPLACE FUNCTION update_updated_at()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = CURRENT_TIMESTAMP;
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	-- トリガーが既に存在していない場合のみ作成する
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'set_updated_at') THEN
			CREATE TRIGGER set_updated_at
			BEFORE UPDATE ON articles
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at();
		END IF;
	END $$;

	-- 記事の「いいね」テーブル
	CREATE TABLE IF NOT EXISTS articleLikes (
		id SERIAL PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		article_id INT NOT NULL REFERENCES articles(id) ON DELETE CASCADE
	);

	-- タグテーブル
	CREATE TABLE IF NOT EXISTS tags (
		id SERIAL PRIMARY KEY,
		word VARCHAR(20) UNIQUE NOT NULL
	);

	-- 記事タグ関連テーブル
	CREATE TABLE IF NOT EXISTS articleTagRelations (
		id SERIAL PRIMARY KEY,
		article_id INT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
		tag_id INT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
		UNIQUE (article_id, tag_id)
	);
`

	// SQL実行
	if err := db.Exec(initSQL).Error; err != nil {
		log.Fatalf("Error executing initSQL: %v", err)
	}

	return db
}
