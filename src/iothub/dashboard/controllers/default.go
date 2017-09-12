package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "github.com/cloustone/sentel"
	c.Data["Email"] = "jenson.zuo@qq.com"
	c.TplName = "index.tpl"
}
