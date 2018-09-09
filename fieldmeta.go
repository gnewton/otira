package otira

type FieldMeta interface {
	Name(string)
	GetName() string
	Length(int)
	GetLength() int
	Unique(bool)
	IsUnique() bool
	Indexed(bool)
	IsIndexed() bool
	NotNullable(bool)
	IsNotNullable() bool
	PrimaryKey(bool)
	IsPrimaryKey() bool
	String() string
	//Table(*Table)
}

type FieldMetaImpl struct {
	name        string
	length      int
	unique      bool
	indexed     bool
	notNullable bool
	primaryKey  bool
	//Table(*Table)
}

func (b *FieldMetaImpl) Name(n string) {
	b.name = n
}
func (b *FieldMetaImpl) GetName() string {
	return b.name
}
func (b *FieldMetaImpl) Length(n int) {
	b.length = n
}

func (b *FieldMetaImpl) GetLength() int {
	return b.length
}

func (b *FieldMetaImpl) IsUnique() bool {
	return b.unique
}
func (b *FieldMetaImpl) Unique(v bool) {
	b.unique = v
}
func (b *FieldMetaImpl) IsIndexed() bool {
	return b.indexed
}

func (b *FieldMetaImpl) Indexed(v bool) {
	b.indexed = v
}
func (b *FieldMetaImpl) IsNotNullable() bool {
	return b.notNullable
}

func (b *FieldMetaImpl) NotNullable(v bool) {
	b.notNullable = v
}
func (b *FieldMetaImpl) IsPrimaryKey() bool {
	return b.primaryKey
}

func (b *FieldMetaImpl) PrimaryKey(v bool) {
	b.primaryKey = v
}
func (b *FieldMetaImpl) String() string {
	return "fff"
}
