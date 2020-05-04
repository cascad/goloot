package parsers

func ChestParser(raw map[string]interface{}) (map[string]int, error) {
	flat := make(map[string]int)
	var inventories map[string]interface{}

	if rawInventories, ok := raw["building"].(map[string]interface{})["inventories"]; ok {
		inventories = rawInventories.(map[string]interface{})
	} else {
		return flat, nil
	}
	parsed, err := ParseInventories(inventories)
	PanicOnErr(err)
	return parsed, err
}
