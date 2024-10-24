package e2e

import (
	"pmon3/model"
)

func init() {
	model.ProcessUsageStatsAccessor = new(processUsageStatsMock)
}

type processUsageStatsMock struct{}

func (p *processUsageStatsMock) GetUsageStats(_ int) (string, string) {
	return "1.1 MB", "1.1%"
}
