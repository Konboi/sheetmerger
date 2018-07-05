package sheetmerger

// File is XXX
type File interface {
	ID() string
	Name() string
}

type DriveService interface {
	Create(name, fileType string, parents []string) (File, error)
	Copy(fileID string, parents []string) (File, error)
	Find(fileID string) (File, error)
	Delete(id string) error
}

type SheetService interface {
	Get(key, sheetName string) (*Sheet, error)
	Replace(sheet *Sheet, column string, replaceKeyValue map[string]string) error
	Append(base, diff *Sheet) error
}
