package infra

import (
	"github.com/ozaki-physics/go-training-composition/ddd02/domain/model"
	"github.com/ozaki-physics/go-training-composition/ddd02/domain/repository"
)

// taskInfra インフラの実体の構造体
type taskInfra struct {
}

// NewTaskRepository task repository のコンストラクタ
// repository の実装 実体はインフラ層の構造体
func NewTaskRepository() repository.TaskRepository {
	return &taskInfra{}
}

// TODO: JSON ファイルに読み書きするように作ろう

func (tr *taskInfra) FindByID(id int) (*model.Task, error) {
	task := &model.Task{ID: id}

	// ダミーとして
	if err := task.Set("hello", "world"); err != nil {
		return nil, err
	}

	return task, nil
}

func (tr *taskInfra) Create(task *model.Task) (*model.Task, error) {

	return task, nil
}

func (tr *taskInfra) Update(task *model.Task) (*model.Task, error) {

	return task, nil
}

func (tr *taskInfra) Delete(task *model.Task) error {

	return nil
}
