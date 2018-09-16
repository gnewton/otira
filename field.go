package otira

type Field struct {
	field       struct{}
	fieldMeta   *FieldMeta
	hasSetValue bool
}
