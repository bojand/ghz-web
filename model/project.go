package model

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/random"
)

// Project represents a project
type Project struct {
	gorm.Model
	Name        string `json:"name" gorm:"unique_index"`
	Description string `json:"description"`
}

// BeforeCreate is a GORM hook called when a model is created
func (p *Project) BeforeCreate() error {
	name := strings.Replace(p.Name, " ", "", -1)
	if name == "" {
		name = random.String(16)
	}
	p.Name = strings.ToLower(name)
	p.Description = strings.TrimSpace(p.Description)

	return nil
}

// BeforeUpdate is a GORM hook called when a model is updated
func (p *Project) BeforeUpdate() error {
	name := strings.Replace(p.Name, " ", "", -1)
	if name == "" {
		return errors.New("Project name cannot be empty")
	}

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
	return p, err
}

// FindByName finds project by name
func (ps *ProjectService) FindByName(name string) (*Project, error) {
	name = strings.ToLower(name)
	p := new(Project)
	err := ps.DB.First(p, "name = ?", name).Error
	return p, err
}

// Create creates a new project
func (ps *ProjectService) Create(p *Project) error {
	return ps.DB.Create(p).Error
}

// Update updates  project
func (ps *ProjectService) Update(p *Project) error {
	projToUpdate := &Project{}
	err := ps.DB.First(projToUpdate, p.ID).Error
	if err != nil {
		return err
	}

	name := strings.Replace(p.Name, " ", "", -1)
	if name == "" {
		p.Name = projToUpdate.Name
	}

	return ps.DB.Model(p).Updates(p).Error
}

// Delete deletes project
func (ps *ProjectService) Delete(p *Project) error {
	return errors.New("Not Implemented Yet")
}
