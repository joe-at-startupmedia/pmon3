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

func Group() *GroupRepo {
	dbInst := db.Db()
	groupOnce.Do(func() {
		if !dbInst.Migrator().HasTable(&model.Group{}) {
			dbInst.AutoMigrate(&model.Group{})
		}
	})
	return &GroupRepo{
		db: dbInst,
	}
}

func GroupOf(p *model.Group) *GroupRepo {
	gr := Group()
	gr.cur = p
	return gr
}

func (gr *GroupRepo) Create(name string) (*model.Group, error) {
	group := &model.Group{Name: name}
	err := gr.db.Save(group).Error
	return group, err
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
		// no soft-deletes
		return gr.db.Unscoped().Delete(g).Error
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

func (gr *GroupRepo) FindOrInsertByNames(names []string) ([]*model.Group, error) {
	if len(names) > 0 {
		groups := make([]*model.Group, len(names))
		for i := range names {
			name := names[i]
			var found model.Group
			err := gr.db.First(&found, "name = ?", name).Error
			if err != nil {
				pmond.Log.Infof("inserting group in database: %s %-v", name, err)
				group, err := gr.Create(name)
				if err != nil {
					pmond.Log.Infof("could not insert group in database: %s %-v", name, err)
					return nil, err
				} else {
					groups[i] = group
				}
			} else {
				groups[i] = &found
			}
		}

		return groups, nil
	} else {
		return nil, nil
	}
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
