package infra

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/ozaki-physics/go-training-composition/ddd02/domain/model"
)

var fileName = "ddd02/db/dummy.json"

// taskJSON model.Task とは別で JSON ファイルの状態と合わせる JSON 専用の構造体を定義する
type taskJSON struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// readJSON JSON ファイルから読み出す
// 読み出すのはなるべく task.go の方でやって このファイル内で呼び出すことはしない
func readJSON() ([]taskJSON, error) {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var tasks []taskJSON
	if err := json.Unmarshal(bytes, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// searchJSON JSON ファイルから検索する
// tasks は  ポインタ(アドレス)を受け取る
// target は ポインタで渡ってくるから 直接代入しても参照元の構造体の値が変化する
func searchJSON(tasks *[]taskJSON, target *model.Task) (*model.Task, error) {
	// tasks 自体がアドレス(ポインタ)だから *tasks でもとの構造体を取得
	for _, t := range *tasks {
		if t.ID == target.ID {
			// JSON ファイルの構造体 から ドメインの構造体 に詰め替える
			// ID しか入っていないモデルの構造体に JSON ファイルから読み取った値を格納
			target.Set(t.Title, t.Content)
			break
		}
	}

	// これはプレゼンテーション層か?
	// error を返すのは正しいが どんな文言を表示するかはプレゼンテーション層の責務だと思われる
	// 該当なしの場合 model.Task が nil を返すだけで usecase 層や presentation 層 で対応しても良いかもしれない
	if target.Title == "" {
		err := errors.New("該当なし")
		return nil, err
	}

	// 引数をそのまま返すのは微妙?
	// だが もとの ID を含んだ構造体を返すため
	return target, nil
}

// writeJSON JSON ファイルの末尾に追加して保存
func writeJSON(tasks *[]taskJSON, target *model.Task) (*model.Task, error) {
	// ドメインの構造体 から JSON ファイルの構造体 に詰め替える
	insertTask := taskJSON{
		ID:      target.ID,
		Title:   target.Title,
		Content: target.Content,
	}

	// 読み込んだ全部入りの配列に追加する
	insertTasks := append(*tasks, insertTask)

	if err := saveJSON(&insertTasks); err != nil {
		return nil, err
	}

	return target, nil
}

// updateJSON JSON ファイル内のデータを更新
func updateJSON(tasks *[]taskJSON, target *model.Task) (*model.Task, error) {
	// ドメインの構造体 から JSON ファイルの構造体 に詰め替える
	updateTask := taskJSON{
		ID:      target.ID,
		Title:   target.Title,
		Content: target.Content,
	}

	// tasks 自体がアドレス(ポインタ)だから *tasks でもとの構造体を取得
	for i, t := range *tasks {
		if t.ID == target.ID {
			// ポインタ元(もとの構造体)の要素を変更するには (*task) とカッコを付ける必要がある
			(*tasks)[i] = updateTask
			break
		}
	}

	if err := saveJSON(tasks); err != nil {
		return nil, err
	}

	// どんな値で更新したか明らかにするために そのまま実際に使った値を返す
	return target, nil
}

// deleteJSON JSON ファイル内のデータを削除
func deleteJSON(tasks *[]taskJSON, target *model.Task) error {
	// golang では 配列の削除メソッドが存在しないらしいから 新たに配列を作る
	for i, t := range *tasks {
		if t.ID == target.ID {
			newTasks := append((*tasks)[:i], (*tasks)[i+1:]...)
			if err := saveJSON(&newTasks); err != nil {
				return err
			}
			break
		}
	}
	return nil
}

// saveJSON JSON ファイルに保存
// infra 層の task.go のメソッド -> infra 層の handleJSON.go のメソッド -> saveJSON() と
// 階層が深いが 本来は infra 層の task.go に全部書くべきと思われる
func saveJSON(tasks *[]taskJSON) error {
	// 全部入りの配列をファイル保存用の bytes に変換
	bytes, err := json.Marshal(*tasks)
	if err != nil {
		return err
	}

	// ファイルを開く(write 権限付き)
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// 全部入り配列を一気に書き込む
	if _, err := file.Write(bytes); err != nil {
		return err
	}

	return nil
}
