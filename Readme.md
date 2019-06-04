## 说明
- 基于beego框架
- 引用cron library  [https://github.com/robfig/cron](https://github.com/robfig/cron)
- 参考代码 [https://github.com/george518/PPGo_Job](https://github.com/george518/PPGo_Job)

## http api
```bash
# 添加任务
/jobs/add?cmd="xx"

# 删除任务
/jobs/remove?id=1

# 暂停任务
/jobs/pause?id=1

# 恢复任务
/jobs/resume?id=1

# 获取任务状态
/jobs/status

# 停止所有任务
/jobs/stop

# 开始任务
/jobs/start
```