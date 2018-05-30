package main

import (
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
)

// Attack reads its Targets from the passed Targeter and attacks them at
// the rate specified for the given duration. When the duration is zero the attack
// runs until Stop is called. Results are sent to the returned channel as soon
// as they arrive and will have their Attack field set to the given name.
func (a *Attacker) Attack(seq uint64, results chan *Result, userConfig UserConfig) {

	name := "mcs" + "_" + strconv.FormatUint(seq, 10) + "_"

	getCount := 0
	seq := uint64(0)
	
	if userConfig.Ratio < 0.01 {
		userConfig.Ratio = 0.01
	} else {
		setCount = 100
		getCount = 100 / userConfig.Ratio
	}


	for {
		// For Set
		for j := 0; j < setCount; j++ {
			res := initResult(name, seq)
			key := name + strconv.FormatUint(seq, 10)
			ok, err := a.mc.Set(key, 0, 300)
			if err != nil {
				res.Code = MEMCACHED_FAILED
				res.Error = err.Error()
			}
			res.Latency = time.Since(res.Timestamp)

			results <- &res	
			seq++
		}

		// For Get
		for j := 0; j < getCount; j++ {
			res := initResult(name, seq)
			ok, err := a.mc.Get(key)
			results <- &res	
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