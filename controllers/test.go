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

func (c *TestController) Json() {
	type Udata struct {
		Id         int    `json:"id"`
		Name       string `json:"-"`
		CreateTime string `json:"create_time"`
	}

	var list = make([]Udata, 0)
	data := Udata{1, "wuzhc", "20190604"}
	list = append(list, data)
	c.Data["json"] = list
	c.ServeJSON()
}
