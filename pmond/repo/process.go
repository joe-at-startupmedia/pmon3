package repo

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"pmon3/pmond"
	"pmon3/pmond/db"
	"pmon3/pmond/model"
	"sync"
)

type ProcessRepo struct {
	db  *gorm.DB
	cur *model.Process
}

var processOnce sync.Once

func Process() *ProcessRepo {
	dbInst := db.Db()
	processOnce.Do(func() {
		if !dbInst.Migrator().HasTable(&model.Process{}) {
			dbInst.Migrator().CreateTable(&model.Process{})
		}
	})
	return &ProcessRepo{
		db: dbInst,
	}
}

func ProcessOf(p *model.Process) *ProcessRepo {
	pr := Process()
	pr.cur = p
	return pr
}

func (pr *ProcessRepo) Save() error {
	return pr.db.Save(&pr.cur).Error
}

func (pr *ProcessRepo) Delete() error {
	return pr.db.Delete(&pr.cur).Error
}

func (pr *ProcessRepo) UpdateStatus(status model.ProcessStatus) error {
	pr.cur.Status = status
	return pr.Save()
}

func (pr *ProcessRepo) FindById(id uint32) (*model.Process, error) {
	var found model.Process
	err := pr.db.First(&found, id).Error
	if err != nil {
		pmond.Log.Infof("could not find process in database: %d %-v", id, err)
		return nil, err
	}
	return &found, nil
}

func (pr *ProcessRepo) FindByIdOrName(idOrName string) (*model.Process, error) {
	var found model.Process
	err := pr.db.First(&found, "id = ? or name = ?", idOrName, idOrName).Error
	if err != nil {
		pmond.Log.Infof("could not find process in database: %s %-v", idOrName, err)
		return nil, err
	}
	return &found, nil
}

func (pr *ProcessRepo) FindByFileAndName(processFile string, name string) (*model.Process, error) {
	var found model.Process
	err := pr.db.First(&found, "process_file = ? AND name = ?", processFile, name).Error
	if err != nil {
		pmond.Log.Infof("could not find process in database: %s or %s %-v", processFile, name, err)
		return nil, err
	}
	return &found, nil
}

func (pr *ProcessRepo) FindByStatus(status model.ProcessStatus) ([]model.Process, error) {
	var all []model.Process
	err := db.Db().Find(&all, "status = ?", status).Error
	if err != nil {
		pmond.Log.Infof("pmon3 can find processes with status %s: %v", status, err)
	}
	return all, err
}

func (pr *ProcessRepo) FindForMonitor() ([]model.Process, error) {
	var all []model.Process
	err := db.Db().Find(&all, "status in (?, ?, ?, ?, ?)",
		model.StatusRunning,
		model.StatusFailed,
		model.StatusQueued,
		model.StatusClosed,
		model.StatusBackoff,
	).Error
	return all, err
}

func (pr *ProcessRepo) FindAll() ([]model.Process, error) {
	var all []model.Process
	err := pr.db.Find(&all).Error
	if err != nil {
		pmond.Log.Infof("cant find processes: %v", err)
	}
	return all, err
}

func (pr *ProcessRepo) FindAndSave() (string, error) {

	originOne, err := pr.FindByFileAndName(pr.cur.ProcessFile, pr.cur.Name)
	if err == nil && originOne.ID > 0 { // process already exists
		pr.cur.ID = originOne.ID
	}

	err = pr.Save()
	if err != nil {
		return "", fmt.Errorf("could not save process: err: %w", err)
	}
	output, err := json.Marshal(pr.cur.RenderTable())
	return string(output), err
}
