package service

import "github.com/bojand/ghz-web/model"

// TestService is the interface for projects
type TestService interface {
	FindByID(id uint, m *model.Test) error
	FindByName(name string, m *model.Test) error
	Create(m *model.Test) error
	Update(m *model.Test) error
	Delete(m *model.Test) error
}
