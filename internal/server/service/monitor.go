package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/quanNaNbk03/monitor-demo/internal/server/mapper"
	"github.com/quanNaNbk03/monitor-demo/internal/server/model"
	"github.com/quanNaNbk03/monitor-demo/internal/server/repository"
	"github.com/quanNaNbk03/monitor-demo/pkg/api/vm"
	"github.com/quanNaNbk03/monitor-demo/pkg/common/errlist"

	"git.ocn.com.vn/ocn/common/httpbase"
	"git.ocn.com.vn/ocn/common/httpbase/ierror"
)

type MonitorVMService interface {
	GetMonitorVMByType(ctx context.Context, id uint64, metric vm.VMMetricType, params vm.GetMonitorVMWithPromParams) (*vm.Monitor, *ierror.Error)
}

type monitorVMService struct {
	prometheusRepo      repository.PrometheusRepository
	hostMetricConfigMap map[vm.VMMetricType]model.MetricByTypeConfig
}

func NewMonitorVMService(prometheusRepo repository.PrometheusRepository) MonitorVMService {
	hostMetricConfigMap := map[vm.VMMetricType]model.MetricByTypeConfig{
		"cpu": {
			Name: "CPU Usage",
			Type: "time_series",
			Unit: "percent",
			Targets: []model.MetricByTypeConfigTarget{
				{
					Expr:         "100 - (avg(irate(node_cpu_seconds_total{neuron_ip=\"$neuron_ip\",mode=\"idle\"}[$__rate_interval])) * 100)",
					LegendFormat: "CPU Usage",
				},
			},
		},
		"diskBytes": {
			Name:  "Disk bytes",
			Type:  "time_series",
			Unit:  "bytes",
			Label: "device",
			Targets: []model.MetricByTypeConfigTarget{
				{
					Expr:         "irate(node_disk_reads_bytes_total{neuron_ip=\"$neuron_ip\",device=\"$device\"}[$__rate_interval])",
					LegendFormat: "Disk Read Bytes - $device",
				},
				{
					Expr:         "irate(node_disk_written_bytes_total{neuron_ip=\"$neuron_ip\",device=\"$device\"}[$__rate_interval])",
					LegendFormat: "Disk Write Bytes - $device",
				},
			},
		},
		"diskOps": {
			Name:  "Disk operation/sec",
			Type:  "time_series",
			Unit:  "iops",
			Label: "device",
			Targets: []model.MetricByTypeConfigTarget{
				{
					Expr:         "irate(node_disk_reads_completed_total{neuron_ip=\"$neuron_ip\",device=\"$device\"}[$__rate_interval])",
					LegendFormat: "Disk Read Operations - $device",
				},
				{
					Expr:         "irate(node_disk_writes_completed_total{neuron_ip=\"$neuron_ip\",device=\"$device\"}[$__rate_interval])",
					LegendFormat: "Disk Write Operations - $device",
				},
			},
		},
		"memory": {
			Name: "Memory Usage",
			Type: "time_series",
			Unit: "bytes",
			Targets: []model.MetricByTypeConfigTarget{
				{
					Expr:         "node_memory_MemTotal_bytes{neuron_ip=\"$neuron_ip\"} - node_memory_MemFree_bytes{neuron_ip=\"$neuron_ip\"} - (node_memory_Cached_bytes{neuron_ip=\"$neuron_ip\"} + node_memory_Buffers_bytes{neuron_ip=\"$neuron_ip\"} + node_memory_SReclaimable_bytes{neuron_ip=\"$neuron_ip\"})",
					LegendFormat: "Memory Usage",
				},
				{
					Expr:         "node_memory_MemTotal_bytes{neuron_ip=\"$neuron_ip\"}",
					LegendFormat: "Memory Total",
				},
				{
					Expr:         "node_memory_Cached_bytes{neuron_ip=\"$neuron_ip\"} + node_memory_Buffers_bytes{neuron_ip=\"$neuron_ip\"} + node_memory_SReclaimable_bytes{neuron_ip=\"$neuron_ip\"}",
					LegendFormat: "Memory Cache+Buffer",
				},
				{
					Expr:         "node_memory_MemFree_bytes{neuron_ip=\"$neuron_ip\"}",
					LegendFormat: "Memory Free",
				},
				{
					Expr:         "(node_memory_SwapTotal_bytes{neuron_ip=\"$neuron_ip\"} - node_memory_SwapFree_bytes{neuron_ip=\"$neuron_ip\"})",
					LegendFormat: "Memory Swap Used",
				},
			},
		},
		"networkBits": {
			Name:  "Network Bits",
			Type:  "time_series",
			Unit:  "bps",
			Label: "device",
			Targets: []model.MetricByTypeConfigTarget{
				{
					Expr:         "irate(node_network_receive_bytes_total{neuron_ip=\"$neuron_ip\",device=\"$device\"}[$__rate_interval])*8",
					LegendFormat: "Network-RX - $device",
				},
				{
					Expr:         "irate(node_network_transmit_bytes_total{neuron_ip=\"$neuron_ip\",device=\"$device\"}[$__rate_interval])*8",
					LegendFormat: "Network-TX - $device",
				},
			},
		},
		"networkPackets": {
			Name:  "Network Packets",
			Type:  "time_series",
			Unit:  "pps",
			Label: "device",
			Targets: []model.MetricByTypeConfigTarget{
				{
					Expr:         "irate(node_network_receive_packets_total{neuron_ip=\"$neuron_ip\",device=\"$device\"}[$__rate_interval])*8",
					LegendFormat: "Network-RX - $device",
				},
				{
					Expr:         "irate(node_network_transmit_packets_total{neuron_ip=\"$neuron_ip\",device=\"$device\"}[$__rate_interval])*8",
					LegendFormat: "Network-TX - $device",
				},
			},
		},
	}
	return &monitorVMService{
		prometheusRepo:      prometheusRepo,
		hostMetricConfigMap: hostMetricConfigMap,
	}
}

func (s *monitorVMService) GetMonitorVMByType(ctx context.Context, id uint64, metric vm.VMMetricType, params vm.GetMonitorVMWithPromParams) (*vm.Monitor, *ierror.Error) {
	neuron, coreErr := s.neuronRepo.GetNeuronByID(ctx, id)
	if coreErr != nil {
		if errors.Is(coreErr, errlist.ErrNeuronNotFound) {
			return nil, httpbase.ErrNotFound(ctx, "neuron not found")
		}
		return nil, httpbase.ErrInternal(ctx, "get neuron error").SetSubError(coreErr)
	}
	config, ok := s.hostMetricConfigMap[metric]
	if !ok {
		return nil, httpbase.ErrBadRequest(ctx, "invalid metric type")
	}

	processedConfig := model.MetricByTypeConfig{
		Name:  config.Name,
		Type:  config.Type,
		Unit:  config.Unit,
		Label: config.Label,
	}

	for _, target := range config.Targets {
		processedConfig.Targets = append(processedConfig.Targets, model.MetricByTypeConfigTarget{
			Expr: strings.ReplaceAll(target.Expr, "$neuron_ip", neuron.IP),
		})
	}

	res, coreErr := s.prometheusRepo.QueryRange(ctx, processedConfig, &model.MonitorByTypeParam{
		Start: time.Unix(params.Start, 0),
		End:   time.Unix(params.End, 0),
		Step:  params.Step,
	})
	if coreErr != nil {
		return nil, httpbase.ErrInternal(ctx, "failed to query monitor by type").SetSubError(coreErr)
	}

	return mapper.ToMonitorDTO(res), nil
}
