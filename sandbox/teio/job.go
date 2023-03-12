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

//go:generate stringer -type=IOType -linecomment
type IOType int

const (
	IOTypeRead  IOType = iota // read
	IOTypeWrite               // write
)

const targetFilenameFormat = "teio.%d"

type Job struct {
	id        int
	blockSize int
	fileSize  int
	directIO  bool
	ioType    IOType
	fp        *os.File
}

func NewJob(id int, blockSize int, fileSize int, directIO bool, ioType IOType) (*Job, error) {
	filename := fmt.Sprintf(targetFilenameFormat, id)

	oFlag := os.O_CREATE | os.O_RDWR | os.O_TRUNC
	if directIO {
		oFlag |= syscall.O_DIRECT
	}
	fp, err := os.OpenFile(filename, oFlag, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	if err := fp.Truncate(int64(fileSize)); err != nil {
		return nil, fmt.Errorf("failed to truncate file: %w", err)
	}

	return &Job{
		id:        id,
		blockSize: blockSize,
		fileSize:  fileSize,
		directIO:  directIO,
		ioType:    ioType,
		fp:        fp,
	}, nil
}

func doWrite(fp *os.File, blockSize, count int) ([]time.Duration, error) {
	lats := make([]time.Duration, count)

	b := makeAlignedSlice(blockSize, 4096)
	for i := 0; i < count; i++ {
		start := time.Now()
		_, err := fp.Write(b[:])
		if err != nil {
			return nil, err
		}
		end := time.Now()
		lats[i] = end.Sub(start)
	}

	return lats, nil
}

func doRead(fp *os.File, blockSize, count int) ([]time.Duration, error) {
	lats := make([]time.Duration, count)

	b := make([]byte, count)
	for i := 0; i < count; i++ {
		start := time.Now()
		_, err := fp.Read(b[:])
		if err != nil {
			return nil, err
		}
		end := time.Now()
		lats[i] = end.Sub(start)
	}

	return lats, nil
}

func (j *Job) Do() (*JobResult, error) {
	defer j.fp.Close()

	count := j.fileSize / j.blockSize

	var f func(*os.File, int, int) ([]time.Duration, error)
	switch j.ioType {
	case IOTypeRead:
		f = doRead
	case IOTypeWrite:
		f = doWrite
	default:
		panic("unreachable")
	}

	lats, err := f(j.fp, j.blockSize, count)
	if err != nil {
		return &JobResult{err: err}, err
	}
	sort.Slice(lats, func(i, j int) bool { return lats[i] < lats[j] })

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
	fmt.Fprintf(w, "  IO type: %s\n", j.job.ioType)
	fmt.Fprintf(w, "  direct IO: %v\n", j.job.directIO)
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
