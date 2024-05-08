package utils

import (
	"sync"
	"time"
)

type Action = func()

type CronStatus int

type Cron struct {
	action   Action
	interval time.Duration
	status   CronStatus

	stopSignal chan bool
	mutex      sync.Mutex
}

func NewCron(a Action, interval time.Duration) *Cron {
	return &Cron{action: a, interval: interval, status: CRON_STOPPED, stopSignal: make(chan bool)}
}

func (c *Cron) Start() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	go func() {
		for {
			c.action()
			select {
			case <-c.stopSignal:
				c.status = CRON_STOPPED
				return
			case <-time.After(c.interval):
				continue
			}
		}
	}()
	c.status = CRON_RUNNING
}

func (c *Cron) Stop() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.stopSignal <- true
	for c.status != CRON_STOPPED {
		time.Sleep(10 * time.Millisecond)
	}
}

const (
	CRON_STOPPED CronStatus = iota
	CRON_RUNNING
)
