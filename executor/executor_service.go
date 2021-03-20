package executor

import (
	"fmt"
	"lendo-sdk/utility"
	"log"
	"sync"
	"time"
)

type ExecutorService struct {
	wg *sync.WaitGroup
	jobChan chan Job
}

func NewExecutorService(workers int, queueLength int) ExecutorService {
	es := ExecutorService{
		&sync.WaitGroup{},
		make(chan Job, queueLength),
	}
	es.initialize(workers)
	return es
}

func (es *ExecutorService) initialize(workers int) {
	es.wg.Add(workers)
	actualWorkers := utility.Max(workers, 1)
	for i := 0; i < actualWorkers; i++ {
		go es.worker(es.jobChan, es.wg)
	}
}

func (es *ExecutorService) shutdown() {
	close(es.jobChan)
	es.wg.Wait()
}

func (es *ExecutorService) queueJob(job Job) bool {
	select {
	case es.jobChan <- job:
		return true
	default:
		return false
	}
}

func (es *ExecutorService) worker(jobChan <-chan Job, wg *sync.WaitGroup) {
	// As soon as the current goroutine finishes (job done!), notify back WaitGroup.
	defer wg.Done()

	fmt.Println("Worker is waiting for jobs")

	for job := range jobChan {
		fmt.Println("Worker picked job", job)
		start := time.Now()

		job.Execute()

		elapsed := time.Since(start)
		log.Printf("Job took %s", elapsed)
	}
}