package sheetmerger

import (
	config "github.com/kayac/go-config"
	"github.com/pkg/errors"
)

type ClientConfig struct {
	Email        string `yaml:"email"`
	PrivateKeyID string `yaml:"private_key_id"`
	PrivateKey   string `yaml:"private_key"`
}

type Config struct {
	Client           ClientConfig `yaml:"client"`
	BaseSheetName    string       `yaml:"base_sheet_name"`
	SheetIndexColumn string       `yaml:"sheet_index_column"`
	BackupFolderID   string       `yaml:"backup_folder_id"`
}

func NewConfig(path string) (Config, error) {
	c := Config{}
	if err := config.LoadWithEnv(&c, path); err != nil {
		return c, errors.Wrap(err, "bmaerror load config file")
	}

	return c, nil
}
