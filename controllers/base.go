package controllers

import (
	"strings"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
	actionName     string
	controllerName string
}

type NextPreparer interface {
	NextPrepare()
}

// controler执行之前操作
func (c *BaseController) Prepare() {
	cname, aname := c.GetControllerAndAction()
	c.controllerName = strings.ToLower(cname[:len(cname)-10])
	c.actionName = strings.ToLower(aname)
}

// 显示模板
func (c *BaseController) Display(tpl ...string) {
	if len(tpl) > 0 {
		c.TplName = tpl[0]
	} else {
		c.Layout = "public/layout.html"
		c.TplName = c.controllerName + "/" + c.actionName + ".html"
	}
}

// 响应数据
func (c *BaseController) RspData(data map[string]interface{}) {
	c.rsp(data, 0, "success")
}

// 响应参数错误
func (c *BaseController) RspParamError(msg string) {
	c.rsp(nil, 21, msg)
}

// 响应成功信息
func (c *BaseController) RspSuccess(msg string) {
	c.rsp(nil, 0, msg)
}

// 响应失败信息
func (c *BaseController) RspFail(msg error) {
	c.rsp(nil, 1, msg.Error())
}

// 响应内容
func (c *BaseController) rsp(data map[string]interface{}, code int, msg string) {
	out := make(map[string]interface{})
	out["code"] = code
	out["msg"] = msg
	out["data"] = data
	c.Data["json"] = out
	c.ServeJSON()
	c.StopRun()
}
