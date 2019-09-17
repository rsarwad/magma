/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/service/config"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/metricsd/confignames"
	promH "magma/orc8r/cloud/go/services/metricsd/prometheus/handlers"

	"github.com/labstack/echo"
	promAPI "github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
)

const (
	MetricsV1Root = obsidian.V1Root + obsidian.MagmaNetworksUrlPart + "/:network_id" + obsidian.UrlSep + "metrics"
)

// GetObsidianHandlers returns all obsidian handlers for metricsd
func GetObsidianHandlers(configMap *config.ConfigMap) []obsidian.Handler {
	var ret []obsidian.Handler
	client, err := promAPI.NewClient(promAPI.Config{Address: configMap.GetRequiredStringParam(confignames.PrometheusQueryAddress)})
	if err != nil {
		ret = append(ret,
			obsidian.Handler{Path: promH.QueryURL, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},
			obsidian.Handler{Path: promH.QueryRangeURL, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},

			// V1
			obsidian.Handler{Path: promH.QueryV1URL, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},
			obsidian.Handler{Path: promH.QueryRangeV1URL, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},
			obsidian.Handler{Path: promH.SeriesV1URL, Methods: obsidian.GET, HandlerFunc: getInitErrorHandler(err)},
		)
	} else {
		pAPI := v1.NewAPI(client)
		ret = append(ret,
			obsidian.Handler{Path: promH.QueryURL, Methods: obsidian.GET, HandlerFunc: promH.GetPrometheusQueryHandler(pAPI)},
			obsidian.Handler{Path: promH.QueryRangeURL, Methods: obsidian.GET, HandlerFunc: promH.GetPrometheusQueryRangeHandler(pAPI)},

			// V1
			obsidian.Handler{Path: promH.QueryV1URL, Methods: obsidian.GET, HandlerFunc: promH.GetPrometheusQueryHandler(pAPI)},
			obsidian.Handler{Path: promH.QueryRangeV1URL, Methods: obsidian.GET, HandlerFunc: promH.GetPrometheusQueryRangeHandler(pAPI)},
			obsidian.Handler{Path: promH.SeriesV1URL, Methods: obsidian.GET, HandlerFunc: promH.GetPrometheusSeriesHandler(pAPI)},
		)
	}

	alertmanagerConfigServiceURL := configMap.GetRequiredStringParam(confignames.AlertmanagerConfigServiceURL)
	prometheusConfigServiceURL := configMap.GetRequiredStringParam(confignames.PrometheusConfigServiceURL)
	alertmanagerURL := configMap.GetRequiredStringParam(confignames.AlertmanagerApiURL)
	ret = append(ret,
		obsidian.Handler{Path: promH.AlertConfigURL, Methods: obsidian.POST, HandlerFunc: promH.GetConfigurePrometheusAlertHandler(prometheusConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertConfigURL, Methods: obsidian.GET, HandlerFunc: promH.GetRetrieveAlertRuleHandler(prometheusConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertConfigURL, Methods: obsidian.DELETE, HandlerFunc: promH.GetDeleteAlertRuleHandler(prometheusConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertUpdateURL, Methods: obsidian.PUT, HandlerFunc: promH.GetUpdateAlertRuleHandler(prometheusConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertBulkUpdateURL, Methods: obsidian.PUT, HandlerFunc: promH.GetBulkUpdateAlertHandler(prometheusConfigServiceURL)},

		obsidian.Handler{Path: promH.FiringAlertURL, Methods: obsidian.GET, HandlerFunc: promH.GetViewFiringAlertHandler(alertmanagerURL)},
		obsidian.Handler{Path: promH.AlertReceiverConfigURL, Methods: obsidian.POST, HandlerFunc: promH.GetConfigureAlertReceiverHandler(alertmanagerConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertReceiverConfigURL, Methods: obsidian.GET, HandlerFunc: promH.GetRetrieveAlertReceiverHandler(alertmanagerConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertReceiverConfigURL, Methods: obsidian.DELETE, HandlerFunc: promH.GetDeleteAlertReceiverHandler(alertmanagerConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertReceiverUpdateURL, Methods: obsidian.PUT, HandlerFunc: promH.GetUpdateAlertReceiverHandler(alertmanagerConfigServiceURL)},

		obsidian.Handler{Path: promH.AlertReceiverConfigURL + "/route", Methods: obsidian.GET, HandlerFunc: promH.GetRetrieveAlertRouteHandler(alertmanagerConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertReceiverConfigURL + "/route", Methods: obsidian.POST, HandlerFunc: promH.GetUpdateAlertRouteHandler(alertmanagerConfigServiceURL)},
	)

	// V1
	ret = append(ret,
		obsidian.Handler{Path: promH.AlertConfigV1URL, Methods: obsidian.POST, HandlerFunc: promH.GetConfigurePrometheusAlertHandler(prometheusConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertConfigV1URL, Methods: obsidian.GET, HandlerFunc: promH.GetRetrieveAlertRuleHandler(prometheusConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertConfigV1URL, Methods: obsidian.DELETE, HandlerFunc: promH.GetDeleteAlertRuleHandler(prometheusConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertUpdateV1URL, Methods: obsidian.PUT, HandlerFunc: promH.GetUpdateAlertRuleHandler(prometheusConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertBulkUpdateV1URL, Methods: obsidian.PUT, HandlerFunc: promH.GetBulkUpdateAlertHandler(prometheusConfigServiceURL)},

		obsidian.Handler{Path: promH.FiringAlertV1URL, Methods: obsidian.GET, HandlerFunc: promH.GetViewFiringAlertHandler(alertmanagerURL)},
		obsidian.Handler{Path: promH.AlertReceiverConfigV1URL, Methods: obsidian.POST, HandlerFunc: promH.GetConfigureAlertReceiverHandler(alertmanagerConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertReceiverConfigV1URL, Methods: obsidian.GET, HandlerFunc: promH.GetRetrieveAlertReceiverHandler(alertmanagerConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertReceiverConfigV1URL, Methods: obsidian.DELETE, HandlerFunc: promH.GetDeleteAlertReceiverHandler(alertmanagerConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertReceiverUpdateV1URL, Methods: obsidian.PUT, HandlerFunc: promH.GetUpdateAlertReceiverHandler(alertmanagerConfigServiceURL)},

		obsidian.Handler{Path: promH.AlertReceiverConfigV1URL + "/route", Methods: obsidian.GET, HandlerFunc: promH.GetRetrieveAlertRouteHandler(alertmanagerConfigServiceURL)},
		obsidian.Handler{Path: promH.AlertReceiverConfigV1URL + "/route", Methods: obsidian.POST, HandlerFunc: promH.GetUpdateAlertRouteHandler(alertmanagerConfigServiceURL)},

		obsidian.Handler{Path: MetricsV1Root + "/push", Methods: obsidian.POST, HandlerFunc: pushHandler},
	)

	return ret
}

func getInitErrorHandler(err error) func(c echo.Context) error {
	return func(c echo.Context) error {
		return obsidian.HttpError(fmt.Errorf("initialization Error: %v", err), 500)
	}
}

func pushHandler(c echo.Context) error {
	nID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	pushedMetrics := make([]*protos.PushedMetric, 0, 0)
	err := json.NewDecoder(c.Request().Body).Decode(&pushedMetrics)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	metrics := protos.PushedMetricsContainer{
		NetworkId: nID,
		Metrics:   pushedMetrics,
	}
	err = metricsd.PushMetrics(metrics)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
