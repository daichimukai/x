teio
====

Simple I/O tester written for technical exam.

Example
-------


```
$ go run cmd/teio/main.go
job0:
  IO type: read
  direct IO: false
  block size: 4096 byte
  total bytes written: 1.00 MiB
  throughput: 2493.77 MiB/s, 638403.99 IOPS
  latency:
    avg: 1 usec
    50%: 1 usec
    90%: 1 usec
    99%: 18 usec
$
```

```
$ go run cmd/teio/main.go --num-jobs=4 --block-size=$((64*1024)) --file-size=$((2*1024*1024)) --rw=write --direct
job0:
  IO type: write
  direct IO: true
  block size: 65536 byte
  total bytes written: 2.00 MiB
  throughput: 98.78 MiB/s, 1580.56 IOPS
  latency:
    avg: 632 usec
    50%: 496 usec
    90%: 524 usec
    99%: 4845 usec
job1:
  IO type: write
  direct IO: true
  block size: 65536 byte
  total bytes written: 2.00 MiB
  throughput: 100.00 MiB/s, 1599.92 IOPS
  latency:
    avg: 625 usec
    50%: 497 usec
    90%: 499 usec
    99%: 4588 usec
job2:
  IO type: write
  direct IO: true
  block size: 65536 byte
  total bytes written: 2.00 MiB
  throughput: 98.36 MiB/s, 1573.80 IOPS
  latency:
    avg: 635 usec
    50%: 496 usec
    90%: 499 usec
    99%: 4932 usec
job3:
  IO type: write
  direct IO: true
  block size: 65536 byte
  total bytes written: 2.00 MiB
  throughput: 98.98 MiB/s, 1583.69 IOPS
  latency:
    avg: 631 usec
    50%: 496 usec
    90%: 502 usec
    99%: 4804 usec
$
```
