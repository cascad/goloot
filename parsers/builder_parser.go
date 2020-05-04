package parsers

func BuilderParser(raw map[string]interface{}, checkUnlocked bool, withCollections bool) (map[string]int, error) {
	flat := make(map[string]int)
	var furnitures map[string]interface{}

	if rawFurnitures, ok := raw["furnitures"].(map[string]interface{})["items"]; ok {
		furnitures = rawFurnitures.(map[string]interface{})
	} else {
		return flat, nil
	}

	for _, v := range furnitures {
		innerData, ok := v.(map[string]interface{})
		PanicOnNok(ok)
		wbsFlat, err := WorkbenchParser(innerData)
		PanicOnErr(err)
		flat = ReduceFlats(flat, wbsFlat)

		if rawUnlocked, ok := innerData["building"].(map[string]interface{})["unlocked"]; ok {
			unlocked := rawUnlocked.(bool)
			if !(checkUnlocked && !unlocked) {
				chestFlat, err := ChestParser(innerData)
				PanicOnErr(err)
				flat = ReduceFlats(flat, chestFlat)
			}
		}
		if withCollections {
			probes, err := CollectionParser(innerData)
			PanicOnErr(err)
			flat = ReduceFlats(flat, probes)
		}

	}

	return flat, nil

}
