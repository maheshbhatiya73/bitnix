package server

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type NetStats struct {
	RxBps     float64
	TxBps     float64
	Interface string
}

type ifaceStats struct {
	rx   int64
	tx   int64
	time time.Time
}

var last = map[string]ifaceStats{}

func readBytes(path string) int64 {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	val, _ := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
	return val
}

func ListInterfaces() []string {
	entries, _ := os.ReadDir("/sys/class/net")
	var interfaces []string
	for _, e := range entries {
		name := e.Name()
		statePath := fmt.Sprintf("/sys/class/net/%s/operstate", name)
		data, err := os.ReadFile(statePath)
		if err == nil && strings.TrimSpace(string(data)) == "up" {
			interfaces = append(interfaces, name)
		}
	}
	return interfaces
}

func GetNetStats(iface string) NetStats {
	if iface == "" {
		return NetStats{}
	}

	rxPath := fmt.Sprintf("/sys/class/net/%s/statistics/rx_bytes", iface)
	txPath := fmt.Sprintf("/sys/class/net/%s/statistics/tx_bytes", iface)

	rxNow := readBytes(rxPath)
	txNow := readBytes(txPath)
	now := time.Now()

	prev, exists := last[iface]
	if !exists {
		last[iface] = ifaceStats{rx: rxNow, tx: txNow, time: now}
		return NetStats{RxBps: 0, TxBps: 0, Interface: iface}
	}

	dt := now.Sub(prev.time).Seconds()
	if dt <= 0 {
		return NetStats{RxBps: 0, TxBps: 0, Interface: iface}
	}

	rxRate := float64(rxNow-prev.rx) * 8 / dt
	txRate := float64(txNow-prev.tx) * 8 / dt

	last[iface] = ifaceStats{rx: rxNow, tx: txNow, time: now}

	return NetStats{
		RxBps:     rxRate,
		TxBps:     txRate,
		Interface: iface,
	}
}
