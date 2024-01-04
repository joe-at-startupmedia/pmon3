package conv

import (
	"fmt"
	"strconv"
)

func StrToUint32(val string) uint32 {
	valInt, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println(err)
	}
	return uint32(valInt)
}

func Uint32ToStr(num uint32) string {
	valStr := strconv.Itoa(int(num))
	return valStr
}
