package routers

import (
	"mint-token-app/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/mint_token", &controllers.MainController{}, "post:MintToken")
	beego.Router("/", &controllers.MainController{})
}
