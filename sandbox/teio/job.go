package teio

import (
	"fmt"
	"io"
	"os"
	"sort"
	"syscall"
	"time"
	"unsafe"
)

const targetFilenameFormat = "teio.%d"

type Job struct {
	id        int
	blockSize int
	fileSize  int
	directIO  bool
	fp        *os.File
}

func NewJob(id int, blockSize int, fileSize int, directIO bool) (*Job, error) {
	filename := fmt.Sprintf(targetFilenameFormat, id)

	oFlag := os.O_CREATE | os.O_RDWR | os.O_TRUNC
	if directIO {
		oFlag |= syscall.O_DIRECT
	}
	fp, err := os.OpenFile(filename, oFlag, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return &Job{
		id:        id,
		blockSize: blockSize,
		fileSize:  fileSize,
		directIO:  directIO,
		fp:        fp,
	}, nil
}

func (j *Job) Do() (*JobResult, error) {
	count := j.fileSize / j.blockSize
	lats := make([]time.Duration, count)

	b := makeAlignedSlice(j.blockSize, 4096)
	for i := 0; i < count; i++ {
		start := time.Now()
		_, err := j.fp.Write(b[:])
		if err != nil {
			return &JobResult{err: err}, err
		}
		end := time.Now()
		lats[i] = end.Sub(start)
	}
	sort.Slice(lats, func(i, j int) bool { return lats[i] < lats[j] })

	j.fp.Close()

	return &JobResult{job: j, lats: lats}, nil
}

type JobResult struct {
	job  *Job
	lats []time.Duration

	err error
}

func (j JobResult) PrettyPrint(w io.Writer) {
	var latSum time.Duration
	for _, lat := range j.lats {
		latSum += lat
	}
	count := len(j.lats)
	writtenMiB := float64(count*j.job.blockSize) / 1024 / 1024
	throughputMiBs := writtenMiB / float64(latSum.Microseconds()) * 1000 * 1000
	iops := float64(count) / float64(latSum.Microseconds()) * 1000 * 1000

	fmt.Fprintf(w, "job%d: \n", j.job.id)
	fmt.Fprintf(w, "  block size: %d byte\n", j.job.blockSize)
	fmt.Fprintf(w, "  total bytes written: %.02f MiB\n", writtenMiB)
	fmt.Fprintf(w, "  throughput: %.02f MiB/s, %.02f IOPS\n", throughputMiBs, iops)
	fmt.Fprintf(w, "  latency:\n")
	fmt.Fprintf(w, "    avg: %d usec\n", latSum.Microseconds()/int64(count))
	fmt.Fprintf(w, "    50%%: %d usec\n", j.lats[count/2].Microseconds())
	fmt.Fprintf(w, "    90%%: %d usec\n", j.lats[count*90/100].Microseconds())
	fmt.Fprintf(w, "    99%%: %d usec\n", j.lats[count*99/100].Microseconds())
}

// makeAlignedSlice returns a slice with @size bytes that aligned to
// @alignment. This is needed since syscalls for fd with O_DIRECT may require
// an alignment limitation for a user-space buffer that varies by filesystems
// and kernels.
func makeAlignedSlice(size int, alignment int) []byte {
	buf := make([]byte, size+alignment)
	violation := uintptr(unsafe.Pointer(unsafe.SliceData(buf))) & uintptr(alignment-1)
	if violation == 0 {
		return buf[:size]
	}
	start := alignment - int(violation)
	end := start + size
	return buf[start:end]
}
