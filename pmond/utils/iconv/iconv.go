package iconv

import (
	"pmon3/pmond"
	"strconv"
)

func MustInt(val string) int {
	valInt, err := strconv.Atoi(val)
	if err != nil {
		pmond.Log.Debug(err)
	}
	return valInt
}
