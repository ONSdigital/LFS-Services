package eurostat

import (
	"fmt"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
	"testing"
)

func TestECOURF16(t *testing.T) {

	resolve.AllowFloat = true
	resolve.AllowLambda = true

	thread := &starlark.Thread{Name: "my thread"}
	globals, err := starlark.ExecFile(thread, "ECOURF16.star", nil, nil)

	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}

	ECOURF16 := globals["ECOURF16"]

	v, err := starlark.Call(thread, ECOURF16, starlark.Tuple{
		starlark.MakeInt(3),
		starlark.MakeInt(3),
		starlark.MakeInt(0),
		starlark.MakeInt(52),
	}, nil)

	if err != nil {
		t.Error("Received error from call to ECOURF16")
		t.FailNow()
	}

	a := v.(starlark.Tuple).Index(1)
	fmt.Printf("Response = %d\n", a)

	fmt.Println("Done")
}
