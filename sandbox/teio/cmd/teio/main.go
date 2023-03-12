package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"time"
)

const (
	targetFilenameFormat = "teio.%d"

	blockSize = 4096
	fileSize  = 1 * 1024 * 1024
	numJobs   = 1
)

type Job struct {
	fp *os.File
}

func NewJob(id uint) (*Job, error) {
	filename := fmt.Sprintf(targetFilenameFormat, 0)
	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC|os.O_SYNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return &Job{fp: fp}, nil
}

func (j *Job) Do() (*JobResult, error) {
	count := fileSize / blockSize
	lats := make([]time.Duration, count)

	var b [blockSize]byte
	for i := 0; i < count; i++ {
		start := time.Now()
		_, err := j.fp.Write(b[:])
		if err != nil && err != io.ErrShortWrite {
			return nil, err
		}
		end := time.Now()
		lats[i] = end.Sub(start)
	}
	sort.Slice(lats, func(i, j int) bool { return lats[i] < lats[j] })

	return &JobResult{lats: lats}, nil
}

type JobResult struct {
	lats []time.Duration
}

func (j JobResult) PrettyPrint(w io.Writer) {
	var latSum time.Duration
	for _, lat := range j.lats {
		latSum += lat
	}
	count := len(j.lats)
	writtenMiB := float64(count*blockSize) / 1024 / 1024
	throughputMiBs := writtenMiB / float64(latSum.Microseconds()) * 1000 * 1000
	iops := float64(count) / float64(latSum.Microseconds()) * 1000 * 1000

	fmt.Fprintf(w, "block size: %d byte\n", blockSize)
	fmt.Fprintf(w, "total bytes written: %.02f MiB\n", writtenMiB)
	fmt.Fprintf(w, "throughput: %.02f MiB/s, %.02f IOPS\n", throughputMiBs, iops)
	fmt.Fprintf(w, "latency:\n")
	fmt.Fprintf(w, "  avg: %d usec\n", latSum.Microseconds()/int64(count))
	fmt.Fprintf(w, "  50%%: %d usec\n", j.lats[count/2].Microseconds())
	fmt.Fprintf(w, "  90%%: %d usec\n", j.lats[count*90/100].Microseconds())
	fmt.Fprintf(w, "  99%%: %d usec\n", j.lats[count*99/100].Microseconds())
}

func main() {
	job, err := NewJob(0)
	if err != nil {
		log.Fatal(err)
	}

	result, err := job.Do()
	if err != nil {
		log.Fatal(err)
	}

	result.PrettyPrint(os.Stdout)
}
