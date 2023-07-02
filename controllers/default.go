package controllers

import (
	"encoding/json"
	"io/ioutil"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

type TokenRequest struct {
	Receiver string `json:"receiver" valid:"Required;Match(/^0x[0-9a-fA-F]{40}$/)"`
	Amount   string `json:"amount" valid:"Required"`
}

func validateReq(c *MainController, req TokenRequest) {
	valid := validation.Validation{}
	b, validationErr := valid.Valid(&req)
	if validationErr != nil {
		logs.Error(validationErr)
		c.Data["json"] = map[string]interface{}{
			"error": validationErr,
		}
		c.ServeJSON()
	}
	if !b {
		for _, err := range valid.Errors {
			logs.Error(err.Key + err.Message)
			c.Data["json"] = map[string]interface{}{
				"error": err.Key + err.Message,
			}
			c.ServeJSON()
		}
	}
}

func (c *MainController) MintToken() {
	body, err := ioutil.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"error": "request body reading error",
		}
		logs.Error("request body reading error")
		c.ServeJSON()
		return
	}

	var req TokenRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"error": "request body parcing error",
		}
		logs.Error("request body parcing error")
		c.ServeJSON()
		return
	}

	// receiver and amount validation
	validateReq(c, req)

	// send response
	c.Data["json"] = map[string]interface{}{
		"token": "generated_token",
	}
	c.ServeJSON()
}
