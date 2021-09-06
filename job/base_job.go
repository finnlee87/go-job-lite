package job

type BaseJob interface {
	Name() string
	Cron() string
	Execute() error
}
