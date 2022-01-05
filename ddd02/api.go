// レイヤードアーキテクチャの練習
//
// See DDDを意識しながらレイヤードアーキテクチャとGoでAPIサーバーを構築する: https://qiita.com/ryokky59/items/6c2b35169fb6acafce15
// ただ DB の扱いがまだ分からないから ファイル IO で頑張ってみる(gorm ライブラリを使わない)
// http 通信も 標準ライブラリで頑張ってみる(echo ライブラリを使わない)
package ddd02

import (
	"fmt"
	"github.com/ozaki-physics/go-training-composition/ddd02/infra"
	"github.com/ozaki-physics/go-training-composition/ddd02/presentation"
	"github.com/ozaki-physics/go-training-composition/ddd02/usecase"
	"log"
	"net/http"
)

func MainApi() {
	taskRepository := infra.NewTaskRepository()
	taskUsecase := usecase.NewTaskUsecase(taskRepository)
	taskHandler := presentation.NewTaskHandler(taskUsecase)
	// プレゼンテーション層の出力が プレゼンテーション層の入力になって 同じ層内で依存しててよくない
	// ただ サンプルも実際は プレゼンテーション層のメソッドに 引数で渡しているから とりあえずこれでいく
	routes := presentation.NewRoutes(taskHandler)

	// ルーティングの定義
	http.HandleFunc("/", routes.InitRouting)
	// サーバ起動
	fmt.Println("Server Start >> http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

// See net/http でサーバーを立て、いくつかのパターンをパースしてみる: https://qiita.com/convto/items/2822d029349cb1b4df93
// http.HandleFunc に渡すときに引数で http メソッド指定とかできないかな

// エラーハンドリングの参考になりそうな記事
// See APIサーバのおけるGoのエラーハンドリングについて考えてみる: https://tutuz-tech.hatenablog.com/entry/2020/03/26/193519
