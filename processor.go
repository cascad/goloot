package main

import (
	"encoding/csv"
	"fmt"
	"github.com/cascad/goloot/parsers"
	"github.com/cascad/goloot/stats_local"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"strconv"
)

func MaxIntSlice(v []float64) (m float64) {
	if len(v) > 0 {
		m = v[0]
	}
	for i := 1; i < len(v); i++ {
		if v[i] > m {
			m = v[i]
		}
	}
	return
}
func MinIntSlice(v []float64) (m float64) {
	if len(v) > 0 {
		m = v[0]
	}
	for i := 1; i < len(v); i++ {
		if v[i] < m {
			m = v[i]
		}
	}
	return
}

func getRand(v int) (res []float64) {
	for i := 0; i < v; i++ {
		res = append(res, float64(rand.Intn(20)))
	}
	return
}

func f2s(input_num float64) string {

	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', -1, 64)
}

func bToMb(b uint64) float64 {
	return float64(b) / 1024 / 1024
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}
func CalcSave(item string, v *[]float64, filename string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	parsers.PanicOnErr(err)
	defer f.Close()
	writer := csv.NewWriter(f)
	defer writer.Flush()

	val := stats_local.Sample{Xs: *v}
	r := make([]string, 0)
	c := len(val.Xs)

	r = append(r, item)

	r = append(r, strconv.Itoa(c))

	min := val.Percentile(0)
	r = append(r, f2s(min))

	p1 := val.Percentile(0.01)
	r = append(r, f2s(p1))

	p5 := val.Percentile(0.05)
	r = append(r, f2s(p5))

	p10 := val.Percentile(0.1)
	r = append(r, f2s(p10))

	p20 := val.Percentile(0.2)
	r = append(r, f2s(p20))

	p30 := val.Percentile(0.3)
	r = append(r, f2s(p30))

	p40 := val.Percentile(0.4)
	r = append(r, f2s(p40))

	p50 := val.Percentile(0.5)
	r = append(r, f2s(p50))

	p60 := val.Percentile(0.6)
	r = append(r, f2s(p60))

	p70 := val.Percentile(0.7)
	r = append(r, f2s(p70))

	p80 := val.Percentile(0.8)
	r = append(r, f2s(p80))

	p90 := val.Percentile(0.9)
	r = append(r, f2s(p90))

	p95 := val.Percentile(0.95)
	r = append(r, f2s(p95))

	p99 := val.Percentile(0.99)
	r = append(r, f2s(p99))

	max := val.Percentile(1)
	r = append(r, f2s(max))

	//fmt.Println("save:", r, *v)
	err = writer.Write(r)
	parsers.PanicOnErr(err)
}

//type Saved struct {
//	Count     uint64
//	Bads      BadUids
//	Db        *mgo.Collection
//	Riak      *goriak.Session
//	Wg        sync.WaitGroup
//	requested uint64
//	processed uint64
//	success   uint64
//	failed    uint64
//}
//func mmc_proc() {
//	uidsFilename := os.Args[1]
//	mmc := memcache.New("127.0.0.1")
//	uids := ReadUidsFile(uidsFilename)
//
//	for i, uid := range uids {
//		if len(uid) > 0 {
//			data, err := mmc.Get(uid)
//			if err == nil {
//
//			}
//		}
//	}
//}

func mainOld() {
	resultFilename := os.Args[1]
	collection := os.Args[2]
	//resultFilename := "result.csv"
	mongoHost := "127.0.0.1"
	// mmc := memcache.New("127.0.0.1")

	var err = os.Remove(resultFilename)

	var tmp map[string]interface{}
	//data := make(map[string][]float64, 0)

	session, err := mgo.Dial(mongoHost)
	parsers.PanicOnErr(err)
	defer session.Close()
	db := session.DB("goloot")
	c := db.C(collection)

	keys := make(map[string]bool, 0)

	iter := c.Find(bson.M{}).Iter()
	for iter.Next(&tmp) {
		for k := range tmp {
			if k != "_id" && k != "uid" {
				keys[k] = true
			}
		}
	}
	log.Println(len(keys), "items")

	tmp = nil
	for item := range keys {
		dbCount, err := c.Find(bson.M{item: bson.M{"$exists": true}}).Count()
		parsers.PanicOnErr(err)
		log.Println(fmt.Sprintf("%s(%v)", item, dbCount))
		var seq []float64

		iter := c.Find(bson.M{item: bson.M{"$exists": true}}).Iter()
		for iter.Next(&tmp) {
			rv := tmp[item]
			var val float64
			switch t := rv.(type) {
			case int, int8, int16, int32, int64:
				val = float64(reflect.ValueOf(t).Int()) // a has type int64
			case uint, uint8, uint16, uint32, uint64:
				val = float64(reflect.ValueOf(t).Uint()) // a has type uint64d
			default:
				log.Println("bad value:", item, rv)
				continue
			}
			if val > 0 {
				seq = append(seq, val)
			}
		}

		if len(seq) > 0 {
			CalcSave(item, &seq, resultFilename)
		}
	}
}
