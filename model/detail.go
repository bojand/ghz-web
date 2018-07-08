package model

import (
	"errors"
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/jinzhu/gorm"
)

// Detail represents a run detail
type Detail struct {
	Model
	RunID   uint    `json:"runID" gorm:"type:integer REFERENCES runs(id)"`
	Run     *Run    `json:"-"`
	Latency float64 `json:"latency" validate:"required"`
	Error   string  `json:"error"`
	Status  string  `json:"status"`
}

// BeforeSave is called by GORM before save
func (d *Detail) BeforeSave(scope *gorm.Scope) error {
	if d.RunID == 0 && d.Run == nil {
		return errors.New("Run must belong to a test")
	}

	d.Error = strings.TrimSpace(d.Error)

	status := strings.TrimSpace(d.Status)
	if status == "" {
		status = "OK"
	}
	d.Status = status

	if scope != nil {
		scope.SetColumn("error", d.Error)
		scope.SetColumn("status", d.Status)
	}

	return nil
}

// DetailService is our implementation
type DetailService struct {
	DB *gorm.DB
}

// Create creates a new run
func (ds *DetailService) Create(r *Detail) error {
	return ds.DB.Create(r).Error
}

// Count returns the total number of runs
func (ds *DetailService) Count(rid uint) (uint, error) {
	count := uint(0)
	err := ds.DB.Model(&Detail{}).Where("run_id = ?", rid).Count(&count).Error
	return count, err
}

// FindByID finds test by id
func (ds *DetailService) FindByID(id uint) (*Detail, error) {
	d := new(Detail)
	err := ds.DB.First(d, id).Error
	if err != nil {
		d = nil
	}
	return d, err
}

// FindByRunID finds tests by project
func (ds *DetailService) FindByRunID(rid, num, page uint) ([]*Detail, error) {
	r := &Run{}
	r.ID = rid

	offset := uint(0)
	if page >= 0 && num >= 0 {
		offset = page * num
	}

	s := make([]*Detail, num)

	err := ds.DB.Offset(offset).Limit(num).Order("id desc").Model(r).Related(&s).Error

	return s, err
}

// FindByRunIDSorted lists tests using sorting
func (ds *DetailService) FindByRunIDSorted(rid, num, page uint, sortField, order string) ([]*Detail, error) {
	if (sortField != "id" && sortField != "latency") || (order != "asc" && order != "desc") {
		return nil, errors.New("Invalid sort parameters")
	}

	offset := uint(0)
	if page >= 0 && num >= 0 {
		offset = page * num
	}

	orderSQL := sortField + " " + order

	r := &Run{}
	r.ID = rid

	s := make([]*Detail, num)

	err := ds.DB.Order(orderSQL).Offset(offset).Limit(num).Model(r).Related(&s).Error

	return s, err
}

// Update updates a run
func (ds *DetailService) Update(d *Detail) error {
	dToUpdate := &Detail{}
	if err := ds.DB.First(dToUpdate, d.ID).Error; err != nil {
		return err
	}

	return ds.DB.Save(d).Error
}

// Delete deletes a detial
func (ds *DetailService) Delete(r *Detail) error {
	return errors.New("Not Implemented Yet")
}

// DeleteAll deletes all details for a run
func (ds *DetailService) DeleteAll(rid uint) error {
	return errors.New("Not Implemented Yet")
}

// CreateBatch creates a batch of details returning the number successfully created,
// and the number failed
func (ds *DetailService) CreateBatch(rid uint, s []*Detail) (uint, uint) {
	NC := 1
	nReq := len(s)

	var nErr uint32

	sem := make(chan bool, NC)

	for _, item := range s {
		sem <- true

		go func(d *Detail) {
			defer func() { <-sem }()

			d.RunID = rid
			err := ds.Create(d)

			if err != nil {
				fmt.Println(err)
				atomic.AddUint32(&nErr, 1)
			}
		}(item)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	errCount := uint(atomic.LoadUint32(&nErr))
	nCreated := uint(nReq) - errCount

	return nCreated, errCount
}
