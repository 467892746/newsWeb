package routers

import (
	"ClassOne/controllers/test"
	"github.com/astaxie/beego"
)

func init() {
	//beego.Router("/", &test.Controller{})
	beego.Router("/register", &test.Controller{}, "get:ShowReg;post:HandleReg")
	beego.Router("/", &test.Controller{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/ShowArticle", &test.Controller{}, "get:ShowArticleList;post:HandleSelect")
	beego.Router("/AddArticle", &test.Controller{}, "get:ShowAddArticle;post:HandleAddArticle")
	beego.Router("/ArticleContent", &test.Controller{}, "get:ArticleContent")
	beego.Router("/DeleteArticle", &test.Controller{}, "get:HandleDelete")
	beego.Router("/UpdateArticle", &test.Controller{}, "get:ShowUpdate;post:HandleUpdate")
	beego.Router("/Logout", &test.Controller{}, "get:Logout")


}
