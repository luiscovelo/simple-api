package dao

import (
	"simple-api/internal/model"

	"gorm.io/gorm"
)

type Dao struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Dao {
	return &Dao{
		db: db,
	}
}

func (d *Dao) Save(msg *model.Message) error {
	return d.db.Create(msg).Error
}

func (d *Dao) Get(id uint) (*model.Message, error) {
	var msg model.Message
	if err := d.db.First(&msg, id).Error; err != nil {
		return nil, err
	}
	return &msg, nil
}
