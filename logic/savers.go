package logic

import (
	. "github.com/cascad/goloot/data_structs"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
)

func Saver(callback func(string, map[string]int, *Aggregate) error, chResults <-chan Result, bads *BadUids, aggr *Aggregate) {
	for res := range chResults {
		aggr.AddRequest(1)
		aggr.AddProcess(1)

		if res.Err == nil {
			err := callback(res.Uid, res.Data, aggr)

			if err != nil {
				aggr.AddFailed(1)
				bads.Add(res.Uid)
				log.Println(fmt.Sprintf("Problem with storage: %v", err))
			} else {
				aggr.AddSuccess(1)
			}
			//results[res.Uid] = res.Data

		} else {
			//log.Println("saver: bad uid", res.Uid, res.Chan)
			bads.Add(res.Uid)
			aggr.AddFailed(1)
		}

	}
}

func MongoAggCallback(uid string, data map[string]int, agg *Aggregate) error {
	bsonV := bson.M{}

	for k, v := range data {
		bsonV[k] = v
	}
	bsonV["uid"] = uid

	err := agg.Db.Insert(bsonV)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func RedisAggCallback(uid string, data map[string]int, agg *Aggregate) error {
	for k, v := range data {
		_, perr := agg.Redis.RPush(k, v).Result()
		_, ierr := agg.Redis.Incr(fmt.Sprintf("stcounter_%s", k)).Result()
		if perr != nil || ierr != nil {
			return fmt.Errorf("RedisAggCallback: %v != %v", perr, ierr)
		}
	}
	_, err := agg.Redis.Incr("real_saves").Result()
	if err != nil {
		return fmt.Errorf("RedisAggCallback: %v", err)
	}
	return nil
}
