package main

import (
	"fmt"
	"github.com/BGrewell/go-conversions"
	"github.com/jianfmax/go-iperf"
	"os"
	"time"
)

func main() {

	//includeServer := true
	proto := "tcp"
	runTime := 5
	omitSec := 0
	//length := "65500"

	c := iperf.NewClient("127.0.0.1")
	//c.SetIncludeServer(includeServer)
	c.SetTimeSec(runTime)
	c.SetOmitSec(omitSec)
	c.SetProto((iperf.Protocol)(proto))
	//c.SetLength(length)
	c.SetJSON(false)
	c.SetIncludeServer(false)
	//c.SetStreams(2)
	c.SetBandwidth("10M")
	reports := c.SetModeLive()

	stopT := time.NewTimer(15 * time.Second)

	go func() {
		for {
			select {
			case report := <-reports:
				fmt.Println(report.String())
			case <-stopT.C:
				break
			}
		}
	}()

	err := c.Start()
	if err != nil {
		fmt.Println("failed to start client")
		os.Exit(-1)
	}

	//time.Sleep(5 * time.Second)
	//c.Stop()

	// Method 1: Wait for the test to finish by pulling from the 'Done' channel which will block until something is put in or it's closed
	<-c.Done

	// Method 2: Poll the c.Running state and wait for it to be 'false'
	//for c.Running {
	//	time.Sleep(100 * time.Millisecond)
	//}

	if c.Report() != nil && c.Report().Error != "" {
		fmt.Println(c.Report().Error)
	} else if c.Report() != nil {
		for _, entry := range c.Report().End.Streams {
			fmt.Println(entry.String())
		}
		for _, entry := range c.Report().ServerOutputJson.End.Streams {
			fmt.Println(entry.String())
		}
		fmt.Printf("DL Rate: %s\n", conversions.IntBitRateToString(int64(c.Report().End.SumReceived.BitsPerSecond)))
		fmt.Printf("UL Rate: %s\n", conversions.IntBitRateToString(int64(c.Report().End.SumSent.BitsPerSecond)))
	}
}
