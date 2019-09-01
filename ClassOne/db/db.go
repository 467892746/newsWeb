package db

import (
	"ClassOne/models"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var DB orm.Ormer

func init()  {
	orm.Debug = true
	orm.RegisterDataBase("default", "mysql", "root:mysql@tcp(127.0.0.1:3306)/newsWeb?charset=utf8")
	orm.RegisterModel(&models.User{}, &models.Article{}, &models.ArticleType{})
	orm.RunSyncdb("default", false, true)
	DB = orm.NewOrm()
}