package teio

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type ScenarioResult struct {
	results []*JobResult
}

func (s ScenarioResult) PrettyPrint(w io.Writer) {
	for _, result := range s.results {
		result.PrettyPrint(os.Stdout)
	}
}

type Scenario struct {
	jobs []*Job
}

func NewScenario(numJobs, blockSize, fileSize int) (*Scenario, error) {
	jobs := make([]*Job, numJobs)

	for i := 0; i < numJobs; i++ {
		job, err := NewJob(i, blockSize, fileSize)
		if err != nil {
			return nil, fmt.Errorf("failed to create job: %w", err)
		}
		jobs[i] = job
	}

	return &Scenario{
		jobs: jobs,
	}, nil
}

func (s *Scenario) Do() (*ScenarioResult, error) {
	results := make([]*JobResult, len(s.jobs))

	startCh := make(chan struct{})
	var wg sync.WaitGroup
	for i, job := range s.jobs {
		wg.Add(1)
		go func(i int, job *Job) {
			<-startCh

			// error is intentionally ignored. @result holds it.
			result, _ := job.Do()
			results[i] = result
			wg.Done()
		}(i, job)
	}
	close(startCh)
	wg.Wait()

	for _, r := range results {
		if r.err != nil {
			return nil, fmt.Errorf("failed to complete jobs: %w", r.err)
		}
	}

	return &ScenarioResult{
		results: results,
	}, nil
}
