package sheetmerger

type mockFile struct {
	id, name string
}

func (f *mockFile) ID() string {
	return f.id
}

func (f *mockFile) Name() string {
	return f.name
}
