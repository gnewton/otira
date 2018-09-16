package otira

import (
	"fmt"
)

type FieldMetaImpl struct {
	name       string
	table      *TableMeta
	length     int
	unique     bool
	indexed    bool
	nullable   bool
	primaryKey bool
	fixed      bool
	//Table(*Table)
}

func (b *FieldMetaImpl) SetName(n string) {
	b.name = n
}
func (b *FieldMetaImpl) Name() string {
	return b.name
}

func (b *FieldMetaImpl) SetTable(table *TableMeta) {
	b.table = table
}

func (b *FieldMetaImpl) Table() *TableMeta {
	return b.table
}

func (b *FieldMetaImpl) SetLength(n int) {
	b.length = n
}

func (b *FieldMetaImpl) Length() int {
	return b.length
}

func (b *FieldMetaImpl) Unique() bool {
	return b.unique
}
func (b *FieldMetaImpl) SetUnique(v bool) {
	b.unique = v
}
func (b *FieldMetaImpl) Indexed() bool {
	return b.indexed
}
func (b *FieldMetaImpl) SetIndexed(v bool) {
	b.indexed = v
}

func (b *FieldMetaImpl) Nullable() bool {
	return b.nullable
}
func (b *FieldMetaImpl) SetNullable(v bool) {
	b.nullable = v
}

func (b *FieldMetaImpl) Fixed() bool {
	return b.fixed
}
func (b *FieldMetaImpl) SetFixed(v bool) {
	b.fixed = v
}

func (b *FieldMetaImpl) PrimaryKey() bool {
	return b.primaryKey
}

func (b *FieldMetaImpl) SetPrimaryKey(v bool) {
	b.primaryKey = v
}
func (b *FieldMetaImpl) String() string {
	//return "Name:" + b.name + " PrimaryKey:" + b.primaryKey
	return fmt.Sprintf("Name: [%s]  PrimaryKey: %t  Unique: %t  Indexed: %t  Nullable: %t", b.name, b.primaryKey, b.unique, b.indexed, b.nullable)
}

//////////////////////////////////////////////////////
type FieldMetaString struct {
	FieldMetaImpl
}

type FieldMetaInt struct {
	FieldMetaImpl
}

type FieldMetaFloat struct {
	FieldMetaImpl
}

type FieldMetaByte struct {
	FieldMetaImpl
}

type FieldMetaTime struct {
	FieldMetaImpl
}

type FieldMetaTimeStamp struct {
	FieldMetaImpl
}

type FieldMetaBool struct {
	FieldMetaImpl
}
