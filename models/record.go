package models

import (
	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	ChatID  int
	SportID int
	Result  int
	Sport   Sport
}

func (r *Record) SaveRecord() (*Record, error) {
	var err error
	err = DB.Create(&r).Error
	if err != nil {
		return &Record{}, err
	}
	return r, nil
}
