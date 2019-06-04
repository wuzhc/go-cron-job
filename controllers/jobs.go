package controllers

import (
	"cron-job/crontab"
	"fmt"
	// "os/exec"
)

type JobsController struct {
	BaseController
}

// 添加定时任务
func (c *JobsController) Add() {
	spec := c.GetString("spec", "*/5 * * * * ?")
	cmd := c.GetString("cmd", "echo \"xxxxx\" >> text.txt")
	timeout, _ := c.GetInt("timeout", 30)
	err := crontab.AddJob(spec, cmd, timeout)
	if err != nil {
		fmt.Println(err)
	}

	c.RspSuccess("添加任务成功")
}

// 暂停定时任务
func (c *JobsController) Pause() {
	id, _ := c.GetInt("id")
	if id <= 0 {
		c.RspParamError("id参数错误")
	}

	err := crontab.PauseJob(id)
	if err != nil {
		c.RspFail(err)
	} else {
		c.RspSuccess("暂停成功")
	}
}

// 恢复任务执行
func (c *JobsController) Resume() {
	id, _ := c.GetInt("id")
	if id <= 0 {
		c.RspParamError("id参数错误")
	}

	err := crontab.ResumeJob(id)
	if err != nil {
		c.RspFail(err)
	} else {
		c.RspSuccess("恢复成功")
	}
}

// 停止所有任务
func (c *JobsController) StopJob() {
	crontab.StopAllJobs()
	c.RspSuccess("停止成功")
}

// 开始所有任务
func (c *JobsController) StartJob() {
	crontab.StartAllJobs()
	c.RspSuccess("开始成功")
}

// 所有任务状态
func (c *JobsController) Status() {
	res := crontab.Status()
	data := make(map[string]interface{})
	data["status"] = res
	c.RspData(data)
}
