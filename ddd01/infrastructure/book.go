package infrastructure

import (
	"context"
	"time"
	"github.com/ozaki-physics/go-training-composition/ddd01/domain/model"
	"github.com/ozaki-physics/go-training-composition/ddd01/domain/repository"
)

type bookInfrastructure struct {}

func NewBookInfrastructure() repository.BookRepository {
	return &bookInfrastructure{}
}

func (bi bookInfrastructure) GetAll(context.Context) ([]*model.Book, error){
	book01 := model.Book{}
	book01.Id = 1
	book01.Title = "クリーンアーキテクチャが分かる本"
	book01.Author = "たろう"
	book01.IssuedAt = time.Now().Add(-24*time.Hour)

	book02 := model.Book{}
	book02.Id = 2
	book02.Title = "レイヤードアーキテクチャが分かる本"
	book02.Author = "はなこ"
	book02.IssuedAt = time.Now().Add(-24*7*time.Hour)

	return []*model.Book{&book01, &book02}, nil
}
