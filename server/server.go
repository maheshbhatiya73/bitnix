package server

import (
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/net"
)

type NetStats struct {
	RxMbps float64
	TxMbps float64
}

var (
	lastRx uint64
	lastTx uint64
	lastTS time.Time
	mu     sync.Mutex
)

func GetNetStats() NetStats {
	mu.Lock()
	defer mu.Unlock()

	stats := NetStats{}

	counters, _ := net.IOCounters(false)
	if len(counters) == 0 {
		return stats
	}

	now := time.Now()
	elapsed := now.Sub(lastTS).Seconds()
	if elapsed == 0 {
		return stats
	}

	currRx := counters[0].BytesRecv
	currTx := counters[0].BytesSent

	if !lastTS.IsZero() {
		stats.RxMbps = float64(currRx-lastRx) * 8 / (1e6 * elapsed)
		stats.TxMbps = float64(currTx-lastTx) * 8 / (1e6 * elapsed)
	}

	lastRx = currRx
	lastTx = currTx
	lastTS = now

	return stats
}
