package mapper

import (
	"github.com/quanNaNbk03/monitor-demo/internal/server/dto"
	"github.com/quanNaNbk03/monitor-demo/internal/server/model"
	"github.com/quanNaNbk03/monitor-demo/pkg/api/vm"
)

func ToMonitorDTO(m *model.Monitor) *dto.Monitor {
	if m == nil {
		return nil
	}
	targets := make([]dto.MonitorTarget, 0, len(m.Targets))
	for _, t := range m.Targets {
		targets = append(targets, dto.MonitorTarget{
			Data: dto.MonitorTargetData{
				Timestamps: t.Data.Timestamps,
				Values:     t.Data.Values,
				Status:     t.Data.Status,
			},
			LegendFormat: t.LegendFormat,
		})
	}
	return &vm.Monitor{
		Name:    m.Name,
		Unit:    m.Unit,
		Type:    m.Type,
		Step:    m.Step,
		Targets: targets,
	}
}
