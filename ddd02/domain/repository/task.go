package repository

import (
	"github.com/ozaki-physics/go-training-composition/ddd02/domain/model"
)

// TaskRepository task のインタフェース
// 実装は インフラ層
type TaskRepository interface {
	Create(task *model.Task) (*model.Task, error)
	FindByID(id int) (*model.Task, error)
	Update(task *model.Task) (*model.Task, error)
	Delete(task *model.Task) error
}
