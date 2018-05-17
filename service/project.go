package service

import "github.com/bojand/ghz-web/model"

// ProjectService is the interface for projects
type ProjectService interface {
	FindByID(id uint) (*model.Project, error)
	FindByName(name string) (*model.Project, error)
	Create(p *model.Project) error
	Update(p *model.Project) error
	Delete(p *model.Project) error
}
