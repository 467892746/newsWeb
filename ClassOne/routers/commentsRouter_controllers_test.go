package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["ClassOne/controllers/test:Controller"] = append(beego.GlobalControllerRouter["ClassOne/controllers/test:Controller"],
        beego.ControllerComments{
            Method: "Login",
            Router: `/*`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
