package repo

import (
	"gorm.io/gorm"
	"pmon3/pmond"
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
			dbInst.AutoMigrate(&model.Group{})
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

func (gr *GroupRepo) Create(name string) error {
	return gr.db.Save(&model.Group{Name: name}).Error
}

func (gr *GroupRepo) Save() error {
	return gr.db.Save(&gr.cur).Error
}

func (gr *GroupRepo) Delete(idOrName string) error {
	g, err := gr.FindByIdOrName(idOrName)
	if err != nil {
		pmond.Log.Infof("could not delete group in database: %s %-v", idOrName, err)
		return err
	} else if g != nil {
		return gr.db.Delete(g).Error
	}
	return nil
}

func (gr *GroupRepo) FindById(id uint32) (*model.Group, error) {
	var found model.Group
	err := gr.db.Preload("Processes").First(&found, id).Error
	if err != nil {
		pmond.Log.Infof("could not find group in database: %d %-v", id, err)
		return nil, err
	}
	return &found, nil
}

func (gr *GroupRepo) FindByIdOrName(idOrName string) (*model.Group, error) {
	var found model.Group
	err := gr.db.Preload("Processes").First(&found, "id = ? or name = ?", idOrName, idOrName).Error
	if err != nil {
		pmond.Log.Infof("could not find group in database: %s %-v", idOrName, err)
		return nil, err
	}
	return &found, nil
}

func (gr *GroupRepo) FindAll() ([]model.Group, error) {
	var all []model.Group
	err := gr.db.Find(&all).Error
	if err != nil {
		pmond.Log.Infof("cant find groups: %v", err)
	}
	return all, err
}

func (gr *GroupRepo) AssignProcess(process *model.Process) error {
	err := gr.db.Model(&gr.cur).Association("Processes").Append(process)
	if err != nil {
		pmond.Log.Infof("cant assign %s to %s: %-v", gr.cur.Name, process.Name, err)
		return err
	}
	return nil
}

func (gr *GroupRepo) RemoveProcess(process *model.Process) error {
	err := gr.db.Model(&gr.cur).Association("Processes").Delete(process)
	if err != nil {
		pmond.Log.Infof("cant delete %s from %s: %-v", process.Name, gr.cur.Name, err)
		return err
	}
	return nil
}
