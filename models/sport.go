package models

import (
	"gorm.io/gorm"
	//"fmt"
)

type Sport struct {
	gorm.Model

	Name string
	Type string
}

func (s *Sport) GetSportByName() (*Sport, error) {
	var err error
	err = DB.Debug().FirstOrCreate(&s, Sport{Name: s.Name, Type: s.Type}).Error
	if err != nil {
		return &Sport{}, err
	}
	return s, nil
}
