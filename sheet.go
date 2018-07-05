package sheetmerger

type Sheet struct {
	key    string
	name   string
	values [][]interface{}
}

func (s *Sheet) Name() string { return s.name }
func (s *Sheet) Key() string  { return s.key }
func (s *Sheet) Columns() []interface{} {
	return s.values[0]
}

func (s *Sheet) Rows() [][]interface{} {
	return s.values[1:]
}

func (s *Sheet) UniqueValuesByColumn(column string) []string {
	index := s.Index(column)
	if index < 0 {
		return nil
	}

	values := []string{}
	valuesMap := map[string]struct{}{}
	for _, row := range s.Rows() {
		if len(row) < 1 { // 空行
			continue
		}

		v := row[index].(string)
		if _, ok := valuesMap[v]; !ok {
			values = append(values, v)
			valuesMap[v] = struct{}{}
		}
	}

	return values
}

func (s *Sheet) ColumnAndRowMaps() []map[string]string {
	values := []map[string]string{}
	for _, row := range s.Rows() {
		value := map[string]string{}
		for i, c := range s.Columns() {
			if c == "" || len(row) <= i {
				value[c.(string)] = ""
			} else {
				value[c.(string)] = row[i].(string)
			}
		}
		values = append(values, value)
	}

	return values
}

func (s *Sheet) DuplicatedColumnValues(column string, diff *Sheet) []interface{} {
	dupCheckMap := map[string]struct{}{}
	for _, k := range s.UniqueValuesByColumn(column) {
		dupCheckMap[k] = struct{}{}
	}

	duplicateValues := []interface{}{}
	for _, k := range diff.UniqueValuesByColumn(column) {
		if _, ok := dupCheckMap[k]; ok {
			duplicateValues = append(duplicateValues, k)
		}
	}

	return duplicateValues
}

func (s *Sheet) Index(column string) int {
	for i, c := range s.Columns() {
		if c == column {
			return i
		}
	}

	return -1
}
