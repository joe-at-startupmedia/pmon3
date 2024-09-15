package repo

import (
	"gorm.io/gorm"
	"pmon3/pmond/db"
	"pmon3/pmond/model"
	"sync"
)

type GroupRepo struct {
	db  *gorm.DB
	cur *model.Group
}

var groupOnce sync.Once
var groupRepo *GroupRepo

func Group() *GroupRepo {
	groupOnce.Do(func() {
		dbInst := db.Db()
		if !dbInst.Migrator().HasTable(&model.Group{}) {
			dbInst.Migrator().CreateTable(&model.Group{})
		}
		groupRepo = &GroupRepo{
			db: dbInst,
		}
	})
	groupRepo.cur = nil
	return groupRepo
}

func GroupOf(p *model.Group) *GroupRepo {
	gr := Group()
	gr.cur = p
	return gr
}
