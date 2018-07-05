package sheetmerger

import (
	"reflect"
	"testing"
)

func TestSheet_replace(t *testing.T) {
	sheet := &Sheet{
		values: [][]interface{}{
			{"id", "name", "key"},
			{"1", "foo", "key1"},
			{"2", "foo", "key2"},
			{"3", "hoge", "key1"},
			{"4", "hoge", "key3"},
		},
	}

	type input struct {
		column  string
		replace map[string]string
	}

	tests := []struct {
		input  input
		output [][]interface{}
	}{
		{
			input: input{
				column: "id",
				replace: map[string]string{
					"1": "2",
					"2": "3",
				},
			},
			output: [][]interface{}{
				{"2", "foo", "key1"},
				{"3", "foo", "key2"},
				{"3", "hoge", "key1"},
				{"4", "hoge", "key3"},
			},
		},
		{
			input: input{
				column: "name",
				replace: map[string]string{
					"foo": "fuga",
				},
			},
			output: [][]interface{}{
				{"2", "fuga", "key1"},
				{"3", "fuga", "key2"},
				{"3", "hoge", "key1"},
				{"4", "hoge", "key3"},
			},
		},
	}

	for _, test := range tests {
		if !reflect.DeepEqual(replaceValue(
			sheet.Rows(),
			sheet.Index(test.input.column),
			test.input.replace,
		), test.output) {
			t.Errorf("error invalid replace rows. got:%v want:%v", sheet.Rows(), test.output)
		}
	}
}
