package service

import "github.com/bojand/ghz-web/model"

// DetailService is the interface for runs
type DetailService interface {
	Count(rid uint) (uint, error)
	FindByID(rid uint) (*model.Detail, error)
	FindByRunID(rid uint, limit, page uint) ([]*model.Detail, error)
	FindByRunIDSorted(rid, num, page uint, sortField, order string) ([]*model.Detail, error)
	Create(m *model.Detail) error
	Update(m *model.Detail) error
	Delete(m *model.Detail) error
	DeleteAll(rid uint) error
}
