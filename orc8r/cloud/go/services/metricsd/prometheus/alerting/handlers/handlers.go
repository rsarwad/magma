/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/alert"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/prometheus/prometheus/pkg/rulefmt"
)

const (
	AlertPath       = rootPath + "/alert"
	AlertUpdatePath = AlertPath + "/:" + RuleNamePathParam
	AlertBulkPath   = AlertPath + "/bulk"

	prometheusReloadPath = "/-/reload"
	ruleNameQueryParam   = "alert_name"
	RuleNamePathParam    = "alert_name"
)

// GetConfigureAlertHandler returns a handler that calls the client method WriteAlert() to
// write the alert configuration from the body of this request
func GetConfigureAlertHandler(client alert.PrometheusAlertClient, prometheusURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		rule, err := decodeRulePostRequest(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		networkID := getNetworkID(c)

		err = client.ValidateRule(rule)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if client.RuleExists(networkID, rule.Alert) {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Rule '%s' already exists", rule.Alert))
		}

		err = client.WriteRule(networkID, rule)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = reloadPrometheus(prometheusURL)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}
}

func GetRetrieveAlertHandler(client alert.PrometheusAlertClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		ruleName := c.QueryParam(ruleNameQueryParam)
		networkID := getNetworkID(c)
		rules, err := client.ReadRules(networkID, ruleName)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		jsonRules, err := rulesToJSON(rules)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, jsonRules)
	}
}

func GetDeleteAlertHandler(client alert.PrometheusAlertClient, prometheusURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		ruleName := c.QueryParam(ruleNameQueryParam)
		networkID := getNetworkID(c)
		if ruleName == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "No rule name provided")
		}
		err := client.DeleteRule(networkID, ruleName)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		err = reloadPrometheus(prometheusURL)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.String(http.StatusOK, fmt.Sprintf("rule %s deleted", ruleName))
	}
}

func GetUpdateAlertHandler(client alert.PrometheusAlertClient, prometheusURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		ruleName := c.Param(RuleNamePathParam)
		networkID := getNetworkID(c)
		if ruleName == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "No rule name provided")
		}

		if !client.RuleExists(networkID, ruleName) {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Rule '%s' does not exist", ruleName))
		}

		rule, err := decodeRulePostRequest(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = client.ValidateRule(rule)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = client.UpdateRule(networkID, rule)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = reloadPrometheus(prometheusURL)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}
}

func GetBulkAlertUpdateHandler(client alert.PrometheusAlertClient, prometheusURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID := getNetworkID(c)

		rules, err := decodeBulkRulesPostRequest(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		for _, rule := range rules {
			err = client.ValidateRule(rule)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
		}

		results, err := client.BulkUpdateRules(networkID, rules)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = reloadPrometheus(prometheusURL)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, results)
	}
}

func decodeRulePostRequest(c echo.Context) (rulefmt.Rule, error) {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return rulefmt.Rule{}, fmt.Errorf("error reading request body: %v", err)
	}
	payload := rulefmt.Rule{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return payload, fmt.Errorf("error unmarshalling payload: %v", err)
	}
	return payload, nil
}

func decodeBulkRulesPostRequest(c echo.Context) ([]rulefmt.Rule, error) {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return []rulefmt.Rule{}, fmt.Errorf("error reading request body: %v", err)
	}
	var payload []rulefmt.Rule
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return payload, fmt.Errorf("error unmarshalling payload: %v", err)
	}
	return payload, nil
}

func reloadPrometheus(url string) error {
	if url == "" {
		glog.Info("Not reloading prometheus. No url given.")
		return nil
	}
	resp, err := http.Post(fmt.Sprintf("http://%s%s", url, prometheusReloadPath), "text/plain", &bytes.Buffer{})
	if err != nil {
		return fmt.Errorf("error reloading prometheus: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("error reloading prometheus (status %d): %s", resp.StatusCode, string(body))
	}
	return nil
}

func getNetworkID(c echo.Context) string {
	return c.Param("network_id")
}

func rulesToJSON(rules []rulefmt.Rule) ([]alert.RuleJSONWrapper, error) {
	ret := make([]alert.RuleJSONWrapper, 0)

	for _, rule := range rules {
		jsonRule, err := rulefmtToJSON(rule)
		if err != nil {
			return ret, err
		}
		ret = append(ret, *jsonRule)
	}
	return ret, nil
}

func rulefmtToJSON(rule rulefmt.Rule) (*alert.RuleJSONWrapper, error) {
	duration, err := time.ParseDuration(rule.For.String())
	if err != nil {
		return nil, err
	}
	return &alert.RuleJSONWrapper{
		Record:      rule.Record,
		Alert:       rule.Alert,
		Expr:        rule.Expr,
		For:         duration.String(),
		Labels:      rule.Labels,
		Annotations: rule.Annotations,
	}, nil

}
