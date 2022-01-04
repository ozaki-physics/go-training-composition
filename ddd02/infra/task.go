package infra

import (
	"errors"

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

func (tr *taskInfra) FindByID(id int) (*model.Task, error) {
	target := &model.Task{ID: id}

	readTasks, err := readJSON()
	if err != nil {
		return nil, err
	}

	// &readTasks は データ量が多いとコピーするのが大変だから ポインタ(アドレス)を渡す
	task, err := searchJSON(&readTasks, target)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (tr *taskInfra) Create(task *model.Task) (*model.Task, error) {
	// ID が ゼロ値の 0 なら新規作成と判断する
	if task.ID != 0 {
		err := errors.New("同じ ID が既に存在します")
		return nil, err
	}

	readTasks, err := readJSON()
	if err != nil {
		return nil, err
	}

	// ID の採番
	// ID が一意にならないといけないのは インフラ都合と考える?
	lastTask := readTasks[len(readTasks)-1]
	task.ID = lastTask.ID + 1

	insertTask, err := writeJSON(&readTasks, task)
	if err != nil {
		return nil, err
	}

	return insertTask, nil
}

func (tr *taskInfra) Update(task *model.Task) (*model.Task, error) {
	// usecase で探して 構造体の更新までしているから ただ保存するだけでよい
	// ただ インフラ都合として 同じ ID で保存しないように調整する
	readTasks, err := readJSON()
	if err != nil {
		return nil, err
	}

	updateTask, err := updateJSON(&readTasks, task)
	if err != nil {
		return nil, err
	}

	return updateTask, nil
}

func (tr *taskInfra) Delete(task *model.Task) error {
	// usecase で探して 構造体の特定までしているから ただ削除するだけでよい
	readTasks, err := readJSON()
	if err != nil {
		return err
	}

	if err := deleteJSON(&readTasks, task); err != nil {
		return err
	}

	return nil
}
