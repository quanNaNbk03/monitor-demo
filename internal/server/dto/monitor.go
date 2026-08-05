package dto

// Monitor defines model for Monitor.
type Monitor struct {
	Name    string          `json:"name,omitempty"`
	Step    int64           `json:"step,omitempty"`
	Targets []MonitorTarget `json:"targets,omitempty"`
	Type    string          `json:"type,omitempty"`
	Unit    string          `json:"unit,omitempty"`
}

// MonitorTarget defines model for MonitorTarget.
type MonitorTarget struct {
	Data         MonitorTargetData `json:"data,omitempty"`
	LegendFormat string            `json:"legendFormat,omitempty"`
}

// MonitorTargetData defines model for MonitorTargetData.
type MonitorTargetData struct {
	Status     string    `json:"status,omitempty"`
	Timestamps []int64   `json:"timestamps,omitempty"`
	Values     []float64 `json:"values,omitempty"`
}
