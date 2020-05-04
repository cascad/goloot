package parsers

import (
	"log"
)

func ParseLootObjects(raw map[string]interface{}, checkUnlocked bool) (map[string]int, error) {
	flat := make(map[string]int)

	for _, rawItem := range raw["items"].(map[string]interface{}) {
		item := rawItem.(map[string]interface{})["item"].(map[string]interface{})
		unlocked := item["unlocked"].(bool)

		if !(checkUnlocked && !unlocked) {
			inventories, ok := item["inventories"].(map[string]interface{})
			PanicOnNok(ok)
			result, err := ParseInventories(inventories)

			if err != nil {
				log.Fatal("pi", err)
			}

			flat = ReduceFlats(flat, result)
		}
	}
	return flat, nil
}
