package iperf

import (
	"fmt"
	"github.com/BGrewell/go-conversions"
	"github.com/BGrewell/tail"
	"log"
	"strconv"
	"strings"
	"time"
)

/*
Connecting to host 127.0.0.1, port 5201
[  5] local 127.0.0.1 port 49759 connected to 127.0.0.1 port 5201
[ ID] Interval           Transfer     Bitrate
[  5]   0.00-1.00   sec   128 KBytes  1.05 Mbits/sec
[  5]   1.00-2.00   sec   128 KBytes  1.05 Mbits/sec
[  5]   2.00-3.00   sec   128 KBytes  1.05 Mbits/sec
[  5]   3.00-4.00   sec   128 KBytes  1.05 Mbits/sec
[  5]   4.00-5.00   sec   128 KBytes  1.05 Mbits/sec
[  5]   5.00-6.00   sec   128 KBytes  1.05 Mbits/sec
[  5]   6.00-7.00   sec   128 KBytes  1.05 Mbits/sec
[  5]   7.00-8.00   sec   128 KBytes  1.05 Mbits/sec
[  5]   8.00-9.00   sec   128 KBytes  1.05 Mbits/sec
[  5]   9.00-10.00  sec   128 KBytes  1.05 Mbits/sec
- - - - - - - - - - - - - - - - - - - - - - - - -
[ ID] Interval           Transfer     Bitrate
[  5]   0.00-10.00  sec  1.25 MBytes  1.05 Mbits/sec                  sender
[  5]   0.00-10.00  sec  1.25 MBytes  1.05 Mbits/sec                  receiver

iperf Done.
*/
func (r *Reporter) runLogProcessor() {
	var err error
	r.tailer, err = tail.TailFile(r.LogFile, tail.Config{
		Follow:    true,
		ReOpen:    true,
		Poll:      false, // on linux we don't need to poll as the fsnotify works properly
		MustExist: true,
	})
	if err != nil {
		log.Fatalf("failed to tail log file: %v", err)
	}

	for {
		select {
		case line := <-r.tailer.Lines:
			if line == nil {
				continue
			}
			r.LineChannel <- line.Text
			if len(line.Text) > 5 {
				id := line.Text[1:4]
				stream, err := strconv.Atoi(strings.TrimSpace(id))
				if err != nil {
					continue
				}
				fields := strings.Fields(line.Text[5:])
				if len(fields) >= 5 {
					if fields[0] == "local" {
						continue
					}
					timeFields := strings.Split(fields[0], "-")
					start, err := strconv.ParseFloat(timeFields[0], 32)
					if err != nil {
						log.Printf("failed to convert start time: %s\n", err)
					}
					end, err := strconv.ParseFloat(timeFields[1], 32)
					transferredStr := fmt.Sprintf("%s%s", fields[2], fields[3])
					transferredBytes, err := conversions.StringBitRateToInt(transferredStr)
					if err != nil {
						log.Printf("failed to convert units: %s\n", err)
					}
					transferredBytes = transferredBytes / 8
					rateStr := fmt.Sprintf("%s%s", fields[4], fields[5])
					rate, err := conversions.StringBitRateToInt(rateStr)
					if err != nil {
						log.Printf("failed to convert units: %s\n", err)
					}
					report := &StreamIntervalReport{
						Socket:        stream,
						StartInterval: float32(start),
						EndInterval:   float32(end),
						Seconds:       float32(end - start),
						Bytes:         int(transferredBytes),
						BitsPerSecond: float64(rate),
					}
					r.ReportingChannel <- report
				}
			}
		case <-time.After(100 * time.Millisecond):
			if !r.running {
				return
			}
		}
	}
}
