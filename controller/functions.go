package controller

import "strconv"

func StringToUint64(s string) (i uint64) {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	} else {
		return i
	}
}

func StringToUint16(s string) (i uint16) {
	return uint16(StringToUint64(s))
}
