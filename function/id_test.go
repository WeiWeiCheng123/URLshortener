package function

import (
	"fmt"
	"testing"
)

func Test_id(t *testing.T) {
	for i := 1; i < 10; i++ {
		fmt.Println(Id())
	}
	
}
