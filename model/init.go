package model

import "pmon3/utils/cpu"

func init() {
	ProcessUsageStatsAccessor = new(processUsageStats)
}

type processUsageStats struct{}

func (p *processUsageStats) GetUsageStats(pid int) (string, string) {
	return cpu.GetUsageStats(pid)
}
