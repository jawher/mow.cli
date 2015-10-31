package cli

import (
	"fmt"
	"reflect"
)

func formatterFor(t reflect.Type) func(interface{}) string {
	switch t.Kind() {
	case reflect.Bool:
		return boolFormatter
	case reflect.String:
		return stringFormatter
	case reflect.Int:
		return intFormatter
	case reflect.Slice:
		switch t.Elem().Kind() {
		case reflect.String:
			return stringsFormatter
		case reflect.Int:
			return intsFormatter
		default:
			panic(fmt.Sprintf("No formatter for %v", t))
		}
	default:
		panic(fmt.Sprintf("No formatter for %v", t))
	}
}

func boolFormatter(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func stringFormatter(v interface{}) string {
	return fmt.Sprintf("%#v", v)
}

func intFormatter(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func stringsFormatter(v interface{}) string {
	res := "["
	strings, _ := v.([]string)
	for idx, s := range strings {
		if idx > 0 {
			res += ", "
		}
		res += fmt.Sprintf("%#v", s)
	}
	return res + "]"
}

func intsFormatter(v interface{}) string {
	res := "["
	ints, _ := v.([]int)
	for idx, s := range ints {
		if idx > 0 {
			res += ", "
		}
		res += fmt.Sprintf("%v", s)
	}
	return res + "]"
}
