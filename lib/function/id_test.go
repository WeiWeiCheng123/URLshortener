package function

import (
	"fmt"
	"sort"
	"testing"
	"time"
)

func Test_id(t *testing.T) {
	a := []string{}
	num_test := 1000
	for i := 0; i < num_test; i++ {
		var res string
		res = Generator()
		time.Sleep(1)
		a = append(a, res)

	}
	sort.Strings(a)
	fmt.Println(a)
	err := 0

	for i := 0; i < num_test-1; i++ {
		if a[i] == a[i+1] {
			err++
		}
	}

	fmt.Println("err = ", err)

	t.Fail()
}
