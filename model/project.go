package model

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/random"
)

// Project represents a project
type Project struct {
	Model
	Name        string `json:"name" gorm:"unique_index;not null"`
	Description string `json:"description"`
}

// BeforeCreate is a GORM hook called when a model is created
func (p *Project) BeforeCreate() error {
	if p.Name == "" {
		p.Name = random.String(16)
	}

	return nil
}

// BeforeUpdate is a GORM hook called when a model is updated
func (p *Project) BeforeUpdate() error {
	if p.Name == "" {
		return errors.New("Project name cannot be empty")
	}

	return nil
}

// BeforeSave is a GORM hook called when a model is created
func (p *Project) BeforeSave() error {
	name := strings.Replace(p.Name, " ", "", -1)
	p.Name = strings.ToLower(name)
	p.Description = strings.TrimSpace(p.Description)

	return nil
}

// ProjectService is our implementation
type ProjectService struct {
	DB *gorm.DB
}

// FindByID finds project by id
func (ps *ProjectService) FindByID(id uint) (*Project, error) {
	p := new(Project)
	err := ps.DB.First(p, id).Error
	if err != nil {
		p = nil
	}
	return p, err
}

// FindByName finds project by name
func (ps *ProjectService) FindByName(name string) (*Project, error) {
	name = strings.ToLower(name)
	p := new(Project)
	err := ps.DB.First(p, "name = ?", name).Error
	if err != nil {
		p = nil
	}
	return p, err
}

// Create creates a new project
func (ps *ProjectService) Create(p *Project) error {
	return ps.DB.Create(p).Error
}

// Update updates  project
func (ps *ProjectService) Update(p *Project) error {
	projToUpdate := &Project{}
	if err := ps.DB.First(projToUpdate, p.ID).Error; err != nil {
		return err
	}

	name := strings.Replace(p.Name, " ", "", -1)
	if name == "" {
		p.Name = projToUpdate.Name
	}

	return ps.DB.Save(p).Error
}

// Delete deletes project
func (ps *ProjectService) Delete(p *Project) error {
	return errors.New("Not Implemented Yet")
}

// List lists projects
func (ps *ProjectService) List(limit, page uint) ([]*Project, error) {
	s := make([]*Project, 0)

	offset := uint(0)
	if page >= 0 && limit >= 0 {
		offset = page * limit
	}

	err := ps.DB.Offset(offset).Limit(limit).Order("name desc").Find(&s).Error

	return s, err
}
