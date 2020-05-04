package parsers

import (
	"encoding/json"
	"fmt"
)

func ParseSyncSave(uid string, byteVal *[]byte) (map[string]int, error) {
	fmt.Println(byteVal)
	return map[string]int{}, nil
}

func ParseSave(uid string, byteVal *[]byte, withCollections bool) (map[string]int, error) {
	var s map[string]interface{}
	err := json.Unmarshal(*byteVal, &s)
	if err != nil {
		return map[string]int{}, err
	}

	flat := make(map[string]int)
	rawLocks := s["locations"]
	if rawLocks == nil {
		return map[string]int{}, err
	}
	locs := rawLocks.(map[string]interface{})

	for locName, loc := range locs {
		checkUnlocked := locName != "home"
		lo := loc.(map[string]interface{})["loot_objects"].(map[string]interface{})
		loParsed, err := ParseLootObjects(lo, checkUnlocked)
		if err != nil {
			return flat, fmt.Errorf("bad parse loot obj -> %s", uid)
		}
		//parsers.PanicOnErr(err)
		flat = ReduceFlats(flat, loParsed)
		builder := loc.(map[string]interface{})["builder"].(map[string]interface{})
		builderParsed, err := BuilderParser(builder, checkUnlocked, withCollections)
		if err != nil {
			return flat, fmt.Errorf("bad parse builder -> %s", uid)
		}
		//parsers.PanicOnErr(err)
		flat = ReduceFlats(flat, builderParsed)
	}
	return flat, nil
}

func ParseBuilds(uid string, byteVal *[]byte) (map[string]int, error) {
	var s map[string]interface{}
	err := json.Unmarshal(*byteVal, &s)
	if err != nil {
		return map[string]int{}, err
	}

	flat := make(map[string]int)
	rawLocks := s["locations"]
	if rawLocks == nil {
		return map[string]int{}, err
	}
	locs := rawLocks.(map[string]interface{})

	var rawHome interface{}

	if data, ok := locs["home"]; !ok {
		return map[string]int{}, nil
	} else {
		rawHome = data
	}
	builder := rawHome.(map[string]interface{})["builder"].(map[string]interface{})

	builderParsed, err := BuilderValueParser(builder)
	if err != nil {
		return flat, fmt.Errorf("bad parse builds -> %s", uid)
	}

	flat = ReduceFlats(flat, builderParsed)

	return flat, nil
}
