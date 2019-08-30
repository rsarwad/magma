/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package plugin

import (
	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/plugin/handlers"
	models4 "magma/lte/cloud/go/plugin/models"
	cellularh "magma/lte/cloud/go/services/cellular/obsidian/handlers"
	cellularState "magma/lte/cloud/go/services/cellular/state"
	meteringdh "magma/lte/cloud/go/services/meteringd_records/obsidian/handlers"
	policydbh "magma/lte/cloud/go/services/policydb/obsidian/handlers"
	models2 "magma/lte/cloud/go/services/policydb/obsidian/models"
	policydbstreamer "magma/lte/cloud/go/services/policydb/streamer"
	"magma/lte/cloud/go/services/subscriberdb"
	subscriberdbh "magma/lte/cloud/go/services/subscriberdb/obsidian/handlers"
	models3 "magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	subscriberdbstreamer "magma/lte/cloud/go/services/subscriberdb/streamer"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/serde"
	srvconfig "magma/orc8r/cloud/go/service/config"
	"magma/orc8r/cloud/go/service/serviceregistry"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/streamer/providers"
)

// LteOrchestratorPlugin implements OrchestratorPlugin for the LTE module
type LteOrchestratorPlugin struct{}

func (*LteOrchestratorPlugin) GetName() string {
	return lte.ModuleName
}

func (*LteOrchestratorPlugin) GetServices() []registry.ServiceLocation {
	serviceLocations, err := serviceregistry.LoadServiceRegistryConfig(lte.ModuleName)
	if err != nil {
		return []registry.ServiceLocation{}
	}
	return serviceLocations
}

func (*LteOrchestratorPlugin) GetSerdes() []serde.Serde {
	return []serde.Serde{
		// TODO: expose enodeb state via swagger model and change serde to swagger serde
		&cellularState.EnodebStateSerde{},
		state.NewStateSerde(lte.SubscriberStateType, &models3.SubscriberState{}),

		// Configurator serdes
		configurator.NewNetworkConfigSerde(lte.CellularNetworkType, &models4.NetworkCellularConfigs{}),
		configurator.NewNetworkEntityConfigSerde(lte.CellularGatewayType, &models4.GatewayCellularConfigs{}),
		configurator.NewNetworkEntityConfigSerde(lte.CellularEnodebType, &models4.EnodebConfiguration{}),

		configurator.NewNetworkEntityConfigSerde(lte.PolicyRuleEntityType, &models2.PolicyRule{}),
		configurator.NewNetworkEntityConfigSerde(lte.BaseNameEntityType, &models2.BaseNameRecord{}),
		configurator.NewNetworkEntityConfigSerde(subscriberdb.EntityType, &models4.LteSubscription{}),
	}
}

func (*LteOrchestratorPlugin) GetMconfigBuilders() []configurator.MconfigBuilder {
	return []configurator.MconfigBuilder{
		&Builder{},
	}
}

func (*LteOrchestratorPlugin) GetMetricsProfiles(metricsConfig *srvconfig.ConfigMap) []metricsd.MetricsProfile {
	return []metricsd.MetricsProfile{}
}

func (*LteOrchestratorPlugin) GetObsidianHandlers(metricsConfig *srvconfig.ConfigMap) []obsidian.Handler {
	return plugin.FlattenHandlerLists(
		cellularh.GetObsidianHandlers(),
		meteringdh.GetObsidianHandlers(),
		policydbh.GetObsidianHandlers(),
		subscriberdbh.GetObsidianHandlers(),
		handlers.GetHandlers(),
	)
}

func (*LteOrchestratorPlugin) GetStreamerProviders() []providers.StreamProvider {
	return []providers.StreamProvider{
		&subscriberdbstreamer.SubscribersProvider{},
		&policydbstreamer.PoliciesProvider{},
		&policydbstreamer.BaseNamesProvider{},
	}
}
