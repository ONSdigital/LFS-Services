package r_times_test

import (
	"fmt"
	"log"
	"pds-go/lfs/r"
	times "pds-go/lfs/r/r_times"
	"testing"
)

/*
Set R_HOME to the correct directory in the environment, e.g.:
R_HOME=/Library/Frameworks/R.framework/Resources
*/

func Test_times(t *testing.T) {

	t.Logf("Starting test - times")

	i := r.RFunctions{}
	defer i.Free()

	i.LoadRSource("times.R")

	res, err := times.Times(5.78, 9.23)

	if err != nil {
		log.Printf("Call to R failed: %s", err)
		panic(err)
	}

	fmt.Printf("Result: %f\n", res)
	t.Logf("Test - add, successful")
}
