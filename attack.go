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