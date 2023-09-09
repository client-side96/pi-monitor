package stats

type Stats struct {
	Temperature float64 `json:"temperature"`
	CPULoad     float64 `json:"cpuLoad"`
	MemoryUsage float64 `json:"memory"`
}
