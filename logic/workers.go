package logic

import (
	. "github.com/cascad/goloot/data_structs"
	. "github.com/cascad/goloot/helpers"
	"github.com/cascad/goloot/parsers"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"github.com/zegl/goriak/v3"
	"io/ioutil"
	"log"
	"os"
)

func Worker(id int, jobs <-chan string, results chan<- Result, rconn *goriak.Session, agg *Aggregate) {
	bucket := "blob"
	for uid := range jobs {
		//fmt.Println("worker:", id, "-", uid)

		rawData, err := GetDataFromRiak(uid, rconn, bucket)
		if err != nil {
			results <- Result{Uid: uid, Err: err, Chan: id}
			continue
		}

		js, err := DecodeDBData(uid, rawData)

		if err != nil {
			results <- Result{Uid: uid, Err: err, Chan: id}
			continue
		}

		result, err := parsers.ParseSave(uid, js, false)
		results <- Result{Uid: uid, Data: result, Err: err, Chan: id}
	}
}

func TestSyncWorker(id int, jobs <-chan string, results chan<- Result, rconn *goriak.Session, agg *Aggregate) {
	var filename string
	for uid := range jobs {
		//fmt.Println("worker:", id, "-", uid)

		if uid == "123" {
			filename = "test_data/raw_test_profile_data_1.bin"
		} else {
			filename = "test_data/raw_test_profile_data_2.bin"
			//filename = "test_data/raw_test_data_4_corrupted.bin"
		}

		if _, err := os.Stat(filename); os.IsNotExist(err) {
			log.Println(fmt.Sprintf("TestSyncWorker: %v", err))
			continue
		}

		rawData, err := ioutil.ReadFile(filename)
		if err != nil {
			results <- Result{Uid: uid, Err: err, Chan: id}
			continue
		}
		rawResult, err := GetCoinsFromSyncSave(&rawData)
		if err != nil {
			agg.AddCorrupted(1)
			agg.Bads.Add(uid)
			results <- Result{Uid: uid, Err: err, Chan: id}
			continue
		}
		result := map[string]int{"coins": rawResult}

		results <- Result{Uid: uid, Data: result, Err: err, Chan: id}
	}
}

func SyncWorker(id int, jobs <-chan string, results chan<- Result, rconn *goriak.Session, agg *Aggregate) {
	bucket := "profile"
	for uid := range jobs {
		//fmt.Println("worker:", id, "-", uid)

		rawData, err := GetDataFromRiak(uid, rconn, bucket)
		if err != nil {
			results <- Result{Uid: uid, Err: err, Chan: id}
			continue
		}

		rawResult, err := GetCoinsFromSyncSave(rawData)
		if err != nil {
			agg.AddCorrupted(1)
			agg.Bads.Add(uid)
			results <- Result{Uid: uid, Err: err, Chan: id}
			continue
		}
		result := map[string]int{"coins": rawResult}

		results <- Result{Uid: uid, Data: result, Err: err, Chan: id}
	}
}

func TestRawWorker(id int, jobs <-chan string, results chan<- Result, rconn *goriak.Session, agg *Aggregate) {
	for uid := range jobs {
		//fmt.Println("worker:", id, "-", uid)

		var filename string
		if uid == "123" {
			filename = "test_data/raw_test_data_1_ok.bin"
		} else {
			filename = "test_data/raw_test_data_2_corrupted.bin"
		}

		//filename := "test_data/raw_test_data_2_corrupted.bin"
		//filename := "test_data/raw_test_data_3_ok.bin"
		//filename := "test_data/raw_test_data_4_corrupted.bin"
		//filename := "test_data/raw_test_data_5_empty.bin"

		if _, err := os.Stat(filename); os.IsNotExist(err) {
			log.Println(err)
		}

		rawData, err := ioutil.ReadFile(filename)

		if err != nil {
			results <- Result{Uid: uid, Err: err, Chan: id}
			continue
		}

		js, err := DecodeDBData(uid, &rawData)

		if err != nil {
			results <- Result{Uid: uid, Err: err, Chan: id}
			continue
		}

		result, err := parsers.ParseBuilds(uid, js)
		//parsers.PanicOnErr(err)
		results <- Result{Uid: uid, Data: result, Err: err, Chan: id}
	}
}

func BuildersWorker(id int, jobs <-chan string, results chan<- Result, rconn *goriak.Session, agg *Aggregate) {
	bucket := "blob"
	for uid := range jobs {
		//fmt.Println("worker:", id, "-", uid)

		//var fCheck bson.M
		//err := agg.Db.Find(bson.M{"uid": uid}).One(&fCheck)
		//if err == nil {
		//	//log.Println("uid:", uid, "already exists in Db")
		//	agg.AddProcess(1)
		//	agg.AddSuccess(1)
		//	continue
		//}

		rawData, err := GetDataFromRiak(uid, rconn, bucket)
		if err != nil {
			results <- Result{Uid: uid, Err: err, Chan: id}
		}

		js, err := DecodeDBData(uid, rawData)

		if err != nil {
			//log.Println("worker: ", id, ", bad uid", uid)
			results <- Result{Uid: uid, Err: err, Chan: id}
		}

		result, err := parsers.ParseBuilds(uid, js)
		//parsers.PanicOnErr(err)
		results <- Result{Uid: uid, Data: result, Err: err, Chan: id}
	}
}

func TestMongoWorker(id int, jobs <-chan string, results chan<- Result, rconn *goriak.Session, agg *Aggregate) {
	for uid := range jobs {
		fmt.Println("worker:", id, "-", uid)

		var fCheck bson.M
		err := agg.Db.Find(bson.M{"uid": uid}).One(&fCheck)
		if err == nil {
			log.Println("uid:", uid, "already exists in Db")
			agg.AddProcess(1)
			agg.AddSuccess(1)
			continue
		}

		//log.Println(uid, "loading..")
		b, err := ioutil.ReadFile(uid + ".json")
		parsers.PanicOnErr(err)

		result, err := parsers.ParseSave(uid, &b, false)
		parsers.PanicOnErr(err)

		//log.Println(uid, "sending..")
		results <- Result{Uid: uid, Data: result, Err: err, Chan: id}
		agg.AddProcess(1)
		//log.Println(uid, "sended")
	}
}
