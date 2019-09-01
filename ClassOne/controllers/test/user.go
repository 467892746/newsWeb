package test

import (
	"ClassOne/db"
	"ClassOne/models"
	"github.com/astaxie/beego"
	"time"
)

func (c *Controller) ShowReg() {
	c.TplName = "register.html"
}

func (c *Controller) HandleReg() {
	name := c.GetString("userName")
	passwd := c.GetString("password")
	if name == "" || passwd == ""{
		beego.Info("用户名或者密码不能为空")
		c.TplName = "register.html"
		return
	}
	user := models.User{
		UserName:name,
		Passwd:passwd,
	}
	_, err := db.DB.Insert(&user)
	if err != nil {
		beego.Info("插入数据失败")
	}
	//c.TplName = "login.html"
	c.Redirect("/?userName=biwentao123&password=biwentao123", 302)
}

func (c *Controller) ShowLogin() {
	name := c.Ctx.GetCookie("userName")
	if name != "" {
		c.Data["userName"] = name
		c.Data["check"] = "checked"
	}
	c.TplName = "login.html"
}

func (c *Controller) HandleLogin() {
	name := c.GetString("userName")
	passWD := c.GetString("password")
	beego.Info(name, passWD)
	if name == "" || passWD == "" {
		beego.Info("用户名或密码不能为空")
		c.TplName = "login.html"
		return
	}
	user := models.User{
		UserName:name,
	}
	err := db.DB.Read(&user, "UserName")
	if err != nil {
		beego.Info("用户名失败")
		c.TplName = "login.html"
		return
	}
	if user.Passwd != passWD {
		beego.Info("密码失败")
		c.TplName = "login.html"
		return
	}
	check := c.GetString("remember")
	if check == "on" {
		c.Ctx.SetCookie("userName", name, time.Second * 3600)
	}else {
		c.Ctx.SetCookie("userName", "sss", -1)
	}
	c.SetSession("userName", name)
	c.Redirect("/ShowArticle", 302)
}

func (c *Controller) Logout() {
	c.DelSession("userName")
	c.Redirect("/", 302)
}