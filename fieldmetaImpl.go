package otira

import (
	"fmt"
	"time"
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

func (b *FieldMetaImpl) String() string {
	//return "Name:" + b.name + " PrimaryKey:" + b.primaryKey
	return fmt.Sprintf("Name: [%s]  PrimaryKey: %t  Unique: %t  Indexed: %t  Nullable: %t", b.name, b.primaryKey, b.unique, b.indexed, b.nullable)
}

// func (b *FieldMetaImpl) IsSameType(interface{}) bool {
// 	return false
// }

//////////////////////////////////////////////////////
type FieldMetaString struct {
	FieldMetaImpl
}

func (fm *FieldMetaString) IsSameType(v interface{}) bool {
	_, ok := v.(string)
	return ok
}

type FieldMetaFloat struct {
	FieldMetaImpl
}

func (fm *FieldMetaFloat) IsSameType(v interface{}) bool {
	_, ok := v.(float32)
	if !ok {
		_, ok = v.(float64)
	}
	return ok
}

type FieldMetaInt struct {
	FieldMetaImpl
}

func (fm *FieldMetaInt) IsSameType(v interface{}) bool {

	_, ok := v.(int)
	if ok {
		return true
	}

	_, ok = v.(int32)
	if ok {
		return true
	}

	_, ok = v.(int64)
	if ok {
		return true
	}

	_, ok = v.(int8)
	if ok {
		return true
	}

	_, ok = v.(int16)
	if ok {
		return true
	}

	_, ok = v.(uint)
	if ok {
		return true
	}

	_, ok = v.(uint8)
	if ok {
		return true
	}

	_, ok = v.(uint16)
	if ok {
		return true
	}

	_, ok = v.(uint32)
	if ok {
		return true
	}

	_, ok = v.(uint64)
	if ok {
		return true
	}

	_, ok = v.(uintptr)
	if ok {
		return true
	}
	return false
}

type FieldMetaByte struct {
	FieldMetaImpl
}

func (fm *FieldMetaByte) IsSameType(v interface{}) bool {
	_, ok := v.([]byte)
	return ok
}

type FieldMetaTime struct {
	FieldMetaImpl
}

func (fm *FieldMetaTime) IsSameType(v interface{}) bool {
	_, ok := v.(time.Time)
	return ok
}

//type FieldMetaTimeStamp struct {
//	FieldMetaImpl
//}

type FieldMetaBool struct {
	FieldMetaImpl
}

func (fm *FieldMetaBool) IsSameType(v interface{}) bool {
	_, ok := v.(bool)
	return ok
}
