package pkg

import "reflect"

// in_array
// 使用反射进行处理，速度较慢，如果业务中要求的逻辑较为简单，尽量自行实现
// 仅适用haystack为slice，array或map类型，其余类型默认返回false
func InArray(needle, haystack interface{}) bool {
	vHay := reflect.ValueOf(haystack)
	vKind := vHay.Kind()
	if vKind == reflect.Slice || vKind == reflect.Array {
		for i := 0; i < vHay.Len(); i++ {
			if reflect.DeepEqual(vHay.Index(i).Interface(), needle) {
				return true
			}
		}
	} else if vKind == reflect.Map {
		iter := vHay.MapRange()
		for iter.Next() {
			if reflect.DeepEqual(iter.Value().Interface(), needle) {
				return true
			}
		}
	}
	return false
}