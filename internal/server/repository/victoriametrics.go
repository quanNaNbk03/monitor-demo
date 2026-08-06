package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	prometheusmodel "github.com/prometheus/common/model"
	"github.com/quanNaNbk03/monitor-demo/internal/server/model"
	"github.com/quanNaNbk03/monitor-demo/pkg/common/errlist"
	"github.com/spf13/viper"

	"git.ocn.com.vn/ocn/common/httpbase/ierror"
	"git.ocn.com.vn/ocn/common/logger"
	"git.ocn.com.vn/ocn/common/prometheus/client"
	"git.ocn.com.vn/ocn/common/prometheus/client/interval"
)

type VictoriaMetricsRepository interface {
	QueryRange(ctx context.Context, metricByTypeConfig model.MetricByTypeConfig, req *model.MonitorByTypeParam) (*model.Monitor, *ierror.CoreError)
}

func NewVictoriaMetricsClient() (VictoriaMetricsRepository, error) {
	c, err := client.NewClient(viper.GetString("victoriaMetric.address"), viper.GetDuration("victoriaMetric.scrapeInterval"))
	if err != nil {
		return nil, fmt.Errorf("error creating victoriametrics client: %v", err)
	}
	return &VictoriaMetricsClient{
		c: c,
	}, nil
}

type VictoriaMetricsClient struct {
	c *client.Client
}

func (c *VictoriaMetricsClient) QueryRange(ctx context.Context, metricByTypeConfig model.MetricByTypeConfig, req *model.MonitorByTypeParam) (*model.Monitor, *ierror.CoreError) {
	var stepDuration time.Duration
	if req.Step > 0 {
		stepDuration = interval.Round(time.Second * time.Duration(req.Step))
	} else {
		stepDuration = interval.CalculateRateInterval(req.Start, req.End, c.c.ScrapeInterval)
	}
	monitorResponse := &model.Monitor{
		Name: metricByTypeConfig.Name,
		Unit: metricByTypeConfig.Unit,
		Type: metricByTypeConfig.Type,
		Step: int64(stepDuration.Seconds()),
	}
	var (
		monitorResponseTargets []*model.MonitorTarget
		labelValues            prometheusmodel.LabelValues
	)
	// Get all values of label
	if metricByTypeConfig.Label != "" {
		var err error
		labelValues, err = c.c.LabelValues(ctx, metricByTypeConfig.Label, nil, req.Start, req.End)
		if err != nil && !errors.Is(err, client.ErrClientIsNil) {
			return nil, errlist.ErrDatabase.WithChild(fmt.Errorf("get label values failed: %w", err))
		}
		if len(labelValues) == 0 {
			return monitorResponse, nil
		}
	}

	for _, target := range metricByTypeConfig.Targets {
		var (
			targetTemps []model.MetricByTypeConfigTarget
		)
		if len(labelValues) > 0 {
			for _, label := range labelValues {
				targetTemps = append(targetTemps, model.MetricByTypeConfigTarget{
					Expr:         strings.ReplaceAll(target.Expr, fmt.Sprintf("$%s", metricByTypeConfig.Label), string(label)),
					LegendFormat: strings.ReplaceAll(target.LegendFormat, fmt.Sprintf("$%s", metricByTypeConfig.Label), string(label)),
				})
			}
		} else {
			targetTemps = append(targetTemps, model.MetricByTypeConfigTarget{
				Expr:         target.Expr,
				LegendFormat: target.LegendFormat,
			})
		}

		for _, targetTemp := range targetTemps {
			modelValue, err := c.c.QueryRange(ctx, targetTemp.Expr, req.Start, req.End, stepDuration)
			if err != nil && !errors.Is(err, client.ErrClientIsNil) {
				logger.Error("query to get metric fail", "err", err, "query", targetTemp.Expr, "start", req.Start.String(), "end", req.End.String())
				continue
			}
			monitorResponseTarget := convertMatrixValueVM(targetTemp, modelValue)
			if monitorResponseTarget != nil {
				monitorResponseTargets = append(monitorResponseTargets, monitorResponseTarget)
			}
		}
	}
	monitorResponse.Targets = monitorResponseTargets
	return monitorResponse, nil
}

func convertMatrixValueVM(m model.MetricByTypeConfigTarget, value prometheusmodel.Value) *model.MonitorTarget {
	valueMatrices, ok := value.(prometheusmodel.Matrix)
	if !ok {
		return nil
	}

	if len(valueMatrices) == 0 || len(valueMatrices[0].Values) == 0 {
		return nil
	}

	var (
		resultTimes  []int64
		resultValues []float64
	)
	for _, valueMatrix := range valueMatrices {
		for _, v := range valueMatrix.Values {
			resultTimes = append(resultTimes, v.Timestamp.Unix())
			resultValues = append(resultValues, float64(v.Value))
		}
	}

	return &model.MonitorTarget{
		Data: model.MonitorTargetData{
			Timestamps: resultTimes,
			Values:     resultValues,
		},
		LegendFormat: m.LegendFormat,
	}
}
