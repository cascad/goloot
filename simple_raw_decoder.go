package main

import (
	. "github.com/cascad/goloot/helpers"
	"github.com/cascad/goloot/parsers"
	"fmt"
	"io/ioutil"
	"log"
)

func SimpleDecoder() map[string]int {
	uid := "123"
	rawData, err := ioutil.ReadFile("test_data/raw_test_data_1_ok.bin")

	//rawData, err := GetDataFromRiak(uid, rconn)
	if err != nil {
		log.Fatal(err)
	}

	js, err := DecodeDBData(uid, &rawData)

	if err != nil {
		log.Fatal(err)
	}

	result, err := parsers.ParseBuilds(uid, js)
	return result
}

func main_old_2() {

	r := SimpleDecoder()
	for k, v := range r {
		fmt.Printf("%s = %d\n", k, v)
	}
}
