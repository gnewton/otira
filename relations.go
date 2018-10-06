package otira

type Relation interface {
	Name() string
}

type BaseRelation struct {
	name       string
	leftTable  *TableMeta
	rightTable *TableMeta
	//rightTableUniqueFields []FieldMeta
	leftKeyField  FieldMeta
	rightKeyField FieldMeta
}

func (rel *BaseRelation) Name() string {
	return rel.name
}

type OneToMany struct {
	BaseRelation
}

func (otm *OneToMany) String() string {
	return "OneToMany:" + otm.name
}

type ManyToMany struct {
	BaseRelation
	joinTable *TableMeta
}

func (mtm *ManyToMany) String() string {
	return "ManyToMany"
}
