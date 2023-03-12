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
)

func main() {
	filename := fmt.Sprintf(targetFilenameFormat, 0)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC|os.O_SYNC, 0644)
	if err != nil {
		log.Fatalf("failed to open file: %s", filename)
	}

	var b [blockSize]byte
	count := fileSize / blockSize

	lats := make([]time.Duration, count)
	for i := 0; i < count; i++ {
		start := time.Now()
		_, err := f.Write(b[:])
		if err != nil && err != io.ErrShortWrite {
			log.Fatalf("failed to write to file: %s", filename)
		}
		end := time.Now()
		lats[i] = end.Sub(start)
	}
	sort.Slice(lats, func(i, j int) bool { return lats[i] < lats[j] })

	var latSum time.Duration
	for _, lat := range lats {
		latSum += lat
	}
	writtenMiB := float64(count*blockSize) / 1024 / 1024
	throughputMiBs := writtenMiB / float64(latSum.Microseconds()) * 1000 * 1000
	iops := float64(count) / float64(latSum.Microseconds()) * 1000 * 1000

	fmt.Printf("block size: %d byte\n", blockSize)
	fmt.Printf("total bytes written: %.02f MiB\n", writtenMiB)
	fmt.Printf("throughput: %.02f MiB/s, %.02f IOPS\n", throughputMiBs, iops)
	fmt.Printf("latency:\n")
	fmt.Printf("  avg: %d usec\n", latSum.Microseconds()/int64(count))
	fmt.Printf("  50%%: %d usec\n", lats[count/2].Microseconds())
	fmt.Printf("  90%%: %d usec\n", lats[count*90/100].Microseconds())
	fmt.Printf("  99%%: %d usec\n", lats[count*99/100].Microseconds())
}
