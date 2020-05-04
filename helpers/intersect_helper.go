package helpers

func HashString(a []string, b []string) []string {
	set := make([]string, 0)
	hash := make(map[string]bool)

	for _, e := range a {
		hash[e] = true
	}

	for _, e := range b {
		if _, found := hash[e]; found {
			set = append(set, e)
		}
	}

	return set
}

func HashNotInArray(a []string, b []string) []string {
	set := make([]string, 0)
	hash := make(map[string]bool)

	for _, e := range a {
		hash[e] = true
	}

	for _, e := range b {
		if _, found := hash[e]; !found {
			set = append(set, e)
		}
	}

	return set
}
