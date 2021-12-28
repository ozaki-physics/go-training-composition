// レイヤードアーキテクチャの練習
//
// See 【Go + レイヤードアーキテクチャー】DDDを意識してWeb APIを実装してみる: https://yyh-gl.github.io/tech-blog/blog/go_web_api/
package ddd01

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ozaki-physics/go-training-composition/ddd01/infrastructure"
	"github.com/ozaki-physics/go-training-composition/ddd01/presentation"
	"github.com/ozaki-physics/go-training-composition/ddd01/usecase"
)

// MainApi レイヤードアーキテクチャのサーバを立てる
func MainApi() {
	// 依存性の注入っぽいことをしている
	bookInfrastructure := infrastructure.NewBookInfrastructure()
	bookUseCase := usecase.NewBookUseCase(bookInfrastructure)
	bookPresentation := presentation.NewBookPresentation(bookUseCase)

	// ルーティングの整理
	http.HandleFunc("/api/v1/books", bookPresentation.ServerHTTP)
	// サーバ起動
	fmt.Println("Server Start >> http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
