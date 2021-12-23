package usecase

import (
	"context"
	"github.com/ozaki-physics/go-training-composition/ddd01/domain/model"
	"github.com/ozaki-physics/go-training-composition/ddd01/domain/repository"
)

type BookUseCase interface {
	GetAll(context.Context) ([]*model.Book, error)
}

type bookUseCase struct {
	bookRepository repository.BookRepository
}

func NewBookUseCase(br repository.BookRepository) BookUseCase {
	return &bookUseCase{
		bookRepository: br,
	}
}

func (bu bookUseCase) GetAll(ctx context.Context) (books []*model.Book, err error){
	books, err = bu.bookRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return books, nil
}
