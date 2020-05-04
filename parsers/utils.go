package parsers

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

func ReduceFlats(main map[string]int, other map[string]int) map[string]int {
	newFlat := make(map[string]int)

	for k, v := range other {
		if val, ok := main[k]; !ok {
			newFlat[k] = v
		} else {
			newFlat[k] = val + v
		}
	}
	for k, v := range main {
		if _, ok := newFlat[k]; !ok {
			newFlat[k] = v
		}
	}

	return newFlat
}
func Pprint(obj map[string]int) {
	strs := []string{}
	for k := range obj {
		strs = append(strs, k)
	}
	sort.Strings(strs)

	var lines []string
	for _, k := range strs {
		v := obj[k]
		lines = append(lines, fmt.Sprintf("%s = %v", k, v))
	}
	log.Println(fmt.Sprintf("\n%s\n", strings.Join(lines, "\n")))
}
