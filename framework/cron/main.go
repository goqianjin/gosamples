package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New(cron.WithSeconds(), cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))
	fmt.Printf("start....\n")
	c.AddFunc("0/1 * * * * *", func() {
		fmt.Println("job1 --> start at " + time.Now().Format("2006-01-02 15:04:05") + " : every 1 seconds executing")
		time.Sleep(3 * time.Second)
		fmt.Println("job1 --> ended at " + time.Now().Format("2006-01-02 15:04:05"))

	})
	c.AddFunc("@every 1s", func() {
		fmt.Println("job2 --> " + time.Now().Format("2006-01-02 15:04:05") + " : every 1 seconds executing")
		time.Sleep(3 * time.Second)
		fmt.Println("job2 --> ended at " + time.Now().Format("2006-01-02 15:04:05"))

	})

	go c.Start()
	defer c.Stop()

	select {
	case <-time.After(time.Second * 10):
		return
	}
	fmt.Printf("end....")
}
