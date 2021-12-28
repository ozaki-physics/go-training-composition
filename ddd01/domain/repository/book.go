package repository

import (
	"context"
	"github.com/ozaki-physics/go-training-composition/ddd01/domain/model"
)

// BookRepository インフラ層 が 実装をするインタフェース
// インフラ層 が ドメイン層のリポジトリ に依存
// DDD ではリポジトリは ドメインの集約だと思っているが参考にした記事は少し違う思想で作っているっぽい
// でもドメインの集約って実質 外部とのやり取りの定義だから インフラ層 のためのインタフェース定義は間違ってないかも
type BookRepository interface {
	GetAll(context.Context) ([]*model.Book, error)
}
