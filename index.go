package otira

type Index struct {
	name       string
	fieldMetas []*FieldMeta
}

func NewIndex(name string, field0, field1 *FieldMeta, fields ...*FieldMeta) *Index {
	index := new(Index)
	index.name = name

	return index

}
