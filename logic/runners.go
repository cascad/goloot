package logic

import (
	. "github.com/cascad/goloot/data_structs"
	"log"
	"time"
)

func ProdRunner(numWorkers int, numJobs int, numResults int, lines *[]string, agg *Aggregate) {
	results := make(chan Result, numResults)
	end := make(chan bool)
	jobs := make(chan string, numJobs)

	for w := 1; w <= numWorkers; w++ {
		go Worker(w, jobs, results, agg.Riak, agg)
	}
	go Reader(lines, jobs, end, agg)
	go Saver(RedisAggCallback, results, &agg.Bads, agg)

	ticker := time.NewTicker(1000 * time.Millisecond)
	forTimeout := time.Now()
	log.Println("started..")
	<-end

	for {
		select {
		case <-ticker.C:
			if agg.Check() == true {
				break
			}
			if time.Now().Sub(forTimeout).Seconds() > 5 {
				break
			}
		}
		break
	}
	//agg.Wg.Wait()
}

func DevRunner(numWorkers int, numJobs int, numResults int, lines *[]string, agg *Aggregate) {
	results := make(chan Result, numResults)
	end := make(chan bool)
	jobs := make(chan string, numJobs)

	for w := 1; w <= numWorkers; w++ {
		go TestRawWorker(w, jobs, results, nil, agg)
	}
	go Reader(lines, jobs, end, agg)
	go Saver(RedisAggCallback, results, &agg.Bads, agg)

	ticker := time.NewTicker(1000 * time.Millisecond)
	forTimeout := time.Now()

	log.Println("started..")
	<-end

	for {
		select {
		case <-ticker.C:
			log.Println(agg.Check())
			if agg.Check() == true {
				break
			}
			if time.Now().Sub(forTimeout).Seconds() > 5 {
				break
			}
		}
		break
	}
	//agg.Wg.Wait()
}

func DevSyncRunner(numWorkers int, numJobs int, numResults int, lines *[]string, agg *Aggregate) {
	results := make(chan Result, numResults)
	end := make(chan bool)
	jobs := make(chan string, numJobs)

	for w := 1; w <= numWorkers; w++ {
		go TestSyncWorker(w, jobs, results, agg.Riak, agg)
	}
	go Reader(lines, jobs, end, agg)
	go Saver(RedisAggCallback, results, &agg.Bads, agg)

	ticker := time.NewTicker(1000 * time.Millisecond)
	forTimeout := time.Now()
	log.Println("started..")
	<-end

	for {
		select {
		case <-ticker.C:
			if agg.Check() == true {
				break
			}
			if time.Now().Sub(forTimeout).Seconds() > 5 {
				break
			}
		}
		break
	}
	//agg.Wg.Wait()
}
func SyncRunner(numWorkers int, numJobs int, numResults int, lines *[]string, agg *Aggregate) {
	results := make(chan Result, numResults)
	end := make(chan bool)
	jobs := make(chan string, numJobs)

	for w := 1; w <= numWorkers; w++ {
		go SyncWorker(w, jobs, results, agg.Riak, agg)
	}
	go Reader(lines, jobs, end, agg)
	go Saver(RedisAggCallback, results, &agg.Bads, agg)

	ticker := time.NewTicker(1000 * time.Millisecond)
	forTimeout := time.Now()
	log.Println("started..")
	<-end

	for {
		select {
		case <-ticker.C:
			if agg.Check() == true {
				break
			}
			if time.Now().Sub(forTimeout).Seconds() > 5 {
				break
			}
		}
		break
	}
	//agg.Wg.Wait()
}
