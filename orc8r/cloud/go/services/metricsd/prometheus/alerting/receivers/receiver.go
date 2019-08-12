/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package receivers

import (
	"fmt"
	"strings"
	"time"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/exporters"

	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/common/model"
)

// Config uses a custom receiver struct to avoid scrubbing of 'secrets' during
// marshaling
type Config struct {
	Global       *config.GlobalConfig  `yaml:"global,omitempty" json:"global,omitempty"`
	Route        *config.Route         `yaml:"route,omitempty" json:"route,omitempty"`
	InhibitRules []*config.InhibitRule `yaml:"inhibit_rules,omitempty" json:"inhibit_rules,omitempty"`
	Receivers    []*Receiver           `yaml:"receivers,omitempty" json:"receivers,omitempty"`
	Templates    []string              `yaml:"templates" json:"templates"`
}

// GetReceiver returns the receiver config with the given name
func (c *Config) GetReceiver(name string) *Receiver {
	for _, rec := range c.Receivers {
		if rec.Name == name {
			return rec
		}
	}
	return nil
}

func (c *Config) GetRouteIdx(name string) int {
	for idx, route := range c.Route.Routes {
		if route.Receiver == name {
			return idx
		}
	}
	return -1
}

func (c *Config) initializeNetworkBaseRoute(route *config.Route, networkID string) error {
	baseRouteName := makeBaseRouteName(networkID)
	if c.GetReceiver(baseRouteName) != nil {
		return fmt.Errorf("Base route for network %s already exists", networkID)
	}

	c.Receivers = append(c.Receivers, &Receiver{Name: baseRouteName})
	route.Receiver = baseRouteName
	route.Match = map[string]string{exporters.NetworkLabelNetwork: networkID}

	c.Route.Routes = append(c.Route.Routes, route)

	return c.Validate()
}

// Validate makes sure that the config is properly formed. Have to do this here
// since alertmanager only does validation during unmarshaling
func (c *Config) Validate() error {
	receiverNames := map[string]struct{}{}

	for _, rcv := range c.Receivers {
		if _, ok := receiverNames[rcv.Name]; ok {
			return fmt.Errorf("notification receiver name '%s' is not unique", rcv.Name)
		}
		for _, sc := range rcv.SlackConfigs {
			err := validateURL(sc.APIURL)
			if err != nil {
				return err
			}
		}
		receiverNames[rcv.Name] = struct{}{}
	}
	if c.Route == nil {
		return fmt.Errorf("no route provided")
	}
	if len(c.Route.Receiver) == 0 {
		return fmt.Errorf("root route must specify a default receiver")
	}
	if len(c.Route.Match) > 0 || len(c.Route.MatchRE) > 0 {
		return fmt.Errorf("root route must not have any matchers")
	}

	// check that all receivers used in routing tree are defined
	return checkReceiver(c.Route, receiverNames)
}

func validateURL(url string) error {
	if !strings.HasPrefix(url, "http") {
		return fmt.Errorf("invalid url: %s", url)
	}
	return nil
}

// checkReceiver returns an error if a node in the routing tree
// references a receiver not in the given map.
func checkReceiver(r *config.Route, receivers map[string]struct{}) error {
	for _, sr := range r.Routes {
		if err := checkReceiver(sr, receivers); err != nil {
			return err
		}
	}
	if r.Receiver == "" {
		return nil
	}
	if _, ok := receivers[r.Receiver]; !ok {
		return fmt.Errorf("undefined receiver %q used in route", r.Receiver)
	}
	return nil
}

// Receiver uses custom notifier configs to allow for marshaling of secrets.
type Receiver struct {
	Name string `yaml:"name" json:"name"`

	SlackConfigs []*SlackConfig `yaml:"slack_configs,omitempty" json:"slack_configs,omitempty"`
}

// Secure replaces the receiver's name with a networkID prefix
func (r *Receiver) Secure(networkID string) {
	r.Name = secureReceiverName(r.Name, networkID)
}

// Unsecure removes the networkID prefix from the receiver name
func (r *Receiver) Unsecure(networkID string) {
	r.Name = unsecureReceiverName(r.Name, networkID)
}

func secureReceiverName(name, networkID string) string {
	return receiverNetworkPrefix(networkID) + name
}

func unsecureReceiverName(name, networkID string) string {
	if strings.HasPrefix(name, receiverNetworkPrefix(networkID)) {
		return name[len(receiverNetworkPrefix(networkID)):]
	}
	return name
}

// SlackConfig uses string instead of SecretURL for the APIURL field so that it
// is marshaled as is instead of being obscured which is how alertmanager handles
// secrets
type SlackConfig struct {
	APIURL      string                `yaml:"api_url" json:"api_url"`
	Channel     string                `yaml:"channel" json:"channel"`
	Username    string                `yaml:"username" json:"username"`
	Color       string                `yaml:"color,omitempty" json:"color,omitempty"`
	Title       string                `yaml:"title,omitempty" json:"title,omitempty"`
	TitleLink   string                `yaml:"title_link,omitempty" json:"title_link,omitempty"`
	Pretext     string                `yaml:"pretext,omitempty" json:"pretext,omitempty"`
	Text        string                `yaml:"text,omitempty" json:"text,omitempty"`
	Fields      []*config.SlackField  `yaml:"fields,omitempty" json:"fields,omitempty"`
	ShortFields bool                  `yaml:"short_fields,omitempty" json:"short_fields,omitempty"`
	Footer      string                `yaml:"footer,omitempty" json:"footer,omitempty"`
	Fallback    string                `yaml:"fallback,omitempty" json:"fallback,omitempty"`
	CallbackID  string                `yaml:"callback_id,omitempty" json:"callback_id,omitempty"`
	IconEmoji   string                `yaml:"icon_emoji,omitempty" json:"icon_emoji,omitempty"`
	IconURL     string                `yaml:"icon_url,omitempty" json:"icon_url,omitempty"`
	ImageURL    string                `yaml:"image_url,omitempty" json:"image_url,omitempty"`
	ThumbURL    string                `yaml:"thumb_url,omitempty" json:"thumb_url,omitempty"`
	LinkNames   bool                  `yaml:"link_names,omitempty" json:"link_names,omitempty"`
	Actions     []*config.SlackAction `yaml:"actions,omitempty" json:"actions,omitempty"`
}

// RouteJSONWrapper Provides a struct to marshal/unmarshal into a rulefmt.Rule
// since rulefmt does not support json encoding
type RouteJSONWrapper struct {
	Receiver string `yaml:"receiver,omitempty" json:"receiver,omitempty"`

	GroupByStr []string          `yaml:"group_by,omitempty" json:"group_by,omitempty"`
	GroupBy    []model.LabelName `yaml:"-" json:"-"`
	GroupByAll bool              `yaml:"-" json:"-"`

	Match    map[string]string        `yaml:"match,omitempty" json:"match,omitempty"`
	MatchRE  map[string]config.Regexp `yaml:"match_re,omitempty" json:"match_re,omitempty"`
	Continue bool                     `yaml:"continue,omitempty" json:"continue,omitempty"`
	Routes   []*RouteJSONWrapper      `yaml:"routes,omitempty" json:"routes,omitempty"`

	GroupWait      string `yaml:"group_wait,omitempty" json:"group_wait,omitempty"`
	GroupInterval  string `yaml:"group_interval,omitempty" json:"group_interval,omitempty"`
	RepeatInterval string `yaml:"repeat_interval,omitempty" json:"repeat_interval,omitempty"`
}

// ToPrometheusConfig converts a json-compatible route specification to a
// prometheus route config
func (r *RouteJSONWrapper) ToPrometheusConfig() (config.Route, error) {
	var groupWait, groupInterval, repeatInterval model.Duration
	var groupWaitP, groupIntervalP, repeatIntervalP *model.Duration
	var err error

	if r.GroupWait != "" {
		groupWait, err = model.ParseDuration(r.GroupWait)
		if err != nil {
			return config.Route{}, fmt.Errorf("Invalid GroupWait '%s': %v", r.GroupWait, err)
		}
	}
	if r.GroupInterval != "" {
		groupInterval, err = model.ParseDuration(r.GroupInterval)
		if err != nil {
			return config.Route{}, fmt.Errorf("Invalid GroupInterval '%s': %v", r.GroupInterval, err)
		}
		if time.Duration(groupInterval) == time.Duration(0) {
			return config.Route{}, fmt.Errorf("GroupInterval cannot be 0")
		}
	}
	if r.RepeatInterval != "" {
		repeatInterval, err = model.ParseDuration(r.RepeatInterval)
		if err != nil {
			return config.Route{}, fmt.Errorf("Invalid RepeatInterval '%s': %v", r.RepeatInterval, err)
		}
		if time.Duration(repeatInterval) == time.Duration(0) {
			return config.Route{}, fmt.Errorf("RepeatInterval cannot be 0")
		}
	}
	groupWaitP = &groupWait
	if time.Duration(groupInterval) != time.Duration(0) {
		groupIntervalP = &groupInterval
	}
	if time.Duration(repeatInterval) != time.Duration(0) {
		repeatIntervalP = &repeatInterval
	}

	var configRoutes []*config.Route
	for _, childRoute := range r.Routes {
		route, err := childRoute.ToPrometheusConfig()
		if err != nil {
			return config.Route{}, fmt.Errorf("error converting child route: %v", err)
		}
		configRoutes = append(configRoutes, &route)
	}

	route := config.Route{
		Receiver:       r.Receiver,
		GroupByStr:     r.GroupByStr,
		GroupBy:        r.GroupBy,
		GroupByAll:     r.GroupByAll,
		Match:          r.Match,
		MatchRE:        r.MatchRE,
		Continue:       r.Continue,
		Routes:         configRoutes,
		GroupWait:      groupWaitP,
		GroupInterval:  groupIntervalP,
		RepeatInterval: repeatIntervalP,
	}
	return route, nil
}
