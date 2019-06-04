package crontab

import (
	"errors"
	"sync"

	"github.com/astaxie/beego"

	"github.com/robfig/cron"
)

var (
	cronManger *cron.Cron
	lock       sync.Mutex
	workPool   chan bool
)

func init() {
	if n, _ := beego.AppConfig.Int("jobs.pool"); n > 0 {
		workPool = make(chan bool, n)
	}

	cronManger = cron.New()
	cronManger.Start()
}

// 添加任务
func AddFunc(spec string, cmd func()) error {
	lock.Lock()
	defer lock.Unlock()

	return cronManger.AddFunc(spec, cmd)
}

// 添加任务
func AddJob(spec string, cmd string, timeout int) error {
	lock.Lock()
	defer lock.Unlock()

	return cronManger.AddJob(spec, NewJob(cmd, timeout))
}

// 暂停任务
func PauseJob(id int) error {
	job := GetJobById(id)
	if job == nil {
		return errors.New("任务不存在")
	}

	job.Status = JOB_PAUSE
	return nil
}

// 恢复任务执行
func ResumeJob(id int) error {
	job := GetJobById(id)
	if job == nil {
		return errors.New("任务不存在")
	}

	job.Status = JOB_READY
	return nil
}

// 删除任务
func RemoveJob(id int) error {
	return nil
}

// 停止所有任务
func StopAllJobs() {
	cronManger.Stop()
}

// 开始所有任务
func StartAllJobs() {
	cronManger.Start()
}

// 所有任务状态
func Status() []map[string]int {
	res := make([]map[string]int, 0)

	entries := cronManger.Entries()
	for _, e := range entries {
		if j, ok := e.Job.(*Job); ok {
			var jobInfo = make(map[string]int)
			jobInfo["failNum"] = j.FailNum
			jobInfo["successNum"] = j.SuccessNum
			jobInfo["jobId"] = j.Id
			res = append(res, jobInfo)
		}
	}

	return res
}

// 获取所有任务
func GetAllJobs() []*Job {
	entries := cronManger.Entries()
	jobs := make([]*Job, 0)

	for _, e := range entries {
		if j, ok := e.Job.(*Job); ok {
			jobs = append(jobs, j)
		}
	}

	return jobs
}

// 根据任务ID获取指定任务
func GetJobById(id int) *Job {
	jobs := GetAllJobs()
	for _, j := range jobs {
		if j.Id == id {
			return j
		}
	}

	return nil
}
