package main

import (
	"github.com/cascad/goloot/parsers"
	"encoding/json"
	"io/ioutil"
	"log"
)

func main_() {
	//fn:= os.Args[1]
	//fn := "test_data/test_data_2.json"
	fn := "test_data/test_data_2.json"
	uid := "1"

	b, err := ioutil.ReadFile(fn)
	parsers.PanicOnErr(err)
	//result, err := parsers.ParseSave(uid, &b, false)
	result, err := parsers.ParseBuilds(uid, &b)
	parsers.PanicOnErr(err)
	parsers.Pprint(result)
}

func main2() {
	//fn:= os.Args[1]
	//fn := "test_data/test_data_2.json"
	fn := "test_data/test_builds.json"

	b, err := ioutil.ReadFile(fn)
	parsers.PanicOnErr(err)

	var s map[string]interface{}
	err = json.Unmarshal(b, &s)
	if err != nil {
		log.Fatal(err)
	}
	result, err := parsers.BuilderValueParser(s)
	parsers.PanicOnErr(err)
	parsers.Pprint(result)
}
