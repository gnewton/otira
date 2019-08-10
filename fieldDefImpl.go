package otira

import (
	"fmt"
	"time"
)

type FieldDefImpl struct {
	fixed      bool
	indexed    bool
	length     int
	name       string
	nullable   bool
	primaryKey bool
	table      *TableDef
	unique     bool
	//Table(*Table)
}

func (b *FieldDefImpl) SetName(n string) {
	b.name = n
}
func (b *FieldDefImpl) Name() string {
	return b.name
}

func (b *FieldDefImpl) SetTable(table *TableDef) {
	b.table = table
}

func (b *FieldDefImpl) Table() *TableDef {
	return b.table
}

func (b *FieldDefImpl) SetLength(n int) {
	b.length = n
}

func (b *FieldDefImpl) Length() int {
	return b.length
}

func (b *FieldDefImpl) Unique() bool {
	return b.unique
}
func (b *FieldDefImpl) SetUnique(v bool) {
	b.unique = v
}
func (b *FieldDefImpl) Indexed() bool {
	return b.indexed
}
func (b *FieldDefImpl) SetIndexed(v bool) {
	b.indexed = v
}

func (b *FieldDefImpl) Nullable() bool {
	return b.nullable
}
func (b *FieldDefImpl) SetNullable(v bool) {
	b.nullable = v
}

func (b *FieldDefImpl) Fixed() bool {
	return b.fixed
}
func (b *FieldDefImpl) SetFixed(v bool) {
	b.fixed = v
}

func (b *FieldDefImpl) String() string {
	//return "Name:" + b.name + " PrimaryKey:" + b.primaryKey
	return fmt.Sprintf("Name: [%s]  PrimaryKey: %t  Unique: %t  Indexed: %t  Nullable: %t", b.name, b.primaryKey, b.unique, b.indexed, b.nullable)
}

// func (b *FieldDefImpl) IsSameType(interface{}) bool {
// 	return false
// }

//////////////////////////////////////////////////////
type FieldDefString struct {
	FieldDefImpl
}

func (fm *FieldDefString) IsSameType(v interface{}) bool {
	_, ok := v.(string)
	return ok
}

type FieldDefFloat struct {
	FieldDefImpl
}

func (fm *FieldDefFloat) IsSameType(v interface{}) bool {
	_, ok := v.(float32)
	if !ok {
		_, ok = v.(float64)
	}
	return ok
}

type FieldDefInt struct {
	FieldDefImpl
}

type FieldDefInt64 struct {
	FieldDefImpl
}

func (fm *FieldDefInt64) IsSameType(v interface{}) bool {
	_, ok := v.(int64)
	if ok {
		return true
	}
	return false
}

func (fm *FieldDefInt) IsSameType(v interface{}) bool {

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

	_, ok = v.(int64)
	if ok {
		return true
	}

	_, ok = v.(uintptr)
	if ok {
		return true
	}
	return false
}

type FieldDefByte struct {
	FieldDefImpl
}

func (fm *FieldDefByte) IsSameType(v interface{}) bool {
	_, ok := v.([]byte)
	return ok
}

type FieldDefTime struct {
	FieldDefImpl
}

func (fm *FieldDefTime) IsSameType(v interface{}) bool {
	_, ok := v.(time.Time)
	return ok
}

//type FieldDefTimeStamp struct {
//	FieldDefImpl
//}

type FieldDefBool struct {
	FieldDefImpl
}

func (fm *FieldDefBool) IsSameType(v interface{}) bool {
	_, ok := v.(bool)
	return ok
}
