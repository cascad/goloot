package main

import (
	"github.com/cascad/goloot/helpers"
	"flag"
	"fmt"
	"github.com/zegl/goriak/v3"
	"io/ioutil"
	"log"
)

func main_old_1() {
	// go run riak_gun.go -u cbedf6e875a94be38138eb421df9f143
	var uid, bucket string
	flag.StringVar(&uid, "u", "", "user id")
	flag.StringVar(&bucket, "b", "profile", "user id")
	flag.Parse()
	// profile, blob

	riakHost := "ldoe-backup-Db-n01"
	session, err := goriak.Connect(goriak.ConnectOpts{
		Address: riakHost,
		//Port:    8087,
	})
	if err != nil {
		panic(err)
	}

	log.Println(fmt.Sprintf("uid: %s", uid))
	log.Println(fmt.Sprintf("bucket: %s", bucket))
	rawData, err := helpers.GetDataFromRiak(uid, session, bucket)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("msg.bin", *rawData, 0644)
	//js, err := helpers.DecodeDBData(uid, rawData)

	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(js)
}
