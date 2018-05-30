package main

import (
	"sort"
	"time"
)

// Result contains the results of a single Target hit.
type Result struct {
	Attack    string        `json:"attack"`
	Seq       uint64        `json:"seq"`
	Code      uint16        `json:"code"`
	Timestamp time.Time     `json:"timestamp"`
	Latency   time.Duration `json:"latency"`
	BytesOut  uint64        `json:"bytes_out"`
	BytesIn   uint64        `json:"bytes_in"`
	Error     string        `json:"error"`
}

// End returns the time at which a Result ended.
func (r *Result) End() time.Time { return r.Timestamp.Add(r.Latency) }

// Equal returns true if the given Result is equal to the receiver.
func (r Result) Equal(other Result) bool {
	return r.Attack == other.Attack &&
		r.Seq == other.Seq &&
		r.Code == other.Code &&
		r.Timestamp.Equal(other.Timestamp) &&
		r.Latency == other.Latency &&
		r.BytesIn == other.BytesIn &&
		r.BytesOut == other.BytesOut &&
		r.Error == other.Error
}

// Results is a slice of Result type elements.
type Results []Result

// Add implements the Add method of the Report interface by appending the given
// Result to the slice.
func (rs *Results) Add(r *Result) { *rs = append(*rs, *r) }

// Close implements the Close method of the Report interface by sorting the
// Results.
func (rs *Results) Close() { sort.Sort(rs) }

// The following methods implement sort.Interface
func (rs Results) Len() int           { return len(rs) }
func (rs Results) Less(i, j int) bool { return rs[i].Timestamp.Before(rs[j].Timestamp) }
func (rs Results) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }