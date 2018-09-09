package otira

type TableMeta struct {
	name       string
	fields     []FieldMeta
	fieldsMap  map[string]FieldMeta
	inited     bool
	oneToMany  []*OneToMany
	manyToMany []*ManyToMany
}

func NewTableMeta(name string) (*TableMeta, error) {
	t := new(TableMeta)
	t.name = name
	return t, nil
}

func (t *TableMeta) Add(f FieldMeta) {
	if t.fields == nil {
		t.fields = make([]FieldMeta, 0)
	}
}
