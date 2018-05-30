package main

import (
	"fmt"
	"time"
	"runtime"
	"strconv"

	//"strings"
	
	"github.com/pangudashu/memcache"
	flag "github.com/spf13/pflag"
)

type UserConfig struct {
	IPs       []string
	Clients   uint64
	Requests  uint64
	DataSize  uint64
	Ratio     float64
}

var userConfig UserConfig

func init() {
	flag.Uint64VarP(&userConfig.Clients, "clients", "c", 100, "Number of parallel connections (default 50)")
	flag.Uint64VarP(&userConfig.Requests, "requests", "n", 100000, "Total number of requests (default 100000)")
	flag.Uint64VarP(&userConfig.DataSize, "data-size", "d", 32, "Data size of SET/GET value in bytes (default 2)")
	flag.Float64Var(&userConfig.Ratio, "ratio", 0.1, "Set:Get ratio (default: 1:10), [0.01 - 100]")
	flag.StringSliceVarP(&userConfig.IPs, "addrs", "", []string{"127.0.0.1:11211"}, "Server hostname (default 127.0.0.1:11211), or can set 127.0.0.1:11211,127.0.0.2:11211")
}

func main() {
    flag.Parse()
    fmt.Println("userConfig: ", userConfig)
    runtime.GOMAXPROCS(runtime.NumCPU())

    mcs := []*memcache.Server{}

    for _, ip := range(userConfig.IPs) {
    	mcs = append(mcs, &memcache.Server{
    		Address: ip,
    		Weight: 50,
    		MaxConn: int(userConfig.Clients),
    		InitConn: int(userConfig.Clients) / 2,
    		IdleTime: time.Hour * 2})
    }

	mc, err := memcache.NewMemcache(mcs)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 设置是否自动剔除无法连接的server，默认不开启(建议开启)
	// 如果开启此选项被踢除的server如果恢复正常将会再次被加入server列表
	mc.SetRemoveBadServer(true)

	a := &Attacker{
	 	stopch:  make(chan struct{}), 
	 	workers: userConfig.Clients,
	 	mc: mc}

	results := make(chan *Result)

	fmt.Println("Start Benchmark")
	for i := uint64(0); i < userConfig.Clients; i++ {
		go a.Attack(i, results, userConfig)
	}

	QPS := 0
	go func() {
		for {
			select {
			case <- results:
			 	QPS = QPS + 11
			}
		}
	}()

	for {
		select {
		case <- time.After(time.Second):
			fmt.Printf("Time=%s ,QPS=%d\r\n", 
				time.Now().Format("2006-01-02 15:04:05"), QPS)
				QPS = 0
			}
	}

	fmt.Println("End Benchmark")
}