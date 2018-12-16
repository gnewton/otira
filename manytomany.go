package otira

func NewManyToMany() *ManyToMany {
	m2m := new(ManyToMany)
	m2m.cache = NewJoinCache()
	return m2m
}
