package model

import (
	"time"
)

type Monitor struct {
	Name    string
	Unit    string
	Type    string
	Step    int64
	Targets []*MonitorTarget
}

type MonitorTarget struct {
	Data         MonitorTargetData
	LegendFormat string
}

type MonitorTargetData struct {
	Timestamps []int64
	Values     []float64
	Status     string // for type status(status history)
}

type MetricByTypeConfig struct {
	Name    string
	Type    string
	Unit    string
	Label   string
	Targets []MetricByTypeConfigTarget
}

type MetricByTypeConfigTarget struct {
	Expr         string
	LegendFormat string
}

type MetricConfig struct {
	Name string
	Expr string
	Unit string
	Type string
}

type MonitorByTypeParam struct {
	Start time.Time
	End   time.Time
	Step  int64
}
