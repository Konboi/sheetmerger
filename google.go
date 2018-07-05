package sheetmerger

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	drive "google.golang.org/api/drive/v3"
	sheets "google.golang.org/api/sheets/v4"
)

type googleDrive struct {
	srv *drive.Service
}

// NewGoogleDrive is XXX
func NewGoogleDrive(cli *http.Client) (DriveService, error) {
	srv, err := drive.New(cli)
	if err != nil {
		return nil, errors.Wrap(err, "error new googl drive service")
	}

	return &googleDrive{srv}, nil
}

func (d *googleDrive) Create(name, fileType string, parents []string) (File, error) {
	f, err := d.srv.Files.Create(&drive.File{
		Name:     name,
		MimeType: fileType,
		Parents:  parents,
	}).Do()
	if err != nil {
		return nil, errors.Wrap(err, "error create file")
	}

	return &file{
		id:   f.Id,
		name: f.Name,
	}, nil
}

func (d *googleDrive) Copy(fileID string, parents []string) (File, error) {
	base, err := d.Find(fileID)
	if err != nil {
		return nil, errors.Wrapf(err, "error find base file id:%s", fileID)
	}

	f, err := d.srv.Files.Copy(fileID, &drive.File{
		Parents: parents,
		Name:    base.Name(),
	}).Do()
	if err != nil {
		return nil, errors.Wrap(err, "error copy file")
	}

	return &file{
		id:   f.Id,
		name: f.Name,
	}, nil
}

func (d *googleDrive) Find(fileID string) (File, error) {
	f, err := d.srv.Files.Get(fileID).Do()
	if err != nil {
		return nil, errors.Wrap(err, "error get file")
	}

	return &file{
		id:   f.Id,
		name: f.Name,
	}, nil
}

func (d *googleDrive) Delete(fileID string) error {
	if err := d.srv.Files.Delete(fileID).Do(); err != nil {
		return errors.Wrap(err, "error delete file")
	}

	return nil
}

func NewGoogleSpreadSheet(cli *http.Client) (SheetService, error) {
	srv, err := sheets.New(cli)
	if err != nil {
		return nil, errors.Wrap(err, "error new google spread sheet")
	}

	return &googleSpreadSheet{srv}, nil
}

type googleSpreadSheet struct {
	srv *sheets.Service
}

func (s *googleSpreadSheet) Get(key, sheetName string) (*Sheet, error) {
	sheet, err := s.srv.Spreadsheets.Values.Get(key, sheetName).Do()
	if err != nil {
		return nil, errors.Wrapf(err, "error get spread sheet. key:%s", key)
	}

	return newSheet(key, sheetName, sheet)
}

func (s *googleSpreadSheet) Append(base, diff *Sheet) error {
	_range := fmt.Sprintf("%s!A%d", base.Name(), len(base.Rows())+3)

	// to google values format
	diffCRMap := diff.ColumnAndRowMaps()
	add := [][]interface{}{}
	for _, row := range diffCRMap {
		v := make([]interface{}, len(row))
		for j, c := range base.Columns() {
			if c != "" {
				v[j] = row[c.(string)]
			} else {
				v[j] = ""
			}
		}
		add = append(add, v)
	}

	_, err := s.srv.Spreadsheets.Values.Append(
		base.Key(),
		_range,
		&sheets.ValueRange{
			MajorDimension: "ROWS",
			Values:         add,
		},
	).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return errors.Wrapf(err, "error append failed. key:%s range:%s", base.Key(), _range)
	}

	return nil
}

func (s *googleSpreadSheet) Replace(sheet *Sheet, column string, replaceKeyValue map[string]string) error {
	index := sheet.Index(column)
	if index < 0 {
		return errors.Errorf("error replace colums:%s is not exists", column)
	}

	update := [][]interface{}{}
	update = append(update, []interface{}{sheet.Columns()[index]})
	for _, row := range replaceValue(sheet.Rows(), index, replaceKeyValue) {
		update = append(update, []interface{}{row[index]})
	}

	_range := fmt.Sprintf("%s!%s:%s", sheet.Name(), index2rangeColumn(index), index2rangeColumn(index))

	return s.update(sheet.Key(), _range, update)
}

func (s *googleSpreadSheet) update(key, _range string, values [][]interface{}) error {
	_, err := s.srv.Spreadsheets.Values.Update(
		key,
		_range,
		&sheets.ValueRange{
			MajorDimension: "ROWS",
			Values:         values,
		}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return errors.Wrapf(err, "error update sheet key:%s range:%s", key, _range)
	}

	return nil
}

func newSheet(key, name string, gs *sheets.ValueRange) (*Sheet, error) {
	values := gs.Values
	if len(values) < 2 {
		return nil, errors.Errorf("error this sheet is empty. (key:%s, name:%s)", key, name)
	}

	colsIntrf := values[0]
	if len(colsIntrf) == 0 {
		return nil, errors.Errorf("error this sheet column is empty. (key:%s name:%s)", key, name)
	}

	return &Sheet{
		key:    key,
		name:   name,
		values: values,
	}, nil
}

func replaceValue(rows [][]interface{}, index int, kvMap map[string]string) [][]interface{} {
	update := [][]interface{}{}
	for _, row := range rows {
		r := row[index]
		if v, ok := kvMap[r.(string)]; ok {
			row[index] = v
		}
		update = append(update, row)
	}

	return update
}

func rangeFromIndexAndValue(name string, index int, value [][]interface{}) string {
	return fmt.Sprintf("%s!%s:%s", name, index2rangeColumn(index), index2rangeColumn(index))
}

func index2rangeColumn(i int) string {
	i = i + 1 // シートのカラムが１始まりなので
	j := 0
	r := ""
	for {
		i = i - 1
		j = i % 26
		i = i / 26
		if 0 < i {
			r = fmt.Sprintf("%s%s", string('A'+j), r)
		} else {
			break
		}
	}

	return fmt.Sprintf("%s%s", string('A'+j), r)
}
