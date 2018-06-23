package service

import "github.com/bojand/ghz-web/model"

// RunService is the interface for runs
type RunService interface {
	Count(pid uint) (uint, error)
	FindByID(id uint) (*model.Run, error)
	FindByTestID(pid uint, limit, page uint) ([]*model.Run, error)
	FindByTestIDSorted(pid, num, page uint, sortField, order string) ([]*model.Run, error)
	Create(m *model.Test) error
	Update(m *model.Test) error
	Delete(m *model.Test) error
}
