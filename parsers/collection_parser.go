package parsers

func CollectionParser(raw map[string]interface{}) (map[string]int, error) {
	flat := make(map[string]int)
	var collection map[string]interface{}

	if rawInventories, ok := raw["building"].(map[string]interface{})["collection"]; ok {
		collection = rawInventories.(map[string]interface{})
		if completed, ok := collection["completed"].(bool); ok {
			if completed == false {
				for k, v := range collection["collections"].(map[string]interface{}) {
					flat[k] = int(v.(float64))
				}
				return flat, nil

			} else {
				return flat, nil
			}
		}
	}
	return flat, nil
}
