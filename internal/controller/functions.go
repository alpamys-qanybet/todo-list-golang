package controller

import "strconv"

func StringToUint16(s string) (uint16, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint16(i), nil
}
