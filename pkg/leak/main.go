package leak

import (
	"fmt"
	"runtime"
	"time"

	"github.com/stackrox/rox/pkg/logging"
)

var (
	logger = logging.LoggerForModule()
)

func k() {
	var i [5000]int
	for x := 0; x < len(i); x++ {
		i[x] = x
	}
	PrintMemUsage()
}

func main1() {
	// Print our starting memory usage (should be around 0mb)
	PrintMemUsage()
	k()

	var overall [][]int
	for i := 0; i < 4; i++ {
		time.Sleep(10 * time.Second)
		// Allocate memory using make() and append to overall (so it doesn't get
		// garbage collected). This is to create an ever increasing memory usage
		// which we can track. We're just using []int as an example.
		a := make([]int, 0, 999999)
		overall = append(overall, a)

		// Print our memory usage at each interval
		PrintMemUsage()
		time.Sleep(time.Second)
	}
	//for x := 0; x < len(i); x++ {
	//	i[x] = 7
	//}

	// Clear our memory and print usage, unless the GC has run 'Alloc' will remain the same
	overall = nil
	PrintMemUsage()

	// Force GC to clear up, should see a memory drop
	runtime.GC()
	PrintMemUsage()
	fmt.Println("test stack")

}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	msg := fmt.Sprintf("\tSys = %v MiB", bToMb(m.Sys))
	msg += fmt.Sprintf("\tHeapIdle = %v MiB", bToMb(m.HeapIdle))
	msg += fmt.Sprintf("\tHeapAlloc = %v MiB", bToMb(m.HeapAlloc))
	msg += fmt.Sprintf("\tHeapFrag = %v MiB", bToMb(m.HeapInuse-m.HeapAlloc))
	msg += fmt.Sprintf("\tStackInuse = %v", bToMb(m.StackInuse))
	msg += fmt.Sprintf("\tNumGC = %v", m.NumGC)
	msg += fmt.Sprintf("\tStackSys+ = %v", bToMb(m.StackSys-m.StackInuse))
	msg += fmt.Sprintf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	logger.Info(msg)
	logger.Infof("Struct %+v\n", m)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
