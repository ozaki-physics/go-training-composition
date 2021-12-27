package model

import "time"

// Book 本当は値オブジェクトで 書籍がどういうモデルなのか表すと良い
type Book struct {
	Id       int64
	Title    string
	Author   string
	IssuedAt time.Time
}
