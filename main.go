package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

type Context struct {
	customerID int64
}

func main() {
	var redisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "notifications_redis:6379")
		},
	}

	pool := work.NewWorkerPool(Context{}, 10, "notifications", redisPool)

	// Add middleware that will be executed for each job
	pool.Middleware((*Context).log)

	// Map the name of jobs to handler functions
	pool.Job("sms", (*Context).SendSMS)
	pool.Job("push", (*Context).PushToMobile)

	// Start processing jobs
	pool.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	// Stop the pool
	pool.Stop()
}

func (c *Context) log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println("Starting job: ", job.Name)
	return next()
}
