package sheetmerger

type file struct {
	id, name string
}

func (f *file) ID() string {
	return f.id
}

func (f *file) Name() string {
	return f.name
}
