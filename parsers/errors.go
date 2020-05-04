package parsers

import "log"

func PanicOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func PanicOnNok(ok bool) {
	if !ok {
		log.Fatal()
	}
}
