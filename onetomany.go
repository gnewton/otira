package otira

func NewOneToMany() *OneToMany {
	one2m := new(OneToMany)
	one2m.cache = NewJoinCache()
	return one2m
}
