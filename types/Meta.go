package types

type MetaCategory struct {
	Key    string
	Counts uint
}

type Meta struct {
	AllCount   uint
	Categories []MetaCategory
}
