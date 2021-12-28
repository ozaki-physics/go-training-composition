// 外部と ユースケース層 の緩衝材
// 今回は http 通信だが プレゼンテーション層 を変えるだけで CLI 対応もできるはず
package presentation

import (
	"encoding/json"
	"github.com/ozaki-physics/go-training-composition/ddd01/usecase"
	"net/http"
	"time"
)

// BookPresentation プレゼンテーション層 が実装するインタフェース
// プレゼンテーション層 が プレゼンテーション層 に依存
type BookPresentation interface {
	ServerHTTP(w http.ResponseWriter, r *http.Request)
}

type bookPresentation struct {
	bookUseCase usecase.BookUseCase
}

// NewBookPresentation プレゼンテーション層 の構造体のポインタを返す
// 戻り値は プレゼンテーション層
// つまり プレゼンテーション層のインタフェース を実装していないと戻り値になれない
func NewBookPresentation(bu usecase.BookUseCase) BookPresentation {
	return &bookPresentation{
		bookUseCase: bu,
	}
}

// ServerHTTP HTTP リクエストを受け取り ユースケース層 を使って処理をして 結果を返す
func (bp bookPresentation) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	// bookField ドメインモデルの構造体に http の関心事である JSON タグを付与したくないため プレゼンテーション層で用意したらしい
	// 正直2重管理では?って思ってしまう
	type bookField struct {
		Id       int64     `json:"id"`
		Title    string    `json:"title"`
		Author   string    `json:"author"`
		IssuedAt time.Time `json:"issued_at"`
	}

	// response API のレスポンス
	type response struct {
		Books []bookField `json:"books"`
	}

	// リクエストのコンテキスト
	ctx := r.Context()

	// ユースケース層 の呼び出し
	books, err := bp.bookUseCase.GetAll(ctx)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// 取得したドメインモデルを レスポンス に変換
	res := new(response)
	for _, book := range books {
		bf := bookField(*book)
		// 構造体のフィールドに追加
		res.Books = append(res.Books, bf)
	}

	// クライアントにレスポンスを返却
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
}
