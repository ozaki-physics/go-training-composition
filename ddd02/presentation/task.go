package presentation

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/ozaki-physics/go-training-composition/ddd02/usecase"
)

// TaskHandler プレゼンテーションのインタフェース
// わざわざ定義する理由は どんなメソッドが実装されているか分かりやすくしたり コンストラクタを経由して実装を強制させたり できるからっぽい
type TaskHandler interface {
	Get(w http.ResponseWriter, r *http.Request) error
	Post(w http.ResponseWriter, r *http.Request) error
	Put(w http.ResponseWriter, r *http.Request) error
	Delete(w http.ResponseWriter, r *http.Request) error
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

func (th *taskHandler) Get(w http.ResponseWriter, r *http.Request) error {
	// URL パラメータ を string で取得
	idString := r.URL.Query().Get("id")
	// int に変換
	id, err := strconv.Atoi(idString)
	if err != nil {
		return err
	}

	// ユースケースのメソッドでタスクを取得
	foundTask, err := th.taskUsecase.FindByID(id)
	if err != nil {
		return err
	}
	res := responseTask{
		ID:      foundTask.ID,
		Title:   foundTask.Title,
		Content: foundTask.Content,
	}

	// http レスポンスに格納する
	if err = json.NewEncoder(w).Encode(res); err != nil {
		return err
	}

	return nil
	// 例
	// curl -X GET http://localhost:8080?id=10
}

func (th *taskHandler) Post(w http.ResponseWriter, r *http.Request) error {
	// Post されたデータを取り出す
	body := r.Body
	defer body.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, body)

	var req requestTask
	json.Unmarshal(buf.Bytes(), &req)

	// 新たに作る
	createTask, err := th.taskUsecase.Create(req.Title, req.Content)
	if err != nil {
		return err
	}

	// ドメインの構造体からレスポンスの構造体に詰め替える
	res := responseTask{
		ID:      createTask.ID,
		Title:   createTask.Title,
		Content: createTask.Content,
	}

	if err = json.NewEncoder(w).Encode(res); err != nil {
		return err
	}

	return nil
	// 例
	// curl -X POST -d '{"title":"insert", "content":"hoge"}' http://localhost:8080
	// bash で日本語入力すると文字化けが発生した
	// curl -X POST -H 'Content-Type: application/json;charset=utf-8' -d '{"title":"インサート", "content":"コンテンツ"}' http://localhost:8080
	// だから UTF-8 で保存した JSON を読み込ませる方法で POST すると 日本語も保存できた
	// (Windows 10のcurlで日本語postしたい時の対処法)[https://zenn.dev/unsoluble_sugar/articles/383a9f7e2c6a52cba2f2]
	// curl -X POST -H 'Content-Type: application/json' -d @ddd02/db/post.json http://localhost:8080
	// つまり go 側で対応することはとりあえず無いと仮定する
}

func (th *taskHandler) Put(w http.ResponseWriter, r *http.Request) error {
	// URL パラメータ を string で取得
	idString := r.URL.Query().Get("id")
	// int に変換
	id, err := strconv.Atoi(idString)
	if err != nil {
		return err
	}

	// Post されたデータを取り出す
	body := r.Body
	defer body.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, body)

	var req requestTask
	json.Unmarshal(buf.Bytes(), &req)

	// 更新する
	updateTask, err := th.taskUsecase.Update(id, req.Title, req.Content)
	if err != nil {
		return err
	}

	// ドメインの構造体からレスポンスの構造体に詰め替える
	res := responseTask{
		ID:      updateTask.ID,
		Title:   updateTask.Title,
		Content: updateTask.Content,
	}

	if err = json.NewEncoder(w).Encode(res); err != nil {
		return err
	}

	return nil
	// 例
	// curl -X PUT -d '{"title":"insert07", "content":"hoge07"}' http://localhost:8080?id=7
	// だから UTF-8 で保存した JSON を読み込ませる方法で PUT すると 日本語も保存できた
	// curl -X PUT -H 'Content-Type: application/json' -d @ddd02/db/post.json http://localhost:8080?id=7
}

func (th *taskHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	// URL パラメータ を string で取得
	idString := r.URL.Query().Get("id")
	// int に変換
	id, err := strconv.Atoi(idString)
	if err != nil {
		return err
	}

	if err := th.taskUsecase.Delete(id); err != nil {
		return err
	}

	return nil
	// 例
	// curl -X DELETE http://localhost:8080?id=4
}
