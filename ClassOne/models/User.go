package models

type User struct {
	Id       int
	UserName string
	Passwd   string
	Articles []*Article `orm:"rel(m2m)"`
}
