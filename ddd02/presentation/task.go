package presentation

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ozaki-physics/go-training-composition/ddd02/usecase"
)

// TaskHandler プレゼンテーションのインタフェース
// わざわざ定義する理由は どんなメソッドが実装されているか分かりやすくしたり コンストラクタを経由して実装を強制させたり できるからっぽい
type TaskHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

// taskHandler プレゼンテーションの実体の構造体
type taskHandler struct {
	taskUsecase usecase.TaskUsecase
}

// requestTask タスクを作るときに使う
type requestTask struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type responseTask struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// NewTaskHandler プレゼンテーションのコンストラクタ
func NewTaskHandler(taskUsecase usecase.TaskUsecase) TaskHandler {
	return &taskHandler{taskUsecase: taskUsecase}
}

func (th *taskHandler) Get(w http.ResponseWriter, r *http.Request) {
	// URL パラメータ を string で取得
	idString := r.URL.Query().Get("id")
	// int に変換
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// ユースケースのメソッドでタスクを取得
	foundTask, err := th.taskUsecase.FindByID(id)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	res := responseTask{
		ID:      foundTask.ID,
		Title:   foundTask.Title,
		Content: foundTask.Content,
	}

	// http レスポンスに格納する
	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func (th *taskHandler) Post(w http.ResponseWriter, r *http.Request) {

}

func (th *taskHandler) Put(w http.ResponseWriter, r *http.Request) {

}

func (th *taskHandler) Delete(w http.ResponseWriter, r *http.Request) {

}
