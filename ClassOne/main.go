package main

import (
	_ "ClassOne/routers"
	"github.com/astaxie/beego"
	_ "ClassOne/db"
	)



func main() {
	beego.Run()
}

