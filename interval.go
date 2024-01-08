package main

import (
	"time"
)

type metricsWithPeriod struct {
	metrics
	interval time.Duration
}

func Period(m metrics, interval time.Duration) metrics {
	return metricsWithPeriod{
		metrics:  m,
		interval: interval,
	}
}

func (m metricsWithPeriod) CheckInterval() time.Duration {
	return m.interval
}
