package server

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ProcessNetStats struct {
	PID     int
	Name    string
	RxBps   float64
	TxBps   float64
}

type procNetStats struct {
	rx   int64
	tx   int64
	time time.Time
}

var lastProcStats = map[int]procNetStats{}

func getAllPIDs() []int {
	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		return nil
	}
	var pids []int
	for _, f := range files {
		if f.IsDir() {
			if pid, err := strconv.Atoi(f.Name()); err == nil {
				pids = append(pids, pid)
			}
		}
	}
	return pids
}

func getProcName(pid int) string {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil || len(data) == 0 {
		return "[unknown]"
	}
	parts := strings.Split(string(data), "\x00")
	return filepath.Base(parts[0])
}

func getProcNetBytes(pid int) (int64, int64) {
	netPath := fmt.Sprintf("/proc/%d/net/dev", pid)
	data, err := os.ReadFile(netPath)
	if err != nil {
		return 0, 0
	}

	var rxTotal, txTotal int64
	lines := strings.Split(string(data), "\n")[2:] // skip headers
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 17 {
			continue
		}
		rx, _ := strconv.ParseInt(fields[1], 10, 64)
		tx, _ := strconv.ParseInt(fields[9], 10, 64)
		rxTotal += rx
		txTotal += tx
	}
	return rxTotal, txTotal
}

func GetAllProcessNetStats() []ProcessNetStats {
	now := time.Now()
	pids := getAllPIDs()
	var results []ProcessNetStats

	for _, pid := range pids {
		rxNow, txNow := getProcNetBytes(pid)
		prev, ok := lastProcStats[pid]
		if !ok {
			lastProcStats[pid] = procNetStats{rx: rxNow, tx: txNow, time: now}
			continue
		}

		dt := now.Sub(prev.time).Seconds()
		if dt <= 0 {
			continue
		}

		rxBps := float64(rxNow-prev.rx) * 8 / dt
		txBps := float64(txNow-prev.tx) * 8 / dt

		lastProcStats[pid] = procNetStats{rx: rxNow, tx: txNow, time: now}

		if rxBps > 0 || txBps > 0 {
			results = append(results, ProcessNetStats{
				PID:   pid,
				Name:  getProcName(pid),
				RxBps: rxBps,
				TxBps: txBps,
			})
		}
	}
	return results
}
