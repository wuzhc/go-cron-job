package crontab

import (
	"bytes"
	"os/exec"
	"time"

	"github.com/astaxie/beego"
)

const (
	JOB_READY = iota
	JOB_RUNNING
	JOB_PAUSE
)

type Job struct {
	Id           int
	Cmd          string
	Timeout      time.Duration
	Status       int //0准备,1进行中,2暂停
	FailNum      int
	SuccessNum   int
	failChan     chan struct{}
	successChan  chan struct{}
	IsConcurrent bool // 如果上一个任务未执行完毕,是否允许执行下一个任务
}

var (
	jobid int           = 0
	jch   chan struct{} = make(chan struct{}, 1)
)

// 生成jobid,并发安全
func GenerateJobId() int {
	jch <- struct{}{}
	defer func() {
		<-jch
	}()

	jobid++
	return jobid
}

// 新建任务初始化
func NewJob(cmd string, timeout int) *Job {
	j := &Job{
		Id:           GenerateJobId(),
		Cmd:          cmd,
		Timeout:      time.Duration(timeout) * time.Second,
		successChan:  make(chan struct{}),
		failChan:     make(chan struct{}),
		IsConcurrent: false,
	}

	// 启动一个goroutine记录成功或失败次数,并发安全
	go func() {
		for {
			select {
			case <-j.failChan:
				j.FailNum++
			case <-j.successChan:
				j.SuccessNum++
			}
		}
	}()

	return j
}

// 执行任务,并且任务设置超时时间
func (j *Job) Run() {
	// 所有job并发个数,如果有设置并发个数
	if workPool != nil {
		workPool <- true
		defer func() {
			<-workPool
		}()
	}

	if j.Status == JOB_PAUSE {
		return
	}

	if j.Status == JOB_RUNNING && j.IsConcurrent == false {
		beego.Warn("不允许并发执行任务", "任务ID为:", j.Id)
		return
	}

	j.Status = JOB_RUNNING
	defer func() {
		j.Status = JOB_READY
	}()

	var stdout, stderr bytes.Buffer
	c := exec.Command("sh", "-c", j.Cmd)
	c.Stdout = &stdout
	c.Stderr = &stderr

	err := c.Start()
	if err != nil {
		beego.Error(err)
		return
	}

	ch := make(chan error)
	go func() {
		ch <- c.Wait()
	}()

	// 设置超时时间
	select {
	case <-time.After(3 * time.Second):
		pid := c.Process.Pid
		j.failChan <- struct{}{}
		if err := c.Process.Kill(); err != nil {
			beego.Error(err, "超时问题", pid)
		} else {
			beego.Warn("任务超时了", pid)
		}
	case err := <-ch:
		if err != nil {
			j.failChan <- struct{}{}
			beego.Error(err, "等待问题")
		} else {
			j.successChan <- struct{}{}
			beego.Info("正常执行", "任务ID为:", j.Id)
			beego.Info(stdout.String())
		}
	}

	return
}
