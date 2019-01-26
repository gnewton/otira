package otira

type Relation interface {
	Name() string
}

type baseRelation struct {
	name                   string
	LeftTable              *TableDef
	RightTable             *TableDef
	RightTableUniqueFields []FieldDef // fields to find out if a record exists; these fields are used in a lookup
	LeftKeyField           FieldDef
	RightKeyField          FieldDef
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
	JoinTable *TableDef
	pkCache   map[string]struct{}
}

func (mtm *ManyToMany) String() string {
	return "ManyToMany"
}
