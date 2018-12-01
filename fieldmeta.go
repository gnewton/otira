package otira

type FieldMeta interface {
	Name() string
	SetName(string)
	Table() *TableMeta
	SetTable(*TableMeta)
	SetLength(int)
	Length() int
	Unique() bool
	SetUnique(bool)
	Indexed() bool
	SetIndexed(bool)
	Nullable() bool
	SetNullable(bool)
	//PrimaryKey() bool
	//SetPrimaryKey(bool)
	String() string
	SetFixed(bool)
	Fixed() bool
	//Table(*Table)
	IsSameType(interface{}) bool
}
