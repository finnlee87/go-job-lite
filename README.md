# go-job-lite

### Install

```shell
go get com.github/finnlee87/go-job-lite
```

![Go Version](https://img.shields.io/badge/go%20version-%3E=1.14-61CFDD.svg?style=flat-square)
[![Join the chat at https://gitter.im/finn87/go-job-lite](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/finn87/go-job-lite?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## How to use

```go
func main() {
	jobEngine := engine.Default()
	jobEngine.Register(ExampleJob{})
	jobEngine.Run()

	ch := make(chan string)
	<-ch
}
```

Job must implement BaseJob

```go
type ExampleJob struct {

}

func (example ExampleJob) Name() string {
	return "example-job"
}

func (example ExampleJob) Cron() string  {
	return "0/1 * * * * * *"
}

func (example ExampleJob) Execute() error {
	lib.Logger.Infoln("example job run.")
	return nil
}
```