package sheetmerger

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	drive "google.golang.org/api/drive/v3"
	sheets "google.golang.org/api/sheets/v4"
)

func New(config Config) (*SheetMerge, error) {
	cnf := &jwt.Config{
		Email:        config.Client.Email,
		PrivateKey:   []byte(config.Client.PrivateKey),
		PrivateKeyID: config.Client.PrivateKeyID,
		TokenURL:     google.JWTTokenURL,
		Scopes: []string{
			drive.DriveScope,
			drive.DriveMetadataScope,
			sheets.SpreadsheetsScope,
		},
	}
	cli := cnf.Client(context.Background())

	backup, err := NewBackup(cli)
	if err != nil {
		return nil, errors.Wrap(err, "error new backup")
	}

	merge, err := NewMerge(cli)
	if err != nil {
		return nil, errors.Wrap(err, "error new merge")
	}

	return &SheetMerge{
		Backuper: backup,
		Merger:   merge,
		config:   config,
	}, nil
}

type SheetMerge struct {
	Backuper
	Merger
	config Config
}

func (sm *SheetMerge) Backup(baseSheetKey, backupFolderName string) error {
	req := &BackupRequest{
		BaseSheetKey:           baseSheetKey,
		BaseSheetName:          sm.config.BaseSheetName,
		SheetIndexColumn:       sm.config.SheetIndexColumn,
		BackupParentFolderID:   sm.config.BackupFolderID,
		BackupParentFolderName: backupFolderName,
	}

	return sm.Backuper.Backup(req)
}

func (sm *SheetMerge) Merge(baseSheetKey, diffSheetKey string, sheetNames ...string) error {
	req := &MergeRequest{
		BaseSheetKey:     baseSheetKey,
		BaseSheetName:    sm.config.BaseSheetName,
		DiffSheetKey:     diffSheetKey,
		SheetNames:       sheetNames,
		IDColumnName:     "id",
		SheetIndexColumn: sm.config.SheetIndexColumn,
	}

	return sm.Merger.MergeBySheetKey(req)
}
