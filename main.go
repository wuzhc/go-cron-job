package main

import (
	_ "cron-job/crontab"
	_ "cron-job/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
