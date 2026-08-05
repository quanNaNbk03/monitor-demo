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

func (Monitor) TableName() string {
	return "monitors"
}

type MonitorTarget struct {
	Data         MonitorTargetData
	LegendFormat string
}

func (MonitorTarget) TableName() string {
	return "monitor_targets"
}

type MonitorTargetData struct {
	Timestamps []int64
	Values     []float64
	Status     string // for type status(status history)
}

func (MonitorTargetData) TableName() string {
	return "monitor_target_datas"
}

type MetricByTypeConfig struct {
	Name    string
	Type    string
	Unit    string
	Label   string
	Targets []MetricByTypeConfigTarget
}

func (MetricByTypeConfig) TableName() string {
	return "metric_by_type_configs"
}

type MetricByTypeConfigTarget struct {
	Expr         string
	LegendFormat string
}

func (MetricByTypeConfigTarget) TableName() string {
	return "metric_by_type_config_targets"
}

type MetricConfig struct {
	Name string
	Expr string
	Unit string
	Type string
}

func (MetricConfig) TableName() string {
	return "metric_configs"
}

type MonitorByTypeParam struct {
	Start time.Time
	End   time.Time
	Step  int64
}

func (MonitorByTypeParam) TableName() string {
	return "monitor_by_type_params"
}
