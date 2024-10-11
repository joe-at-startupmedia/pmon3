package file

import (
	"os"
	"strconv"
)

func ModeIntFromString(fileModeStr string) (int64, error) {
	// Parse the octal string to an integer
	return strconv.ParseInt(fileModeStr, 8, 32)
}

func ModeFromString(fileModeStr string) (os.FileMode, error) {
	// Parse the octal string to an integer
	modeInt, err := ModeIntFromString(fileModeStr)
	if err != nil {
		return 0, err
	}
	// Convert the integer to FileMode
	return os.FileMode(modeInt), nil
}
