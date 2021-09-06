package engine

import (
	"com.github/finnlee87/go-job-lite/lib"
	"com.github/finnlee87/go-job-lite/model"
	"database/sql"
	"time"
)

func GetRangeJobs(timeRange int) []model.Job {
	now := time.Now()
	endTime := now.Add(time.Duration(timeRange) * time.Second)
	rows, err := lib.DB.Query("select id, name, description, cron, next_time, status, " +
		"concurrent from job where (next_time < ? or next_time is null) and status = 1", endTime)

	if err != nil {
		lib.Logger.Errorln(err)
		return []model.Job{}
	}
	var jobs []model.Job
	for rows.Next() {
		var job model.Job
		var nextTimeStr string
		rows.Scan(&job.Id, &job.Name, &job.Description, &job.Cron, &nextTimeStr, &job.Status, &job.Concurrent)
		job.NextTime = nextTimeStr
		jobs = append(jobs, job)
	}
	return jobs
}

type DistributedHandler func()

func doDistributedLock(lockName string, inLockVal int, handler DistributedHandler, cleanFunc DistributedHandler)  {
	lib.DoTransaction(func(tx *sql.Tx) error {
		var lockVal int
		row := tx.QueryRow("select lock_val from job_lock where lock_name=? for update", lockName)
		row.Scan(&lockVal)
		if lockVal < inLockVal {
			handler()
		} else {
			cleanFunc()
		}
		//update
		tx.Exec("update job_lock set lock_val=? where lock_name=?", inLockVal, lockName)
		return nil
	})
}

func addJobLog(jobName string, startTime string, status int) (int64, error) {
	exec, _ := lib.DB.Exec("insert into job_log(job_name, start_time, status) values(?, ?, ?)", jobName, startTime, status)
	id, err := exec.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func updateJobLog(jobLogId int64, endTime string, status int, errorMsg string) error {
	_, err := lib.DB.Exec("update job_log set end_time=?, status=?, error_msg=? where id=?", endTime, status, errorMsg, jobLogId)
	if err != nil {
		return err
	}
	return nil
}

func addJob(name string, cron string) {
	var id int64
	row := lib.DB.QueryRow("select id from job where name=?", name)
	row.Scan(&id)
	if id != 0 {
		lib.DB.Exec("update job set name=?, cron=? where id=?", name, cron, id)
	} else {
		lib.DB.Exec("insert into job(name, cron) values(?,?)", name, cron)
	}
}
