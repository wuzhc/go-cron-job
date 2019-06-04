package routers

import (
	"cron-job/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.AutoRouter(&controllers.JobsController{})
	beego.AutoRouter(&controllers.TestController{})
}
