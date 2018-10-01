package service

import "github.com/bojand/ghz-web/model"

// RunService is the interface for runs
type RunService interface {
	Count(tid uint) (uint, error)
	FindLatest(tid uint) (*model.Run, error)
	FindByID(id uint) (*model.Run, error)
	FindByTestID(tid uint, limit, page uint, populate bool) ([]*model.Run, error)
	FindByTestIDSorted(tid, num, page uint, sortField, order string, histogram bool, latency bool) ([]*model.Run, error)
	Create(m *model.Run) error
	Update(m *model.Run) error
	Delete(m *model.Run) error
}
