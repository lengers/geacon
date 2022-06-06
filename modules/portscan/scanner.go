package modules

import (
	"fmt"
	"log"
	"time"

	"github.com/XinRoom/go-portScan/core/host"
	"github.com/XinRoom/go-portScan/core/port"
	"github.com/XinRoom/iprange"
)

type greeting string

type Parameters struct {
	targets       string
	ports         string
	discoverytype string
	socketcount   int
}

func (params Parameters) Module() {
	fmt.Println("Hello Universe")
	single := make(chan struct{})
	retChan := make(chan port.OpenIpPort, 65535)
	go func() {
		for {
			select {
			case ret := <-retChan:
				if ret.Port == 0 {
					single <- struct{}{}
					return
				}
				fmt.Println(ret)
			default:
				time.Sleep(time.Millisecond * 10)
			}
		}
	}()

	ports, err := port.ParsePortRangeStr(params.ports)
	if err != nil {
		fmt.Println(err)
	}

	// parse ip
	it, _, _ := iprange.NewIter(params.targets)

	// scanner
	ss, err := port.NewTcpScanner(retChan, port.DefaultTcpOption)
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	for i := uint64(0); i < it.TotalNum(); i++ { // ip索引
		ip := it.GetIpByIndex(i)
		if !host.IsLive(ip.String()) { // ping
			continue
		}
		for _, _porterange := range ports { // port
			for _, _port := range _porterange { // port
				ss.WaitLimiter()
				ss.Scan(ip, _port) // syn
			}
		}
	}
	ss.Close()
	fmt.Println(time.Since(start))
}
