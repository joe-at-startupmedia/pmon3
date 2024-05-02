package model

import "gorm.io/gorm"

type Pmond struct {
	gorm.Model
	Version string
}

func (Pmond) TableName() string {
	return "pmond"
}
