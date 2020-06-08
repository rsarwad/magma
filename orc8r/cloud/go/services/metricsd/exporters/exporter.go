/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package exporters provides an interface for converting protobuf metrics to
// timeseries datapoints and writing these datapoints to storage.
package exporters

import (
	"magma/orc8r/lib/go/protos"

	dto "github.com/prometheus/client_model/go"
)

// Exporter exports metrics received by the metricsd servicer to a datasink.
type Exporter interface {
	// Submit metrics to the exporter.
	// This method must be thread-safe.
	Submit(metrics []MetricAndContext) error

	// Start the metrics export loop.
	// This method is async, i.e. the callee begins a goroutine and
	// returns immediately.
	Start()
}

// MetricAndContext wraps a metric family and metric context
type MetricAndContext struct {
	Family  *dto.MetricFamily
	Context MetricsContext
}

// MetricsContext provides information to the exporter about where this metric
// comes from.
type MetricsContext struct {
	MetricName        string
	AdditionalContext AdditionalMetricContext
}

type AdditionalMetricContext interface {
	isExtraMetricContext()
}

type CloudMetricContext struct {
	// CloudHost is the hostname of the cloud host which originated the metric.
	CloudHost string
}

func (c *CloudMetricContext) isExtraMetricContext() {}

type GatewayMetricContext struct {
	NetworkID, GatewayID string
}

func (c *GatewayMetricContext) isExtraMetricContext() {}

type PushedMetricContext struct {
	NetworkID string
}

func (c *PushedMetricContext) isExtraMetricContext() {}

// ConvertMetricAndContextToProto converts metricAndContext objects to their
// protobuf representation.
func ConvertMetricAndContextToProto(metric MetricAndContext) *protos.MetricAndContext {
	metricAndContext := &protos.MetricAndContext{
		Family: metric.Family,
		Context: &protos.MetricContext{
			MetricName: metric.Context.MetricName,
		},
	}
	switch additionalCtx := metric.Context.AdditionalContext.(type) {
	case *CloudMetricContext:
		metricAndContext.Context.MetricOriginContext = &protos.MetricContext_CloudMetric{
			CloudMetric: &protos.CloudMetricContext{
				CloudHost: additionalCtx.CloudHost,
			},
		}
	case *GatewayMetricContext:
		metricAndContext.Context.MetricOriginContext = &protos.MetricContext_GatewayMetric{
			GatewayMetric: &protos.GatewayMetricContext{
				NetworkId: additionalCtx.NetworkID,
				GatewayId: additionalCtx.GatewayID,
			},
		}
	case *PushedMetricContext:
		metricAndContext.Context.MetricOriginContext = &protos.MetricContext_PushedMetric{
			PushedMetric: &protos.PushedMetricContext{
				NetworkId: additionalCtx.NetworkID,
			},
		}
	}
	return metricAndContext
}
