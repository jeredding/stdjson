package rewriter

import (
	"reflect"
	"fmt"
)

func ClonedMap(m map[string]interface{}) (map[string]interface{}, error) {
	res := cloneRecursive(reflect.ValueOf(m))

	if !res.CanInterface() {
		return nil, fmt.Errorf("unable to clone: %v, got %v", m, res)
	}

	if resMap, ok := res.Interface().(map[string]interface{}); ok {
		return resMap, nil
	}

	return nil, fmt.Errorf("unable to clone %v. got: %v", m, res)
}

func cloneRecursive(in reflect.Value) reflect.Value {
	switch in.Kind() {
	case reflect.Map:
		res := reflect.MakeMap(in.Type())
		for _, key := range in.MapKeys() {
			ii := in.MapIndex(key)
			res.SetMapIndex(key, cloneRecursive(ii))
		}
		return res
	case reflect.Slice:
		res := reflect.MakeSlice(in.Type(), in.Len(), in.Cap())
		for i := 0; i < in.Len(); i++ {
			ri := res.Index(i)
			ri.Set(cloneRecursive(in.Index(i)))
		}
		return res
	case reflect.Interface:
		iv := in.Elem()
		if !iv.IsValid() {
			return in
		}
		return cloneRecursive(iv)
	default:
		return in
	}
}
