package usecase

import (
	"context"
	"github.com/ozaki-physics/go-training-composition/ddd01/domain/model"
	"github.com/ozaki-physics/go-training-composition/ddd01/domain/repository"
)

// BookUseCase ユースケース層 が 実装をするインタフェース
// ユースケース層 が ユースケース層 に依存
type BookUseCase interface {
	GetAll(context.Context) ([]*model.Book, error)
}

type bookUseCase struct {
	bookRepository repository.BookRepository
}

// NewBookUseCase ユースケース層 の構造体のポインタを返す
// 戻り値は ユースケース層
// つまり ユースケース層のインタフェース を実装してないと戻り値になれない
func NewBookUseCase(br repository.BookRepository) BookUseCase {
	return &bookUseCase{
		bookRepository: br,
	}
}

// GetAll ユースケース層の実装
func (bu bookUseCase) GetAll(ctx context.Context) (books []*model.Book, err error) {
	// リポジトリ層のインタフェースを使う
	// ユースケース層 は リポジトリ層 に依存している
	// 実体は インタフェース層 で 依存性の逆転ができている
	books, err = bu.bookRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return books, nil
}
