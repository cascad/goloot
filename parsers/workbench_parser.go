package parsers

func WorkbenchParser(raw map[string]interface{}) (map[string]int, error) {
	flat := make(map[string]int)
	var wbs map[string]interface{}

	if rawWbs, ok := raw["building"].(map[string]interface{})["workbench"]; ok {
		wbs = rawWbs.(map[string]interface{})
	} else {
		return flat, nil
	}

	for _, invType := range []string{"input_inventories", "output_inventories", "fuel_inventories", "feed_inventories", "inventories"} {

		var rawInfo map[string]interface{}
		if raw, ok := wbs[invType]; ok {
			rawInfo = raw.(map[string]interface{})
			res, err := ParseInventories(rawInfo)
			PanicOnErr(err)
			flat = ReduceFlats(flat, res)
		}
	}

	return flat, nil

}
