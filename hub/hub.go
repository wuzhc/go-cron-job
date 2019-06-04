package hub

import (
	"cron-job/cron"
	"errors"
	"sync"

	"github.com/astaxie/beego"
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
func RemoveJob(id int) {
	cronManger.RemoveJob(func(e *cron.Entry) bool {
		if j, ok := e.Job.(*Job); ok {
			if j.Id == id {
				return true
			}
		}
		return false
	})
}

// 停止所有任务
func StopAllJobs() {
	cronManger.Stop()
}

// 开始所有任务
func StartAllJobs() {
	cronManger.Start()
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

type jobInfo struct {
	Id         int    `json:"id"`
	Cmd        string `json:"cmd"`
	Status     int    `json:"status"` //0准备,1进行中,2暂停
	FailNum    int    `json:"fail_num"`
	SuccessNum int    `json:"success_num"`
}

// 所有任务状态
func Status() []jobInfo {
	res := make([]jobInfo, 0)

	entries := cronManger.Entries()
	for _, e := range entries {
		if j, ok := e.Job.(*Job); ok {
			info := jobInfo{
				Id:         j.Id,
				Cmd:        j.Cmd,
				Status:     j.Status,
				FailNum:    j.FailNum,
				SuccessNum: j.SuccessNum,
			}
			res = append(res, info)
		}
	}

	return res
}
