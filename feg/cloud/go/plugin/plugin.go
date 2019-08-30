/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package plugin exposes the OrchestratorPlugin implementation for the module.
// This is so unit tests can register the plugin without building and loading
// it from disk.
package plugin

import (
	"magma/feg/cloud/go/feg"
	fegh "magma/feg/cloud/go/services/controller/obsidian/handlers"
	"magma/feg/cloud/go/services/controller/obsidian/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/serde"
	srvconfig "magma/orc8r/cloud/go/service/config"
	"magma/orc8r/cloud/go/service/serviceregistry"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/streamer/providers"
)

// FegOrchestratorPlugin is an implementation of OrchestratorPlugin for the
// feg module
type FegOrchestratorPlugin struct{}

func (*FegOrchestratorPlugin) GetName() string {
	return feg.ModuleName
}

func (*FegOrchestratorPlugin) GetServices() []registry.ServiceLocation {
	serviceLocations, err := serviceregistry.LoadServiceRegistryConfig(feg.ModuleName)
	if err != nil {
		return []registry.ServiceLocation{}
	}
	return serviceLocations
}

func (*FegOrchestratorPlugin) GetSerdes() []serde.Serde {
	return []serde.Serde{
		// configurator serdes
		configurator.NewNetworkConfigSerde(feg.FegNetworkType, &models.NetworkFederationConfigs{}),
		configurator.NewNetworkEntityConfigSerde(feg.FegGatewayType, &models.GatewayFegConfigs{}),
	}
}

func (*FegOrchestratorPlugin) GetMconfigBuilders() []configurator.MconfigBuilder {
	return []configurator.MconfigBuilder{
		&Builder{},
	}
}

func (*FegOrchestratorPlugin) GetMetricsProfiles(metricsConfig *srvconfig.ConfigMap) []metricsd.MetricsProfile {
	return []metricsd.MetricsProfile{}
}

func (*FegOrchestratorPlugin) GetObsidianHandlers(metricsConfig *srvconfig.ConfigMap) []obsidian.Handler {
	return plugin.FlattenHandlerLists(
		fegh.GetObsidianHandlers(),
	)
}

func (*FegOrchestratorPlugin) GetStreamerProviders() []providers.StreamProvider {
	return []providers.StreamProvider{}
}
