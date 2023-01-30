package models

type ForwardDeclPoint struct {
	Type string
	Database string
	Table string
}

type ForwardDecl struct {
	Name string
	From ForwardDeclPoint
	To ForwardDeclPoint
	Watch bool
}
