package logic

import (
	. "github.com/cascad/goloot/data_structs"
	"fmt"
	"log"
)

func Reader(lines *[]string, chJobs chan<- string, end chan<- bool, agg *Aggregate) {
	ln := len(*lines)
	for i, uid := range *lines {
		c := i + 1
		if (c)%1000 == 0 {
			ar := agg.Requested()
			ap := agg.Processed()
			as := agg.Success()
			af := agg.Failed()
			ac := agg.Corrupted()
			log.Println(fmt.Sprintf("%.5f%% | %d/%d | [%d %d %d %d %d]", float64(c)/float64(ln)*100, c, ln, ar, ap, as, af, ac))
		}
		chJobs <- uid
	}
	end <- true
}
