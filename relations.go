package otira

type Relation struct {
	leftTable              *TableMeta
	rightTable             *TableMeta
	rightTableUniqueFields []FieldMeta
}

type OneToMany struct {
	Relation
	leftKeyField  FieldMeta
	rightKeyField string
}

type ManyToMany struct {
	Relation
	leftKeyField, rightKeyField string
}
