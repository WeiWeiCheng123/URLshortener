package function

import (
	"fmt"
	"testing"
)

func Test_id(t *testing.T) {

	for i := 1; i < 10; i++ {
		i := Id()
		fmt.Println(i)
	}

	t.Error("123")

}

//c85ppbetmvehb53ava10
//c85ppgetmvejv84l3dp0
//c85ppgetmvejv84l3dp0
