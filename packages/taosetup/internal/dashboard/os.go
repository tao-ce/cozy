package dashboard

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func getUptime() time.Duration {
	uptime, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0
	}
	uptimeFloat, err := strconv.ParseFloat(strings.Fields(strings.TrimSpace(string(uptime)))[0], 64)
	if err != nil {
		fmt.Println("error parsing uptime:", err)
		return 0
	}
	return time.Duration(uptimeFloat * float64(time.Second))
}
