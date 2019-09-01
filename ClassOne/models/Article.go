package models

import "time"

type Article struct {
	Id int  `orm:"pk;auto"`
	Title string `orm:"size(20)" `  //文章标题
	Content string `orm:"size(500)"`//文章内容
	Img string  `orm:"siez(50);null"`//图片 路径
	//Type int  //类型
	AddTime time.Time `orm:"type(datetime);auto_now_add"` //文章发布时间
	Count int64 `orm:"default(0)"`//阅读量
	ArticleType *ArticleType `orm:"rel(fk)"`
	Users []*User `orm:"reverse(many)"`
}

type ArticleType struct {
	Id int
	TypeName string `orm:"size(20)"`
	Article []*Article `orm:"reverse(many)"`
}

