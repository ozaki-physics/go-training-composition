// model の構造体やビジネスロジックがおいてある
// 今回は title はブランクがダメということがわかる
package model

import (
	"errors"
)

type Task struct {
	ID      int
	Title   string
	Content string
}

// NewTask task のコンストラクタ
func NewTask(title, content string) (*Task, error) {
	if title == "" {
		return nil, errors.New("title を入力してください")
	}

	task := &Task{
		Title:   title,
		Content: content,
	}

	return task, nil
}

// Set task のセッター
func (t *Task) Set(title, content string) error {
	if title == "" {
		return errors.New("title を入力してください")
	}

	t.Title = title
	t.Content = content

	return nil
}
