package sheetmerger

import (
	"fmt"
	"reflect"
	"testing"
)

func fileID(k, name string) string {
	return fmt.Sprintf("%s:%s", k, name)
}

func Test_index(t *testing.T) {
	sheet := &Sheet{
		values: [][]interface{}{
			{"id", "name", "key"},
		},
	}
	tests := []struct {
		input  string
		output int
	}{
		{
			input:  "id",
			output: 0,
		},
		{
			input:  "key",
			output: 2,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("case:%d", i), func(t *testing.T) {
			if sheet.Index(test.input) != test.output {
				t.Errorf("error invalid index got:%d want:%d", sheet.Index(test.input), test.output)
			}
		})
	}
}

func TestColumnAndRowMaps(t *testing.T) {

	tests := []struct {
		input  *Sheet
		output []map[string]string
	}{
		{
			input: &Sheet{
				values: [][]interface{}{
					{"id", "", "name", "key"},
					{"1", "", "foo", "key1"},
					{"2", "", "foo", "key2"},
				},
			},
			output: []map[string]string{
				{"id": "1", "": "", "name": "foo", "key": "key1"},
				{"id": "2", "": "", "name": "foo", "key": "key2"},
			},
		},
		{
			input: &Sheet{
				values: [][]interface{}{
					{"id", "", "name", "", "key"},
					{"3", "", "hoge", "", "key3"},
					{"4", "", "fuga", "", "key4"},
				},
			},
			output: []map[string]string{
				{"id": "3", "": "", "name": "hoge", "key": "key3"},
				{"id": "4", "": "", "name": "fuga", "key": "key4"},
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("case:%d", i), func(t *testing.T) {
			if !reflect.DeepEqual(test.input.ColumnAndRowMaps(), test.output) {
				t.Errorf("error invalid result. got:%v want:%v",
					test.input.ColumnAndRowMaps(),
					test.output,
				)
			}
		})
	}
}

func TestSheeterUniqueValuesByColumns(t *testing.T) {
	type input struct {
		column string
		sheet  *Sheet
	}
	tests := []struct {
		input  input
		output []string
	}{
		{
			input: input{
				column: "key",
				sheet: &Sheet{
					values: [][]interface{}{
						{"id", "name", "key"},
						{"1", "foo", "key1"},
						{"2", "foo", "key2"},
						{"3", "hoge", "key1"},
						{"4", "hoge", "key3"},
					},
				},
			},
			output: []string{"key1", "key2", "key3"},
		},
		{
			input: input{
				column: "name",
				sheet: &Sheet{
					values: [][]interface{}{
						{"id", "name", "key"},
						{"1", "foo", "key1"},
						{"2", "foo", "key2"},
						{"3", "hoge", "key1"},
						{"4", "hoge", "key3"},
					},
				},
			},
			output: []string{"foo", "hoge"},
		},
		{
			input: input{
				column: "fuga",
				sheet: &Sheet{
					values: [][]interface{}{
						{"id", "name", "key"},
						{"1", "foo", "key1"},
						{"2", "foo", "key2"},
					},
				},
			},
			output: nil,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("case:%d", i), func(t *testing.T) {
			values := test.input.sheet.UniqueValuesByColumn(test.input.column)
			if !reflect.DeepEqual(test.output, values) {
				t.Errorf("error invalid result. got:%q want:%q", values, test.output)
			}
		})
	}
}

func TestSheetDuplicateValues(t *testing.T) {
	type input struct {
		column   string
		from, to *Sheet
	}
	tests := []struct {
		input  input
		output []interface{}
	}{
		{
			input: input{
				column: "key",
				from: &Sheet{
					values: [][]interface{}{
						{"id", "key"},
						{"1", "key1"},
						{"2", "key2"},
						{"3", "key1"},
					},
				},
				to: &Sheet{
					values: [][]interface{}{
						{"id", "key"},
						{"1", "key3"},
						{"2", "key4"},
						{"3", "key3"},
					},
				},
			},
			output: []interface{}{},
		},
		{
			input: input{
				column: "id",
				from: &Sheet{
					values: [][]interface{}{
						{"id", "key"},
						{"1", "key1"},
						{"2", "key2"},
						{"3", "key1"},
					},
				},
				to: &Sheet{
					values: [][]interface{}{
						{"id", "key"},
						{"1", "key3"},
						{"2", "key4"},
						{"3", "key3"},
					},
				},
			},
			output: []interface{}{"1", "2", "3"},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("case:%d", i), func(t *testing.T) {
			duplicatedValues := test.input.from.DuplicatedColumnValues(test.input.column, test.input.to)
			if !reflect.DeepEqual(duplicatedValues, test.output) {
				t.Error("error invalid result")
			}
		})
	}
}
