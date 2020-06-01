package database

import (
	"github.com/jinzhu/gorm"
)

// IModel name of interface
type IModel interface {
	TableName() string
}

// IGormRepo base repository
type IGormRepo interface {
	Create(model IModel) error
	Update(model IModel) error
	Delete(model IModel) error
	FindOne(condition interface{}, model IModel) error
	GetDB() *gorm.DB
	InitDB(*gorm.DB) *GormRepo
}

// GormRepo model
type GormRepo struct {
	DB *gorm.DB
}

// NewGormRepo return a new drive instance
func NewGormRepo(db *gorm.DB) IGormRepo {
	return &GormRepo{DB: db}
}

// InitDB initialize a database
func (r *GormRepo) InitDB(db *gorm.DB) *GormRepo {
	r.DB = db
	return r
}

// Create func for create a new table
func (r *GormRepo) Create(model IModel) error {
	return r.DB.Create(model).Error
}

// FindOne func for find a single registry, use limit for that
func (r *GormRepo) FindOne(condition interface{}, model IModel) error {
	return r.DB.Where(condition).Limit(1).Find(model).Error
}

// Update func for upgrade a registry and write on database
func (r *GormRepo) Update(model IModel) error {
	return r.DB.Model(model).Update(model).Error
}

// Delete func remove registry from database
func (r *GormRepo) Delete(model IModel) error {
	return r.DB.Delete(model).Error
}

// GetDB return gorm engine
func (r *GormRepo) GetDB() *gorm.DB {
	return r.DB
}
