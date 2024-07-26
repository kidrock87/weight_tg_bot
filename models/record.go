package models

import (
	"errors"

	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	ChatID  int
	SportID int
	Result  int
	Sport   Sport
	Oresult string
}

func (r *Record) SaveRecord() (*Record, error) {
	var err error
	err = DB.Create(&r).Error
	if err != nil {
		return &Record{}, err
	}
	return r, nil
}

func (r *Record) FindRecordBySportAndChat(ChatId, SportId int) (*Record, error) {
	if err := DB.Debug().Where("sport_id = ? AND chat_id >= ?", SportId, ChatId).Preload("Sport").Find(&r).Error; err != nil {
		return r, errors.New("Record not found!")
	}
	return r, nil
}
