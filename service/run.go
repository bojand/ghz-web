package service

import "github.com/bojand/ghz-web/model"

// RunService is the interface for runs
type RunService interface {
	Count(pid uint) (uint, error)
	FindByID(id uint) (*model.Run, error)
	FindByTestID(pid uint, limit, page uint) ([]*model.Run, error)
	FindByTestIDSorted(pid, num, page uint, sortField, order string) ([]*model.Run, error)
	Create(m *model.Run) error
	Update(m *model.Run) error
	Delete(m *model.Run) error
}
