package controllers

import (
	"encoding/json"
	"io/ioutil"

	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

type TokenRequest struct {
	Receiver string `json:"receiver"`
	Amount   string `json:"amount"`
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.vip"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

func (c *MainController) MintToken() {
	body, err := ioutil.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"error": "request body reading error",
		}
		c.ServeJSON()
		return
	}

	var req TokenRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"error": "request body parcing error",
		}
		c.ServeJSON()
		return
	}

	// receiver and amount validation

	// send response
	c.Data["json"] = map[string]interface{}{
		"token": "generated_token",
	}
	c.ServeJSON()
}
