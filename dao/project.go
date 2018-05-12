package dao

import (
	"strings"

	"github.com/labstack/gommon/random"

	"github.com/bojand/ghz-web/model"
	"github.com/jinzhu/gorm"
)

// ProjectService is our implementation
type ProjectService struct {
	DB *gorm.DB
}

// FindByID finds project by id
func (ps *ProjectService) FindByID(id uint, p *model.Project) error {
	return ps.DB.First(p, id).Error
}

// FindByName finds project by name
func (ps *ProjectService) FindByName(name string, p *model.Project) error {
	name = strings.ToLower(name)
	return ps.DB.First(p, "name = ?", name).Error
}

// Create creates a new project
func (ps *ProjectService) Create(p *model.Project) error {
	if strings.TrimSpace(p.Name) == "" {
		p.Name = random.String(16)
	}
	p.Name = strings.ToLower(p.Name)
	return ps.DB.Create(p).Error
}

// Update updates  project
func (ps *ProjectService) Update(p *model.Project) error {
	projToUpdate := model.Project{}
	err := ps.DB.First(&projToUpdate, p.ID).Error
	if err != nil {
		return err
	}

	if strings.TrimSpace(p.Name) == "" {
		p.Name = projToUpdate.Name
	}

	p.Name = strings.ToLower(p.Name)

	return ps.DB.Model(projToUpdate).Updates(p).Error
}
