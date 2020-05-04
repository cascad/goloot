package parsers

import (
	"strconv"
)

type Building struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Grade int `json:"grade"`
}

type Item struct {
	Id       string   `json:"description_id"`
	Building Building `json:"building"`
}

func BuilderValueParser(raw map[string]interface{}) (map[string]int, error) {
	flat := make(map[string]int)
	var foundaments map[string]interface{}
	//var walls map[string]interface{}

	if raw, ok := raw["foundaments"].(map[string]interface{})["items"]; ok {
		foundaments = raw.(map[string]interface{})

		for _, v := range foundaments {
			rawData := v.(map[string]interface{})
			id := rawData["description_id"].(string)
			grade := rawData["building"].(map[string]interface{})["grade"].(float64)

			key := "f_" + id + "_" + strconv.FormatFloat(grade, 'f', 0, 64)

			if _, ok = flat[key]; ok {
				flat[key] += 1
			} else {
				flat[key] = 1
			}
		}
	}

	if raw, ok := raw["walls"].(map[string]interface{})["items"]; ok {
		foundaments = raw.(map[string]interface{})

		for _, v := range foundaments {
			rawData := v.(map[string]interface{})
			id := rawData["description_id"].(string)
			grade := rawData["building"].(map[string]interface{})["grade"].(float64)

			key := "w_" + id + "_" + strconv.FormatFloat(grade, 'f', 0, 64)

			if _, ok = flat[key]; ok {
				flat[key] += 1
			} else {
				flat[key] = 1
			}
		}
	}

	return flat, nil

}
