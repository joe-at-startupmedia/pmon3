package conf

import (
	"github.com/joe-at-startupmedia/xipc"
	"os"
	"pmon3/utils/file"
)

type FileOwnershipConfig struct {
	config        *Config
	User          string `yaml:"user,omitempty"`
	Group         string `yaml:"group,omitempty"`
	DirectoryMode string `yaml:"directory_mode" default:"0770"`
	FileMode      string `yaml:"file_mode" default:"0660"`
}

func (foc *FileOwnershipConfig) GetDirectoryMode() os.FileMode {
	if foc.DirectoryMode != "" {
		dirMode, err := file.ModeFromString(foc.DirectoryMode)
		if err == nil {
			return dirMode
		}
		foc.config.GetLogger().Errorf("error parsing DirectoryMode from %s with err %s", foc.DirectoryMode, err)
	}
	return os.FileMode(0770) //a safe default
}

func (foc *FileOwnershipConfig) GetFileMode() os.FileMode {
	if foc.FileMode != "" {
		fileMode, err := file.ModeFromString(foc.FileMode)
		if err == nil {
			return fileMode
		}
		foc.config.GetLogger().Errorf("error parsing FileMode from %s with err %s", foc.FileMode, err)
	}
	return os.FileMode(0660) //a safe default
}

func (c *Config) GetLogsFileOwnershipConfig() FileOwnershipConfig {
	lc := c.Logs
	foc := FileOwnershipConfig{
		config:        c,
		User:          lc.User,
		Group:         lc.Group,
		DirectoryMode: lc.DirectoryMode,
		FileMode:      lc.FileMode,
	}
	c.inheritFromDefaultPermissions(&foc)
	return foc
}

func (c *Config) GetDataFileOwnershipConfig() FileOwnershipConfig {
	dc := c.Data
	foc := FileOwnershipConfig{
		config:        c,
		User:          dc.User,
		Group:         dc.Group,
		DirectoryMode: dc.DirectoryMode,
		FileMode:      dc.FileMode,
	}
	c.inheritFromDefaultPermissions(&foc)
	return foc
}

func (c *Config) GetMessageQueueFileOwnershipConfig() FileOwnershipConfig {
	mqc := c.MessageQueue
	foc := FileOwnershipConfig{
		config:        c,
		User:          mqc.User,
		Group:         mqc.Group,
		DirectoryMode: mqc.DirectoryMode,
		FileMode:      mqc.FileMode,
	}
	c.inheritFromDefaultPermissions(&foc)
	return foc
}

func (c *Config) inheritFromDefaultPermissions(foc *FileOwnershipConfig) {
	if foc.User == "" {
		foc.User = c.Permissions.User
	}
	if foc.Group == "" {
		foc.Group = c.Permissions.Group
	}
	if foc.FileMode == "" {
		foc.FileMode = c.Permissions.FileMode
	}
	if foc.DirectoryMode == "" {
		foc.DirectoryMode = c.Permissions.DirectoryMode
	}
}

func (foc *FileOwnershipConfig) CreateDirectoryIfNonExistent(directory string) error {
	_, err := os.Stat(directory)
	dirMode := foc.GetDirectoryMode()
	if os.IsNotExist(err) {
		foc.config.GetLogger().Debugf("Creating directory %s with mode %s", directory, dirMode)
		if err = os.MkdirAll(directory, dirMode); err != nil {
			return err
		}
	}
	return foc.ApplyDirectoryPermissions(directory)
}

func (foc *FileOwnershipConfig) ApplyDirectoryPermissions(directory string) error {
	if _, err := os.Stat(directory); err != nil {
		return err
	}
	dirMode := foc.GetDirectoryMode()
	foc.config.GetLogger().Debugf("Applying permissions to directory (%s) with %-v", directory, dirMode)
	if foc.User != "" || foc.Group != "" {
		ownership := xipc.Ownership{
			Group:    foc.Group,
			Username: foc.User,
		}
		if err := ownership.ApplyPermissions(directory, int(dirMode)); err != nil {
			return err
		}
	} else if err := os.Chmod(directory, dirMode); err != nil {
		return err
	}

	return nil
}

func (foc *FileOwnershipConfig) ApplyFilePermissions(file string) error {
	if _, err := os.Stat(file); err != nil {
		return err
	}
	fileMode := foc.GetFileMode()
	foc.config.GetLogger().Debugf("Applying permissions to file (%s) with %-v", file, fileMode)
	if foc.User != "" || foc.Group != "" {
		ownership := xipc.Ownership{
			Group:    foc.Group,
			Username: foc.User,
		}
		if err := ownership.ApplyPermissions(file, int(fileMode)); err != nil {
			return err
		}
	} else if err := os.Chmod(file, fileMode); err != nil {
		return err
	}

	return nil
}
