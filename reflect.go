package cli

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func vconv(s string, to reflect.Type) (reflect.Value, error) {
	switch to.Kind() {
	case reflect.String:
		return reflect.ValueOf(s), nil
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(b), nil
	case reflect.Int:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(int(i)), nil
	case reflect.Slice:
		res := reflect.New(to)
		vs := strings.Split(s, ",")
		for _, v := range vs {
			conv, err := vconv(strings.TrimSpace(v), to.Elem())
			if err != nil {
				return reflect.Value{}, err
			}
			res.Elem().Set(reflect.Append(res.Elem(), conv))
		}
		return res.Elem(), nil
	default:
		panic(fmt.Sprintf("Unhandled conversion to %v", to))
	}
}

func vset(into reflect.Value, s string) error {
	dest := into.Elem()

	switch dest.Type().Kind() {
	case reflect.Slice:
		v, err := vconv(s, dest.Type().Elem())
		if err != nil {
			return err
		}
		dest.Set(reflect.Append(dest, v))
	default:
		conv, err := vconv(s, dest.Type())
		if err != nil {
			return err
		}
		dest.Set(conv)
	}
	return nil
}

func vinit(into reflect.Value, envVars string, defaultValue interface{}) {
	if len(envVars) > 0 {
		for _, rev := range strings.Split(envVars, " ") {
			ev := strings.TrimSpace(rev)
			if len(ev) > 0 {
				v := os.Getenv(ev)
				if len(v) > 0 {
					conv, err := vconv(v, into.Elem().Type())
					if err == nil {
						into.Elem().Set(conv)
						return
					}
				}
			}
		}

	}
	into.Elem().Set(reflect.ValueOf(defaultValue))
}
