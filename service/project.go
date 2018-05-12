package service

import "github.com/bojand/ghz-web/model"

// ProjectService is the interface for projects
type ProjectService interface {
	FindByID(id uint, p *model.Project) error
	FindByName(name string, p *model.Project) error
	Create(p *model.Project) error
	Update(p *model.Project) error
}
