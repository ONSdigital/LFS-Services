package eurostat

import (
	"fmt"
	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
	"log"
	"testing"
)

func TestOtto(t *testing.T) {
	vm := otto.New()

	script, err := vm.Compile("ECOURF16.js", nil)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot compile script %s", err))
	}
	_, err = vm.Run(script)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot run script %s", err))
	}

	if value, err := vm.Call("ECOURF16", nil, 3, 0, 0); err == nil {
		fmt.Printf("Response: %s", value)
	}

}
