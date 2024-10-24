package cpu

import (
	"fmt"
	"github.com/struCoder/pidusage"
	"pmon3/utils/conv"
)

func GetUsageStats(pid int) (string, string) {
	cpuVal := "0"
	memVal := "0"

	info, err := pidusage.GetStat(pid)
	if err != nil {
		return cpuVal, memVal
	}

	cpuVal = conv.FloatToStr(info.CPU)

	if info.Memory <= 1024 {
		memVal = conv.FloatToStr(info.Memory)
	} else if info.Memory <= 1024*1024 {
		memVal = fmt.Sprintf("%.1f KB", info.Memory/float64(1024))
	} else if info.Memory <= 1024*1024*1024 {
		memVal = fmt.Sprintf("%.1f MB", info.Memory/float64(1024*1024))
	} else if info.Memory <= 1024*1024*1024*1024 {
		memVal = fmt.Sprintf("%.1f GB", info.Memory/float64(1024*1024*1024))
	} else {
		memVal = conv.FloatToStr(info.Memory)
	}

	return cpuVal + "%", memVal
}
