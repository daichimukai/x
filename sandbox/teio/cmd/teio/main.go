package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/daichimukai/x/sandbox/teio"
)

var (
	numJobs   = flag.Int("num-jobs", 1, "the number of jobs which run concurrently")
	blockSize = flag.Int("block-size", 4096, "size of I/O unit")
	fileSize  = flag.Int("file-size", 1*1024*1024, "do I/Os up to this size in bytes")
)

func main() {
	flag.Parse()

	jobs := make([]*teio.Job, *numJobs)
	results := make([]*teio.JobResult, *numJobs)

	for i := 0; i < *numJobs; i++ {
		job, err := teio.NewJob(i, *blockSize, *fileSize)
		if err != nil {
			log.Fatal(err)
		}
		jobs[i] = job
	}

	startCh := make(chan struct{})
	var wg sync.WaitGroup
	for i, job := range jobs {
		wg.Add(1)
		go func(i int, job *teio.Job) {
			<-startCh

			result, err := job.Do()
			if err != nil {
				log.Fatal(err)
			}

			results[i] = result
			wg.Done()
		}(i, job)
	}
	close(startCh)
	wg.Wait()

	for _, result := range results {
		result.PrettyPrint(os.Stdout)
	}
}
