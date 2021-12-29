package usecase

import (
	"github.com/ozaki-physics/go-training-composition/ddd02/domain/model"
	"github.com/ozaki-physics/go-training-composition/ddd02/domain/repository"
)

// TaskUsecase ユースケースのインタフェース
// わざわざ定義する理由は どんなメソッドが実装されているか分かりやすくしたり コンストラクタを経由して実装を強制させたり できるからっぽい
// そして 各ユースケースメソッドの中で 実際に動くのはインフラ層のメソッドだけど
// ユースケース層からはリポジトリ層のメソッドを使っているように見えるだけなのが依存関係が美しい
// そのリポジトリ層は ユースケースのコンストラクタで依存性の注入をしている
type TaskUsecase interface {
	FindByID(id int) (*model.Task, error)
	Create(title, content string) (*model.Task, error)
	Update(id int, title, content string) (*model.Task, error)
	Delete(id int) error
}

// taskUsecase ユースケースの実体の構造体
type taskUsecase struct {
	taskRepo repository.TaskRepository
}

// NewTaskUsecase ユースケースのコンストラクタ
func NewTaskUsecase(taskRepo repository.TaskRepository) TaskUsecase {
	return &taskUsecase{taskRepo: taskRepo}
}

func (tu *taskUsecase) FindByID(id int) (*model.Task, error) {
	// リポジトリのメソッドでタスクを取得(実際に動くのは インフラ層のメソッド)
	foundTask, err := tu.taskRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return foundTask, nil
}

func (tu *taskUsecase) Create(title, content string) (*model.Task, error) {
	// ドメイン層のメソッドで 新しいタスクを生成
	task, err := model.NewTask(title, content)
	if err != nil {
		return nil, err
	}

	// リポジトリのメソッドで新しいタスクを保存(実際に動くのは インフラ層のメソッド)
	createTask, err := tu.taskRepo.Create(task)
	if err != nil {
		return nil, err
	}

	return createTask, nil
}

func (tu *taskUsecase) Update(id int, title, content string) (*model.Task, error) {
	// リポジトリのメソッドでタスクを取得(実際に動くのは インフラ層のメソッド)
	targetTask, err := tu.taskRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// ドメイン層のメソッドで モデルに更新用の値を格納
	err = targetTask.Set(title, content)
	if err != nil {
		return nil, err
	}

	// リポジトリのメソッドでタスクを更新(実際に動くのは インフラ層のメソッド)
	updateTask, err := tu.taskRepo.Update(targetTask)
	if err != nil {
		return nil, err
	}

	return updateTask, nil
}

func (tu *taskUsecase) Delete(id int) error {
	// リポジトリのメソッドでタスクを取得(実際に動くのは インフラ層のメソッド)
	task, err := tu.taskRepo.FindByID(id)
	if err != nil {
		return err
	}

	// リポジトリのメソッドでタスクを削除(実際に動くのは インフラ層のメソッド)
	err = tu.taskRepo.Delete(task)
	if err != nil {
		return err
	}

	return nil
}
