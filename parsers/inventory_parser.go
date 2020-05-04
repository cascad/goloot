package parsers

func ParseInventories(raw map[string]interface{}) (map[string]int, error) {
	flat := make(map[string]int)

	for _, v := range raw {
		invItem, ok := v.(map[string]interface{})
		PanicOnNok(ok)
		//inventoryId := inv_item["id"].(string)
		inventory, ok := invItem["inventory"].(map[string]interface{})
		PanicOnNok(ok)

		//inventoryId := inv_item["id"].(string)
		if len(inventory) == 0 {
			continue
		}
		cells := inventory["cells"].(map[string]interface{})
		for _, rawCell := range cells {
			cell := rawCell.(map[string]interface{})
			if len(cell) == 0 {
				continue
			}
			stackId := cell["stack_id"].(string)
			rawStack, ok := cell["stack"]
			var amount int

			if !ok {
				amount = 1

			} else {
				stack := rawStack.(map[string]interface{})
				rawAmount, ok := stack["amount"]

				if ok {
					amount = int(rawAmount.(float64))
				} else {
					amount = 1
				}
			}
			if val, ok := flat[stackId]; !ok {
				flat[stackId] = amount
			} else {
				flat[stackId] = val + amount
			}

		}
	}
	return flat, nil

}

//func test_inv_parser() {
//	rawData := map[string]interface{"0": {"id": "backpack_8", "inventory": {"cells": {"0": {"stack_id": "wood", "stack": {"amount": 20}}, "1": {"stack_id": "wood", "stack": {"amount": 20}}, "2": {"stack_id": "furniture_couch", "stack": {}},"3": {"stack_id": "furniture_couch","stack": {}},"4": {"stack_id": "furniture_couch","stack": {}},"5": {"stack_id": "furniture_couch","stack": {}},"6": {"stack_id": "furniture_couch","stack": {}},"7": {"stack_id": "furniture_couch","stack": {}}}}}}
//	parsedData, err := ParseInventories(rawData.(map[string]interface{}))
//	log.Fatal(err != nil)
//
//	checkData := map[string]float64{
//		"furniture_couch": 6, "wood": 40}
//
//	for k, v := range parsedData {
//		if v != checkData[k] {
//			log.Fatal(v != checkData[k])
//		}
//	}
//
//	log.Println(ParseInventories(&rawData))
//}
