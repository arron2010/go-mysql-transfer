package collections

import "go-mysql-transfer/util/stringutil"

func Contain(array []string, v interface{}) bool {
	vvv := stringutil.ToString(v)
	for _, vv := range array {
		if vv == vvv {
			return true
		}
	}
	return false
}

func Uint16ToUInt32(a []uint16) []uint32 {
	b := make([]uint32, len(a), len(a))
	for i := 0; i < len(a); i++ {
		b[i] = uint32(a[i])
	}
	return b
}
