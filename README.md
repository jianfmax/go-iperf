# go-iperf
A Go based wrapper around iperf3

```shell
go get -u github.com/jianfmax/go-iperf
```

## Basic Usage

basic client setup
```go
func main() {
	
	c := iperf.NewClient("192.168.0.10")
	c.SetJSON(true)
	c.SetIncludeServer(true)
	c.SetStreams(4)
	c.SetTimeSec(30)
	c.SetInterval(1)
	
	err := c.Start()
	if err != nil {
        fmt.Printf("failed to start client: %v\n", err)
        os.Exit(-1)
	}
	
	<- c.Done
	
	fmt.Println(c.Report().String())
}
```

basic server setup
```go
func main() {
	
	s := iperf.NewServer()
	err := s.Start()
	if err != nil {
        fmt.Printf("failed to start server: %v\n", err)
        os.Exit(-1)
    }
    
    for s.Running() {
    	time.Sleep(100 * time.Millisecond)
    }
    
    fmt.Println("server finished")
}
```

client with live results printing
```go
func main() {
	
	c := iperf.NewClient("192.168.0.10")
	c.SetJSON(true)
	c.SetIncludeServer(true)
	c.SetStreams(4)
	c.SetTimeSec(30)
	c.SetInterval(1)
	liveReports := c.SetModeLive()
	
	go func() {
	    for report := range liveReports {
	        fmt.Println(report.String())	
            }   	
        }   
	
	err := c.Start()
	if err != nil {
            fmt.Printf("failed to start client: %v\n", err)
            os.Exit(-1)
	}
	
	<- c.Done
	
	fmt.Println(c.Report().String())
}
```

building binary data package with iperf binaries
```
go-bindata -pkg iperf -prefix "embedded/" embedded/
```
