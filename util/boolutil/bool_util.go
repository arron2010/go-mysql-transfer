package boolutil

func BoolToInt(b bool) uint32 {
	if b {
		return 1
	} else {
		return 0
	}
}
