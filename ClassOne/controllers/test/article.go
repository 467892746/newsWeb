package test

import (
	"ClassOne/db"
	"ClassOne/models"
	"bytes"
	"encoding/gob"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"math"
	"path"
	"strconv"
	"time"
)

func (c *Controller) HandleSelect() {
	typeName := c.GetString("select")
	if typeName == "" {
		c.Ctx.WriteString("系统错误")
		return
	}
	var articles []models.Article
	db.DB.QueryTable(&models.ArticleType{}).RelatedSel("ArticleType").
		Filter("ArticleType_TypeName", typeName).All(&articles)

}

func (c *Controller) ShowArticleList() {
	userName := c.GetSession("userName")
	if userName == nil {
		c.Redirect("/", 302)
		return
	}
	qs := db.DB.QueryTable(&models.Article{})
	var articles []models.Article



	count, err := qs.Count()
	if err != nil {
		c.Ctx.WriteString("系统错误")
		return
	}

	pageSize := 1

	pageIndex := 1
	start := pageSize*(pageIndex - 1)
	qs = qs.RelatedSel("ArticleType").Limit(pageSize, start)

	pageCount := float64(count)/float64(pageSize)
	pageCount1 := math.Ceil(pageCount)

	qs.All(&articles)
	c.Data["userName"] = userName

	var types []models.ArticleType

	conn, err := redis.Dial("tcp", ":6379")
	rel, err := redis.Bytes(conn.Do("get", "types"))
	if err != nil {
		c.Ctx.WriteString("获取redis数据错误")
		return
	}
	dec := gob.NewDecoder(bytes.NewReader(rel))
	dec.Decode(&types)
	if len(types) == 0 {
		db.DB.QueryTable(&models.ArticleType{}).All(&types)
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		err = enc.Encode(types)
		_, err = conn.Do("set", "types", buffer.Bytes())

		if err != nil {
			c.Ctx.WriteString("redis数据库操作错误")
			return
		}
	}






	c.Data["types"] = types
	c.Data["count"] = count
	c.Data["pageCount"] = pageCount1
	c.Data["articles"] = articles
	c.Layout = "layout.html"
	c.TplName = "index.html"
}

func (c *Controller) ShowAddArticle() {
	var types []models.ArticleType
	db.DB.QueryTable(&models.ArticleType{}).All(&types)
	c.Data["types"] = types
	c.TplName = "add.html"
}

func (c *Controller) HandleAddArticle() {
	artileName := c.GetString("articleName")
	artiContent := c.GetString("content")
	f, h, err := c.GetFile("uploadname")
	defer f.Close()
	ext := path.Ext(h.Filename)
	beego.Info(ext)
	switch ext {
	case ".jpg", ".png", ".jpeg":
	default:
		c.Ctx.WriteString("上传文件格式不正确")
		return
	}
	if h.Size > 5000000 {
		 c.Ctx.WriteString("文件太大， 不允许上传")
	}
	fileName := time.Now().Format("2006-01-02 15:04:05")


	c.SaveToFile("uploadname", "./static/img/" + fileName+ext)
	if err != nil {
		beego.Info("上传文件失败")
		return
	}
	article := models.Article{
		Title:artileName,
		Content:artiContent,
		Img:"./static/img/" + fileName+ext,
	}

	typeName := c.GetString("select")
	if typeName == "" {
		c.Ctx.WriteString("系统错误")
		return
	}
	var articleType models.ArticleType
	articleType.TypeName = typeName
	err = db.DB.Read(&articleType, "TypeName")
	if err != nil {
		c.Ctx.WriteString("系统错误")
		return
	}

	article.ArticleType = &articleType

	_, err = db.DB.Insert(&article)
	if err != nil {
		beego.Info("插入数据失败")
		return
	}

	c.Redirect("/ShowArticle", 302)

	beego.Info(artileName, artiContent)
	c.Ctx.WriteString("上传文件成功")
}

func (c *Controller) ArticleContent() {
	id, _ := c.GetInt("id")
	article := models.Article{
		Id:id,
	}
	err := db.DB.Read(&article)
	if err != nil {
		c.Ctx.WriteString("查询数据为空")
		return
	}
	article.Count += 1
	m2m := db.DB.QueryM2M(&article, "Users")
	userName := c.GetSession("userName")
	user := models.User{UserName:userName.(string),}

	db.DB.Read(&user, "UserName")
	_, err = m2m.Add(&user)
	if err != nil {
		c.Ctx.WriteString("系统错误")
		return
	}


	db.DB.Update(&article)
	var users models.User
	db.DB.QueryTable("User").Filter("Articles_Article_Id", id).Distinct().All(&users)
	c.Data["article"] = article
	c.Layout = "layout.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["contentHead"] = "head.html"
	c.TplName = "content.html"
}

func (c *Controller) HandleDelete() {
	id, _ := c.GetInt("id")
	article := models.Article{Id:id}
	db.DB.Delete(&article)
	c.Redirect("/ShowArticle", 302)
}

func (c *Controller) ShowUpdate() {
	id, _ := c.GetInt("id")
	beego.Info(id)
	if id == 0 {
		c.Ctx.WriteString("数据错误")
		return
	}
	article := models.Article{
		Id:id,
	}
	err := db.DB.Read(&article)
	if err != nil {
		c.Ctx.WriteString("数据错误")
		return
	}
	c.Data["article"] = article
	c.TplName = "update.html"
}

func (c *Controller) HandleUpdate() {
	name := c.GetString("articleName")
	content := c.GetString("content")
	if name == "" || content == "" {
		 c.Ctx.WriteString("更新数据失败")
		return
	}
	f, h, err := c.GetFile("uploadName")
	if err != nil {
		c.Ctx.WriteString("上传文件失败")
		return
	}
	defer f.Close()
	if h.Size > 500000 {
		c.Ctx.WriteString("文件太大")
		return
	}
	ext := path.Ext(h.Filename)
	switch ext {
	case ".jpg", ".png", ".jpeg":
	default:
		c.Ctx.WriteString("上传文件格式不正确")
		return
	}
	fileName := time.Now().UnixNano()
	c.SaveToFile("uoloadName", "./static/img/" + strconv.Itoa(int(fileName)) + ext)
	id, _ := c.GetInt("id")
	article := models.Article{Id:id}
	err = db.DB.Read(&article)
	if err != nil {
		c.Ctx.WriteString("数据错误")
		return
	}
	article.Title = name
	article.Content = content
	article.Img = "./static/img/" + strconv.Itoa(int(fileName)) + ext
	_, err = db.DB.Update(&article)
	if err != nil {
		c.Ctx.WriteString("更新失败")
		return
	}
}