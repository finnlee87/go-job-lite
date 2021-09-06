package model

type JobLog struct {
	BaseModel
	JobId int64 `json:"job_id"`
	Status int `json:"status"`
	StartTime string `json:"start_time"`
	EndTime string `json:"end_time"`
	ErrorMsg string `json:"error_msg"`
}
