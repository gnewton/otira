package otira

type Relation interface {
	Name() string
}

type baseRelation struct {
	name                   string
	leftTable              *TableMeta
	rightTable             *TableMeta
	rightTableUniqueFields []FieldMeta // fields to find out if a record exists; these fields are used in a lookup
	leftKeyField           FieldMeta
	rightKeyField          FieldMeta
	cache                  *joinCache
}

func (rel *baseRelation) Name() string {
	return rel.name
}

type OneToMany struct {
	baseRelation
}

func (otm *OneToMany) String() string {
	return "OneToMany:" + otm.name
}

type ManyToMany struct {
	baseRelation
	joinTable *TableMeta
	pkCache   map[string]struct{}
}

func (mtm *ManyToMany) String() string {
	return "ManyToMany"
}
