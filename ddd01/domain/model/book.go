package model

import "time"

type Book struct {
	Id       int64
	Title    string
	Author   string
	IssuedAt time.Time
}
