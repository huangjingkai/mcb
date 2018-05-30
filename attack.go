package main

import (
	"strconv"
	"time"

	"github.com/pangudashu/memcache"
)

// Attacker is an attack executor which wraps an http.Client
type Attacker struct {
	mc        *memcache.Memcache
	stopch    chan struct{}
	workers   uint64
}

const (
	MEMCACHED_SUCCESS = 0
	MEMCACHED_FAILED = 1

	MEMCACHED_GET = 0
	MEMCACHED_SET = 1
)

// Attack reads its Targets from the passed Targeter and attacks them at
// the rate specified for the given duration. When the duration is zero the attack
// runs until Stop is called. Results are sent to the returned channel as soon
// as they arrive and will have their Attack field set to the given name.
func (a *Attacker) Attack(no uint64, results chan *Result, userConfig UserConfig) {

	name := "mcs" + "_" + strconv.FormatUint(no, 10) + "_"

	seq := uint64(0)
	
	if userConfig.Ratio < 0.01 {
		userConfig.Ratio = 0.01
	}

	setCount := 100
	getCount := int(100 / userConfig.Ratio)

	for {
		key := name + strconv.FormatUint(seq, 10)
		// For Set
		for j := 0; j < setCount; j++ {
			res := a.initResult(name, seq)
			res.Command = "SET"
			_, err := a.mc.Set(key, 0, 300)
			if err != nil {
				res.Code = MEMCACHED_FAILED
				res.Error = err.Error()
			}
			res.Latency = time.Since(res.Timestamp)

			results <- res	

			seq++
			key = name + strconv.FormatUint(seq, 10)
		}

		// For Get
		for j := 0; j < getCount; j++ {
			res := a.initResult(name, seq)
			res.Command = "GET"
			_, _, err := a.mc.Get(key)
			if err != nil {
				res.Code = MEMCACHED_FAILED
				res.Error = err.Error()
			}
			res.Latency = time.Since(res.Timestamp)

			results <- res	
			seq++
		}
		// END
	}
}

func (a *Attacker) initResult(name string, seq uint64) *Result {
	return &Result{
			Attack: name,
			Seq: seq,
			Timestamp: time.Now(),
			Code: MEMCACHED_SUCCESS,
	}
}