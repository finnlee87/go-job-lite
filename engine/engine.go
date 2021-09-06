package engine

import (
	"context"
	"fmt"
	"github.com/finnlee87/go-job-lite/job"
	"github.com/finnlee87/go-job-lite/lib"
	"github.com/gin-gonic/gin"
	"github.com/gorhill/cronexpr"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const secondPerMinute = 60

type Engine struct {
	TimeThreshold   int
	JobLoopInterval int
	JobMap          map[string]job.BaseJob
	TimeWheel       [60] []string
	Mutex           *sync.Mutex
}

func Default() *Engine {
	return &Engine{JobLoopInterval: 5, TimeThreshold: 5, JobMap: map[string]job.BaseJob{}, Mutex: &sync.Mutex{}}
}

func (engine *Engine) Register(baseJob job.BaseJob) {
	engine.JobMap[baseJob.Name()] = baseJob
	//update db
	addJob(baseJob.Name(), baseJob.Cron())
}

func (engine *Engine) Run() {
	r := gin.Default()
	//add http api.
	r.GET("/jobs")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGSTOP)

	ctx, cancel := context.WithCancel(context.Background())
	go func(cancel context.CancelFunc, sig chan os.Signal) {
		<-sig
		cancel()
	}(cancel, sig)

	now := time.Now().Truncate(time.Second)
	second := now.Second()
	period := secondPerMinute - second
	timeTip := int(float64(engine.JobLoopInterval) / 2)
	putTimeWheelTime := secondPerMinute - timeTip
	if second >= putTimeWheelTime {
		period = period + secondPerMinute
	}
	sleep := period - timeTip

	putTimeWheel(now.Add(time.Duration(5) * time.Second), engine, float64(period - 5))

	go func(ctx context.Context, sleep int, engine *Engine) {
		time.Sleep(time.Duration(sleep)* time.Second)
		//计算next time 并放入time wheel
		for {
			select {
				case <- ctx.Done():
					lib.Logger.Infoln("stop put job to time wheel.")
					return
			default:
				lib.Logger.Debugln("start put job to time ring.")
				t := time.Now()
				putTimeWheel(t.Add(time.Duration(1) * time.Minute).Truncate(time.Minute).Add(time.Duration(-1)*time.Second), engine, secondPerMinute)
 				time.Sleep(t.Add(time.Duration(2) * time.Minute).Truncate(time.Minute).Add(time.Duration(-timeTip)*time.Second).Sub(t))
			}
		}

	}(ctx, sleep, engine)

	//run job
	go func(ctx context.Context, sleep int, engine *Engine) {
		//time.Sleep(time.Duration(sleep)* time.Second)
		for {
			lib.Logger.Debugln("dispatch job.")
			select {
			case <-ctx.Done():
				lib.Logger.Infoln("stop run job.")
				return
			default:
				now := time.Now()
				second := now.Second()
				lockSecond := second / engine.JobLoopInterval * engine.JobLoopInterval
				lib.Logger.Debugf("second %d", second)
				var secondStr string
				if second < 10 {
					secondStr = "0" + strconv.Itoa(lockSecond)
				} else {
					secondStr = strconv.Itoa(lockSecond)
				}
				lockVal := fmt.Sprintf("%s%s", now.Format("200601021504"), secondStr)
				lockValInt, _ := strconv.Atoi(lockVal)
				doDistributedLock("job_execute", lockValInt, func() {
					for i := 0; i < (engine.JobLoopInterval - second % engine.JobLoopInterval); i++ {
						engine.Mutex.Lock()
						names := engine.TimeWheel[second + i]
						for _, name := range names {
							go func(name string, engine *Engine, sleep int) {
								time.Sleep(time.Duration(sleep) * time.Second)
								baseJob := engine.JobMap[name]
								if baseJob != nil {
									jobLogId, _ := addJobLog(name, time.Now().Format("2006-01-02 15:04:05"), 1)
									err := baseJob.Execute()
									if err != nil {
										updateJobLog(jobLogId,  time.Now().Format("2006-01-02 15:04:05"), 0, err.Error())
									} else {
										updateJobLog(jobLogId,  time.Now().Format("2006-01-02 15:04:05"), 1, "")
									}
								}
							}(name, engine, i)
						}
						engine.TimeWheel[second + i] = make([]string, 0)
						engine.Mutex.Unlock()
					}
				}, func() {
					for i := 0; i < (engine.JobLoopInterval - second % engine.JobLoopInterval); i++ {
						engine.Mutex.Lock()
						engine.TimeWheel[second + i] = make([]string, 0)
						engine.Mutex.Unlock()
					}
				})

				t := time.Now().Truncate(time.Second).Add(time.Duration(engine.JobLoopInterval) * time.Second)
				time.Sleep(t.Add(time.Duration(0 - t.Second() % engine.JobLoopInterval) * time.Second).Sub(now))
			}
		}
	}(ctx, secondPerMinute - second, engine)

}

func putTimeWheel(time time.Time, engine *Engine, period float64) {
	for _, baseJob := range engine.JobMap {
		cron := baseJob.Cron()
		nextTime := cronexpr.MustParse(cron).Next(time)
		if nextTime.Sub(time).Seconds() < period {
			engine.Mutex.Lock()
			engine.TimeWheel[nextTime.Second()] = append(engine.TimeWheel[nextTime.Second()], baseJob.Name())
			engine.Mutex.Unlock()
		}
		for {
			nextTime = cronexpr.MustParse(cron).Next(nextTime)
			if nextTime.Sub(time).Seconds() < period {
				engine.Mutex.Lock()
				engine.TimeWheel[nextTime.Second()] = append(engine.TimeWheel[nextTime.Second()], baseJob.Name())
				engine.Mutex.Unlock()
			} else {
				break
			}
		}
	}
}






