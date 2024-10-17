package repo

import (
	"gorm.io/gorm"
	model2 "pmon3/model"
	"pmon3/pmond"
	"pmon3/pmond/db"
	"sync"
)

type GroupRepo struct {
	db  *gorm.DB
	cur *model2.Group
}

var groupOnce sync.Once

func Group() *GroupRepo {
	dbInst := db.Db()
	groupOnce.Do(func() {
		if !dbInst.Migrator().HasTable(&model2.Group{}) {
			dbInst.AutoMigrate(&model2.Group{})
		}
	})
	return &GroupRepo{
		db: dbInst,
	}
}

func GroupOf(p *model2.Group) *GroupRepo {
	gr := Group()
	gr.cur = p
	return gr
}

func (gr *GroupRepo) Create(name string) (*model2.Group, error) {
	group := &model2.Group{Name: name}
	err := gr.db.Save(group).Error
	return group, err
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

func (gr *GroupRepo) FindByIdOrName(idOrName string) (*model2.Group, error) {
	var found model2.Group
	err := gr.db.Preload("Processes").First(&found, "id = ? or name = ?", idOrName, idOrName).Error
	if err != nil {
		pmond.Log.Infof("could not find group in database: %s %-v", idOrName, err)
		return nil, err
	}
	return &found, nil
}

func (gr *GroupRepo) FindOrInsertByNames(names []string) ([]*model2.Group, error) {
	if len(names) > 0 {
		groups := make([]*model2.Group, len(names))
		for i := range names {
			name := names[i]
			var found model2.Group
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

func (gr *GroupRepo) FindAll() ([]model2.Group, error) {
	var all []model2.Group
	err := gr.db.Find(&all).Error
	if err != nil {
		pmond.Log.Infof("cant find groups: %v", err)
	}
	return all, err
}

func (gr *GroupRepo) AssignProcess(process *model2.Process) error {
	err := gr.db.Model(&gr.cur).Association("Processes").Append(process)
	if err != nil {
		pmond.Log.Infof("cant assign %s to %s: %-v", gr.cur.Name, process.Name, err)
		return err
	}
	return nil
}

func (gr *GroupRepo) RemoveProcess(process *model2.Process) error {
	err := gr.db.Model(&gr.cur).Association("Processes").Delete(process)
	if err != nil {
		pmond.Log.Infof("cant delete %s from %s: %-v", process.Name, gr.cur.Name, err)
		return err
	}
	return nil
}
