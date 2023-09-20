package test

import (
	"crm-util-go/common"
	"fmt"
)

func TestReflect(data interface{}) {
	// Get value syntax: interfaceVariable.(type) such as data.(float64)

	switch data.(type) {
	case bool:
		fmt.Println("This is a boolean value: ", data.(bool))
	case int:
		fmt.Println("This is my nice int value: ", data.(int))
	case float64:
		fmt.Println(data.(float64))
	case complex128:
		fmt.Println(data.(complex128))
	case string:
		fmt.Println(data.(string))
	case chan int:
		fmt.Println(data.(chan int))
	default:
		fmt.Println("Unknown type")
	}

	t, v := common.GetReflect(data)
	fmt.Println(fmt.Sprintf("type: %v, value: %v", t, v))

	// Checking type assertions
	var i interface{} = 42
	val, ok := i.(string)

	if ok == false {
		fmt.Println("Wrong type assertion!")
	} else {
		fmt.Println(val)
	}
}
