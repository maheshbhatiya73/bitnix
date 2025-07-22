package tui
import "sort"
import "bitnix/server"

type SortBy int

const (
	SortByTotal SortBy = iota
	SortByRx
	SortByTx
)

func GetAllProcessNetStatsSorted(by SortBy) []server.ProcessNetStats {
	procs := server.GetAllProcessNetStats()
	sort.Slice(procs, func(i, j int) bool {
		switch by {
		case SortByRx:
			return procs[i].RxBps > procs[j].RxBps
		case SortByTx:
			return procs[i].TxBps > procs[j].TxBps
		default:
			return (procs[i].RxBps + procs[i].TxBps) > (procs[j].RxBps + procs[j].TxBps)
		}
	})
	if len(procs) > 30 {
		return procs[:30]
	}
	return procs
}
