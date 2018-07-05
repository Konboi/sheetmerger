package sheetmerger

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

const (
	GoogleDriveFileTypeFolder = "application/vnd.google-apps.folder"
)

type BackupRequest struct {
	BaseSheetKey           string
	BaseSheetName          string
	SheetIndexColumn       string
	BackupParentFolderID   string
	BackupParentFolderName string
}

func (r *BackupRequest) IsValid() error {
	if r.BaseSheetKey == "" {
		return fmt.Errorf("error base sheet key is required.")
	}
	if r.BackupParentFolderID == "" {
		return fmt.Errorf("error backup parent folder id is required.")
	}

	return nil
}

// Backuper is XXX
type Backuper interface {
	Backup(req *BackupRequest) error
}

type backup struct {
	drive DriveService
	sheet SheetService
}

func NewBackup(cli *http.Client) (Backuper, error) {
	drive, err := NewGoogleDrive(cli)
	if err != nil {
		return nil, errors.Wrap(err, "error new drive")
	}

	sheet, err := NewGoogleSpreadSheet(cli)
	if err != nil {
		return nil, errors.Wrap(err, "error new spreadsheet")
	}

	return &backup{
		drive: drive,
		sheet: sheet,
	}, nil
}

func (b *backup) Backup(req *BackupRequest) error {
	if err := req.IsValid(); err != nil {
		return errors.Wrap(err, "error request is invalid")
	}

	log.Println("start backup")

	rootFolder, err := b.drive.Find(req.BackupParentFolderID)
	if err != nil {
		return errors.Wrap(err, "error find root folder")
	}

	backupFolder, err := b.drive.Create(req.BackupParentFolderName, GoogleDriveFileTypeFolder, []string{
		rootFolder.ID(),
	})
	if err != nil {
		return errors.Wrap(err, "error create backup folder")
	}

	baseSheetInfo, err := b.drive.Find(req.BaseSheetKey)
	if err != nil {
		return errors.Wrap(err, "error find base sheet")
	}

	baseSheet, err := b.sheet.Get(req.BaseSheetKey, req.BaseSheetName)
	if err != nil {
		return errors.Wrap(err, "error find base spreadsheet")
	}

	copiedSheetFile, err := b.drive.Copy(baseSheetInfo.ID(), []string{
		backupFolder.ID(),
	})
	if err != nil {
		return errors.Wrap(err, "error copy base sheet file")
	}

	copiedSheet, err := b.sheet.Get(copiedSheetFile.ID(), req.BaseSheetName)
	if err != nil {
		return errors.Wrap(err, "error get copy file")
	}

	replaceIDMap := make(map[string]string)
	otherSheetKeys := baseSheet.UniqueValuesByColumn("key")
	for _, key := range otherSheetKeys {
		copied, err := b.drive.Copy(key, []string{
			backupFolder.ID(),
		})
		if err != nil {
			return errors.Wrapf(err, "error copy file key:%s", key)
		}

		replaceIDMap[key] = copied.ID()
		log.Printf("copy %s:%s done", copied.Name(), copied.ID())
	}

	if err := b.sheet.Replace(copiedSheet, req.SheetIndexColumn, replaceIDMap); err != nil {
		return errors.Wrap(err, "error replace")
	}

	if 0 < len(baseSheet.DuplicatedColumnValues(req.SheetIndexColumn, copiedSheet)) {
		if err := b.drive.Delete(backupFolder.ID()); err != nil {
			return errors.Wrapf(err, "error delete folder (id:%s)", backupFolder.ID())
		}
		return errors.Errorf("error copy sheet:%s is failed key duplicated %.\nbase keys:[%q]\ncopied keys:[%q]", copiedSheetFile.ID(),
			baseSheet.DuplicatedColumnValues(req.SheetIndexColumn, copiedSheet),
			baseSheet.UniqueValuesByColumn("key"),
			copiedSheet.UniqueValuesByColumn("key"),
		)
	}

	log.Println("finish backup")

	return nil
}
