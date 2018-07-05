package sheetmerger

import (
	"log"
	"net/http"

	"github.com/pkg/errors"
)

// Merger is XXX
type Merger interface {
	MergeBySheetKey(req *MergeRequest) error
}

func NewMerge(cli *http.Client) (Merger, error) {
	sheeter, err := NewGoogleSpreadSheet(cli)
	if err != nil {
		return nil, errors.Wrap(err, "error new google spreadsheet")
	}

	return &merger{
		sheet: sheeter,
	}, nil
}

type MergeRequest struct {
	BaseSheetKey     string
	BaseSheetName    string
	DiffSheetKey     string
	SheetIndexColumn string
	IDColumnName     string
	SheetNames       []string
}

type merger struct {
	sheet SheetService
}

func (m *merger) MergeBySheetKey(req *MergeRequest) error {
	baseIndexSheet, err := m.sheet.Get(req.BaseSheetKey, req.BaseSheetName)
	if err != nil {
		return errors.Wrap(err, "error get base index sheet")
	}

	diffIndexSheet, err := m.sheet.Get(req.DiffSheetKey, req.BaseSheetName)
	if err != nil {
		return errors.Wrap(err, "error get diff index sheet")
	}

	for _, name := range req.SheetNames {
		if err := m.merge(baseIndexSheet, diffIndexSheet, name, req.SheetIndexColumn, req.IDColumnName); err != nil {
			return errors.Wrapf(err, "error merge failed base:%s diff:%s sheet:%s",
				baseIndexSheet.Key(),
				diffIndexSheet.Key(),
				name,
			)
		}
		log.Printf("done merge sheet:%s", name)
	}

	return nil
}

func (m *merger) merge(baseSheet, diffSheet *Sheet, name, indexColumn, idColumn string) error {
	baseSheetIndex := baseSheet.ColumnAndRowMaps()

	baseIndexRows, err := indexRows(name, baseSheetIndex)
	if err != nil {
		return errors.Errorf("error base sheet %s can not merge", name)
	}

	diffSheetIndex := diffSheet.ColumnAndRowMaps()
	diffIndexRows, err := indexRows(name, diffSheetIndex)
	if err != nil {
		return errors.Errorf("error diff sheet %s can not merge", name)
	}

	b, err := m.sheet.Get(baseIndexRows[0][indexColumn], name)
	if err != nil {
		return errors.Errorf("error get base %s sheet key:%s", name, baseSheetIndex[0][indexColumn])
	}

	d, err := m.sheet.Get(diffIndexRows[0][indexColumn], name)
	if err != nil {
		return errors.Errorf("error get diff %s sheet key:%s", name, diffSheetIndex[0][indexColumn])
	}

	if 0 < len(b.DuplicatedColumnValues(idColumn, d)) {
		return errors.Errorf("error sheet:%s id duplicated %q", name, b.DuplicatedColumnValues(idColumn, d))
	}

	return m.sheet.Append(b, d)
}

func indexRows(name string, maps []map[string]string) ([]map[string]string, error) {
	indexRows := []map[string]string{}
	for _, m := range maps {
		if m["sheet"] == name {
			indexRows = append(indexRows, m)
		}
	}

	if 1 < len(indexRows) {
		return nil, errors.New("error exists multi index rows")
	}
	if len(indexRows) == 0 {
		return nil, errors.New("error not exists index rows")
	}

	return indexRows, nil
}
