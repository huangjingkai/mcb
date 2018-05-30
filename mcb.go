package main

import (
	"fmt"
	"time"
	"strconv"

	//"strings"
	
	"github.com/pangudashu/memcache"
	flag "github.com/spf13/pflag"
)

type UserConfig struct {
	IPs       []string
	Clients   int
	Requests  int64
	DataSize  int
}

var userConfig UserConfig

func init() {
	flag.IntVar(&userConfig.Clients, "clients", 50, "Number of parallel connections (default 50)")
	flag.Int64Var(&userConfig.Requests, "requests", 100000, "Total number of requests (default 100000)")
	flag.IntVar(&userConfig.DataSize, "dataSize", 2, "Data size of SET/GET value in bytes (default 2)")
	flag.StringSliceVar(&userConfig.IPs, "ips", []string{"127.0.0.1:11211"}, "Server hostname (default 127.0.0.1:11211), or can set 127.0.0.1:11211,127.0.0.2:11211")
}

func main() {
    flag.Parse()
    fmt.Println(userConfig)

    mcs := []*memcache.Server{}

    for _, ip := range(userConfig.IPs) {
    	mcs = append(mcs, &memcache.Server{
    		Address: ip,
    		Weight: 50,
    		MaxConn: userConfig.Clients,
    		InitConn: userConfig.Clients / 2,
    		IdleTime: time.Hour * 2})
    }

	mc, err := memcache.NewMemcache(mcs)
	if err != nil {
		fmt.Println(err)
		return
	}

	//设置是否自动剔除无法连接的server，默认不开启(建议开启)
	//如果开启此选项被踢除的server如果恢复正常将会再次被加入server列表
	mc.SetRemoveBadServer(true)

	beforeTimeMS := time.Now().UnixNano()

	for i := 0; i < int(userConfig.Requests); i++ {
		res, err := mc.Set("test" + strconv.Itoa(i), 0, 300)
		if err != nil {
			fmt.Println("mc.Set error, ", res, err)
		}
	}
	afterTimeMS := time.Now().UnixNano()

	QPS := userConfig.Requests * 1e9  / (afterTimeMS - beforeTimeMS)

    fmt.Println(time.Unix(beforeTimeMS/1e9, 0).String()) //输出当前英文时间戳格式  
    fmt.Println(time.Unix(afterTimeMS/1e9, 0).String()) //输出当前英文时间戳格式  

	fmt.Printf("QPS: %d\r\nuserConfig.Requests=%d\r\nafterTime=%s %d\r\nbeforeTime=%s %d\r\n", 
		QPS, 
		userConfig.Requests, 
		time.Unix(beforeTimeMS/1e9, 0).String(), 
		afterTimeMS, time.Unix(afterTimeMS/1e9, 0).String(), 
		beforeTimeMS)
}