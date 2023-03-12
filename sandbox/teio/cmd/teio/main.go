package main

import (
	"flag"
	"log"
	"os"

	"github.com/daichimukai/x/sandbox/teio"
)

var (
	numJobs   = flag.Int("num-jobs", 1, "the number of jobs which run concurrently")
	blockSize = flag.Int("block-size", 4096, "size of I/O unit")
	fileSize  = flag.Int("file-size", 1*1024*1024, "do I/Os up to this size in bytes")
)

func main() {
	flag.Parse()

	s, err := teio.NewScenario(*numJobs, *blockSize, *fileSize)
	if err != nil {
		log.Fatalf("failed to initialize a scenario: %v", err)
	}
	r, err := s.Do()
	if err != nil {
		log.Fatalf("failed to complate the scenario: %v", err)
	}
	r.PrettyPrint(os.Stdout)
}
