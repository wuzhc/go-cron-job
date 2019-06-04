package main

import (
	_ "cron-job/hub"
	_ "cron-job/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
