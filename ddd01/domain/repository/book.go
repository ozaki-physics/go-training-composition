package repository

import (
	"context"
	"github.com/ozaki-physics/go-training-composition/ddd01/domain/model"
)

type BookRepository interface {
	GetAll(context.Context) ([]*model.Book, error)
}
