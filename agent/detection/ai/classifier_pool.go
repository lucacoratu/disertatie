package detection

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/lucacoratu/disertatie/agent/logging"
)

/*
 * ClassifierPool keeps several python classifier processes alive so that
 * multiple requests can be classified in parallel. Each process still handles
 * one request at a time (its pipe is a single stream), so a request waits only
 * when all the processes are busy.
 */
type ClassifierPool struct {
	logger    logging.ILogger
	processes []*ClassifierProcess
	idle      chan *ClassifierProcess //The processes which are not classifying a request at the moment
}

// Creates a pool of size processes running the same classifier script
// If size is not positive, one process is started for each CPU (each process uses about 130 MB of memory)
func NewClassifierPool(logger logging.ILogger, interpreter string, scriptPath string, size int) *ClassifierPool {
	if size <= 0 {
		size = runtime.NumCPU()
	}
	pool := &ClassifierPool{logger: logger, idle: make(chan *ClassifierProcess, size)}
	for i := 0; i < size; i++ {
		pool.processes = append(pool.processes, NewClassifierProcess(logger, interpreter, scriptPath))
	}
	return pool
}

// Start launches all the processes in parallel and waits until every model is loaded
func (pool *ClassifierPool) Start() error {
	errs := make(chan error, len(pool.processes))
	var wg sync.WaitGroup
	for _, process := range pool.processes {
		wg.Add(1)
		go func(process *ClassifierProcess) {
			defer wg.Done()
			errs <- process.Start()
		}(process)
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			pool.Stop()
			return err
		}
	}

	for _, process := range pool.processes {
		pool.idle <- process
	}
	pool.logger.Info(fmt.Sprintf("Started %d classifier processes", len(pool.processes)))
	return nil
}

// Classify waits for an idle process, sends it the features and returns the classification
func (pool *ClassifierPool) Classify(features string) (string, error) {
	process := <-pool.idle
	defer func() { pool.idle <- process }()
	return process.Classify(features)
}

// Stop terminates all the python processes
func (pool *ClassifierPool) Stop() {
	for _, process := range pool.processes {
		process.Stop()
	}
}
