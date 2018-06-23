package service

import "github.com/bojand/ghz-web/model"

// TestService is the interface for tests
type TestService interface {
	Count(pid uint) (uint, error)
	FindByID(id uint) (*model.Test, error)
	FindByName(name string) (*model.Test, error)
	FindByProjectID(pid uint, limit, page uint) ([]*model.Test, error)
	FindByProjectIDSorted(pid, num, page uint, sortField, order string) ([]*model.Test, error)
	Create(m *model.Test) error
	Update(m *model.Test) error
	Delete(m *model.Test) error
}
