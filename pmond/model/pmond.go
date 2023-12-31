package model

import "github.com/jinzhu/gorm"

type Pmond struct {
	gorm.Model
	Version string
}

func (Pmond) TableName() string {
	return "pmond"
}
