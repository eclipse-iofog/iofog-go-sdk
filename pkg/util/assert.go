package util

import "fmt"

func AssertInt(in interface{}) int {
	out, ok := in.(int)
	if !ok {
		floatOut, ok := in.(float64)
		if ok {
			return int(floatOut)
		}
		panic(fmt.Sprintf(pkg.errorVariableNotInteger, in))
	}
	return out
}

func AssertBool(in interface{}) bool {
	out, ok := in.(bool)
	if !ok {
		panic(fmt.Sprintf(pkg.errorVariableNotBool, in))
	}
	return out
}
