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
	directIO  = flag.Bool("direct", false, "do direct IO")
	rw        = flag.String("rw", "read", "I/O type, read or write (sequential)")
)

func main() {
	flag.Parse()

	var jobType teio.IOType
	switch *rw {
	case "read":
		jobType = teio.IOTypeRead
	case "write":
		jobType = teio.IOTypeWrite
	default:
		log.Fatalf("unknown I/O type: %s", *rw)
	}

	s, err := teio.NewScenario(*numJobs, *blockSize, *fileSize, *directIO, jobType)
	if err != nil {
		log.Fatalf("failed to initialize a scenario: %v", err)
	}
	r, err := s.Do()
	if err != nil {
		log.Fatalf("failed to complate the scenario: %v", err)
	}
	r.PrettyPrint(os.Stdout)
}
