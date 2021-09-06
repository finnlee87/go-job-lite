package model

type Job struct {
	BaseModel
	Name string `json:"name"`
	Description string `json:"description"`
	Cron string `json:"cron"`
	NextTime string `json:"next_time"`
	Status int `json:"status"`
	Concurrent int `json:"concurrent"`
}
