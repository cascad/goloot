package data_structs

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"gopkg.in/mgo.v2"
	"github.com/zegl/goriak/v3"
	"sync"
	"sync/atomic"
)

type Aggregate struct {
	Count     uint64
	Bads      BadUids
	Db        *mgo.Collection
	Riak      *goriak.Session
	Redis     *redis.Client
	Wg        sync.WaitGroup
	requested uint64 // request from riak
	processed uint64 // non-empty save
	success   uint64 // success parse
	failed    uint64 // fail parse
	corrupted uint64 // fail parse in code
	// requested = failed + success
}

func (c *Aggregate) AddCount(v int) uint64 {
	return atomic.AddUint64(&c.Count, uint64(v))
}
func (c *Aggregate) AddProcess(v uint64) uint64 {
	return atomic.AddUint64(&c.processed, v)
}
func (c *Aggregate) AddRequest(v uint64) uint64 {
	return atomic.AddUint64(&c.requested, v)
}
func (c *Aggregate) AddSuccess(v uint64) uint64 {
	return atomic.AddUint64(&c.success, v)
}
func (c *Aggregate) AddFailed(v uint64) uint64 {
	return atomic.AddUint64(&c.failed, v)
}
func (c *Aggregate) AddCorrupted(v uint64) uint64 {
	return atomic.AddUint64(&c.corrupted, v)
}

func (c *Aggregate) Processed() uint64 {
	return atomic.LoadUint64(&c.processed)
}
func (c *Aggregate) Requested() uint64 {
	return atomic.LoadUint64(&c.requested)
}
func (c *Aggregate) Success() uint64 {
	return atomic.LoadUint64(&c.success)
}
func (c *Aggregate) Failed() uint64 {
	return atomic.LoadUint64(&c.failed)
}
func (c *Aggregate) Corrupted() uint64 {
	return atomic.LoadUint64(&c.corrupted)
}
func (a *Aggregate) Check() bool {
	r := a.Requested()
	//p := a.Processed()
	return r == a.Success()+a.Failed()
}

type Result struct {
	Uid  string
	Data map[string]int
	Err  error
	Chan int
}

type BadUids struct {
	bads []string
}

func (c *BadUids) Add(uid string) {
	c.bads = append(c.bads, uid)
}

type Save struct {
	Player   json.RawMessage     `json:"omitempty"`
	Location map[string]Location `json:"omitempty"`
}

type Location struct {
	LootObjects map[string]interface{}
	Builder     map[string]interface{}
}

type SyncSave struct {
	Raw map[string]map[string]int `msgpack:"resources"`
}

func (s *SyncSave) Value() (int, error) {
	if c, ok := s.Raw["coins"]; !ok {
		return 0, fmt.Errorf("missing key coins")
	} else {
		if a, ok := c["amount"]; !ok {
			return 0, fmt.Errorf("missing key amount")
		} else {
			return a, nil
		}
	}
}
