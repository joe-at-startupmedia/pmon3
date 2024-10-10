package process

import (
	"pmon3/conf"
	"pmon3/pmond"
	"sync"
	"time"
)

type FlapDetector struct {
	lastForgiven       time.Time
	restarted          int
	countdown          int
	thresholdRestarted int
	thresholdCountdown int
	thresholdDecrement int
}

var mutex = sync.Mutex{}
var flapDetectors = map[uint32]*FlapDetector{}

func GetFlapDetectorByProcessId(processId uint32, conf *conf.Config) *FlapDetector {
	mutex.Lock()
	fd := flapDetectors[processId]
	if fd == nil {
		fd = &FlapDetector{
			time.Now(),
			0,
			0,
			int(conf.FlapDetection.ThresholdRestarted),
			int(conf.FlapDetection.ThresholdCountdown),
			int(conf.FlapDetection.ThresholdDecrement),
		}
		flapDetectors[processId] = fd
	}
	mutex.Unlock()
	return fd
}

func (fd *FlapDetector) ShouldBackOff(evaluationInterval time.Duration) bool {
	if fd.countdown == 1 {
		fd.countdown = 0
		fd.restarted = 0
	} else if fd.countdown > 1 {
		fd.countdown--
		pmond.Log.Debugf("BACKOFF decremented: %d", fd.countdown)

		if fd.thresholdDecrement > 0 && fd.restarted > 0 && time.Since(fd.lastForgiven) > time.Duration(fd.thresholdDecrement)*evaluationInterval {
			pmond.Log.Infof("BACKOFF decrementing restarted: %d, time since %d %d", fd.countdown, time.Since(fd.lastForgiven), time.Duration(fd.thresholdDecrement)*time.Millisecond*time.Duration(pmond.Config.ProcessMonitorInterval))
			fd.restarted--
			fd.lastForgiven = time.Now()
		}
		if fd.restarted >= fd.thresholdRestarted {
			return true
		}
	}

	return false
}

func (fd *FlapDetector) RestartProcess() {
	fd.restarted++
	fd.lastForgiven = time.Now()
	pmond.Log.Infof("fd: %-v", fd)
	if fd.restarted >= fd.thresholdRestarted {
		//begin backoff
		pmond.Log.Infof("BACKOFF triggered")
		fd.countdown = fd.thresholdCountdown
		fd.lastForgiven = time.Now()
	}
}
