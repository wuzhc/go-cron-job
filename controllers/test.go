package controllers

import (
	"github.com/astaxie/beego"
)

type TestController struct {
	BaseController
}

func (c *TestController) Err() {
	beego.Info("xxx", "222", "333")
	beego.Debug("xxx", "222", "333")
	beego.Warning("xxx", "222", "333")
	beego.Warn("xxx", "222", "333")
	c.RspSuccess("success")
}
