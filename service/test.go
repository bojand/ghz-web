package service

import "github.com/bojand/ghz-web/model"

// TestService is the interface for projects
type TestService interface {
	FindByID(id uint) (*model.Test, error)
	FindByName(name string) (*model.Test, error)
	FindByProjectID(pid uint, limit, page int) ([]*model.Test, error)
	Create(m *model.Test) error
	Update(m *model.Test) error
	Delete(m *model.Test) error
}
