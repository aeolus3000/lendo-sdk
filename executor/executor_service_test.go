package executor

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

const (
	gracePeriod = 1 * time.Second
)

type TestJob struct {
	WaitTime time.Duration
	Counter  *Counter
}

func (tj *TestJob) Execute() {
	fmt.Println("Waiting for ", tj.WaitTime)
	time.Sleep(tj.WaitTime)
	tj.Counter.Count()
}

func (tj *TestJob) String() string {
	return fmt.Sprintf("Job %fs", tj.WaitTime.Seconds())
}

type Counter struct {
	count uint64
}

func (c *Counter) Count() {
	atomic.AddUint64(&c.count, 1)
}

func TestNewExecutorService(t *testing.T) {
	type args struct {
		workers     int
		queueLength int
		jobs        uint64
		jobDuration time.Duration
	}
	tests := []struct {
		name     string
		args     args
		expected uint64
	}{
		{"Working fine",
			args{
				workers:     4,
				queueLength: 4,
				jobs:        4,
				jobDuration: 2 * time.Second,
			},
			4},
		{"Not enough time",
			args{
				workers:     4,
				queueLength: 4,
				jobs:        6,
				jobDuration: 2 * time.Second,
			},
			4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			es := NewExecutorService(tt.args.workers, tt.args.queueLength)
			defer es.shutdown()
			c := Counter{0}
			scheduleJobs(&es, &c, tt.args.jobs, tt.args.jobDuration)
			timeout := tt.args.jobDuration + gracePeriod
			evaluateTimeout(timeout, &c, tt.args.jobs)
			if c.count != tt.expected {
				t.Errorf("Expected %d jobs to finish but only %d job(s) finished", tt.expected, c.count)
			}
		})
	}
}

func evaluateTimeout(timeout time.Duration, counter *Counter, expectedCount uint64) {
	done := make(chan int, 1)
	go waitForResult(done, counter, expectedCount)
	select {
	case <-done:
		break
	case <-time.After(timeout):
		break
	}
}

func waitForResult(done chan int, counter *Counter, expectedCount uint64) {
	for {
		if counter.count == expectedCount {
			done <- 0
			break
		}
	}
}

func scheduleJobs(es *ExecutorService, counter *Counter, jobs uint64, jobDuration time.Duration) {
	var i uint64
	for i = 0; i < jobs; i++ {
		tj := TestJob{
			WaitTime: jobDuration,
			Counter:  counter,
		}
		es.queueJob(&tj)
	}
}
