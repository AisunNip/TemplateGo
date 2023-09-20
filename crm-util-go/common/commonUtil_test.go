package common_test

import (
	"crm-util-go/common"
	"fmt"
	"math"
	"testing"
)

func TestRound(t *testing.T) {
	data := 1.128888
	fmt.Println(common.Round(data))
}

func TestCeil(t *testing.T) {
	data := 1.128888
	fmt.Println(common.Ceil(data))
	fmt.Println(math.Ceil(data * 100))
	fmt.Println(math.Ceil(data))
}
