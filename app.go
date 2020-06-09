package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"sync"

	ds "github.com/cascad/goloot/data_structs"
	"github.com/cascad/goloot/erlang"
	"github.com/cascad/goloot/helpers"
	"github.com/cascad/goloot/logic"
	"github.com/cascad/goloot/parsers"
	"github.com/go-redis/redis"
)

func Dev1(runner func(int, int, int, *[]string, *ds.Aggregate)) {
	var agg ds.Aggregate
	agg.Wg = sync.WaitGroup{}

	numProc := flag.Int("p", 1, "core num, default=1")
	numWorkers := flag.Int("w", 1, "num workers, default=1")
	flush := flag.Int("c", 1, "clear storage, default=False")
	numDB := flag.Int("db", 0, "num of db, default=0")
	readerBuffers := flag.Int("rb", 0, "reader buffer size, default=0")
	saverBuffers := flag.Int("sb", 0, "saver buffer size, default=0")
	filename := flag.String("f", "test_data/test_file_uids.txt", "file with uids")
	flag.Parse()

	runtime.GOMAXPROCS(*numProc)

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",     // no password set
		DB:       *numDB, // use default DB
	})

	defer client.Close()

	if *flush == 1 {
		log.Println("Flushing DB")
		client.FlushAll()
	}

	agg.Redis = client

	lines := helpers.ReadUidsFile(*filename)
	counter := len(lines)
	agg.AddCount(counter)

	counter, err := agg.Redis.Get("counter").Int()
	if err.Error() == "redis: nil" || counter == 0 {
		agg.Redis.Set("counter", counter, 0)
	}

	completedCounter := "real_saves"
	completed, err := agg.Redis.Get(completedCounter).Int()
	if err != nil {
		if err.Error() != "redis: nil" {
			log.Fatal(fmt.Sprintf("bad get %s -> %v", completedCounter, err))
		}
	}

	needUids := lines[completed:]
	toProcess := len(needUids)

	log.Println(fmt.Sprintf("%d/%d/%d", counter, completed, toProcess))
	log.Println(fmt.Sprintf("uids(%d)[%d:]", counter, completed))

	//logic.DevRunner(numWorkers, readerBuffers, saverBuffers, &needUids, &agg)
	runner(*numWorkers, *readerBuffers, *saverBuffers, &needUids, &agg)

	log.Println(agg.Check(), agg.Requested(), agg.Processed(), agg.Success(), agg.Failed())
}

func Prod(runner func(int, int, int, *[]string, *ds.Aggregate)) {
	// go run app.go -f utest.txt -c 1
	var agg ds.Aggregate
	agg.Wg = sync.WaitGroup{}

	numProc := flag.Int("p", 1, "core num, default=1")
	numWorkers := flag.Int("w", 1, "num workers, default=1")
	flush := flag.Int("c", 1, "clear storage, default=False")
	numDB := flag.Int("db", 0, "num of db, default=0")
	readerBuffers := flag.Int("rb", 0, "reader buffer size, default=0")
	saverBuffers := flag.Int("sb", 0, "saver buffer size, default=0")
	filename := flag.String("f", "", "file with uids")
	flag.Parse()

	runtime.GOMAXPROCS(*numProc)

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",     // no password set
		DB:       *numDB, // use default DB
	})
	defer client.Close()
	if *flush == 1 {
		log.Println("Flushing DB")
		client.FlushAll()
	}
	agg.Redis = client

	riakHost := "ldoe-backup-Db-n01"
	riakSession, err := erlang.Connect(riakHost)
	parsers.PanicOnErr(err)
	agg.Riak = riakSession

	lines := helpers.ReadUidsFile(*filename)
	counter := len(lines)
	agg.AddCount(counter)

	counter, err = agg.Redis.Get("counter").Int()
	if err.Error() == "redis: nil" || counter == 0 {
		agg.Redis.Set("counter", counter, 0)
	}

	completedCounter := "real_saves"
	completed, err := agg.Redis.Get(completedCounter).Int()
	if err != nil {
		if err.Error() != "redis: nil" {
			log.Fatal(fmt.Sprintf("bad get %s -> %v", completedCounter, err))
		}
	}

	needUids := lines[completed:]
	toProcess := len(needUids)

	log.Println(fmt.Sprintf("%d/%d/%d", counter, completed, toProcess))
	log.Println(fmt.Sprintf("uids(%d)[%d:]", counter, completed))

	runner(*numWorkers, *readerBuffers, *saverBuffers, &needUids, &agg)
	log.Println(agg.Check(), agg.Requested(), agg.Processed(), agg.Success(), agg.Failed())
}

func main() {
	//Prod(logicProdRunner)
	//Dev1(logic.DevRunner)
	// Prod(logic.SyncRunner)
	Dev1(logic.DevSyncRunner)
}
