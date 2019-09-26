/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package plugin_test

import (
	"testing"

	"magma/cwf/cloud/go/cwf"
	"magma/cwf/cloud/go/plugin"
	"magma/cwf/cloud/go/plugin/models"
	fegmconfig "magma/feg/cloud/go/protos/mconfig"
	ltemconfig "magma/lte/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_Build(t *testing.T) {
	builder := &plugin.Builder{}

	// empty case: no cwf associated to magmad gateway
	nw := configurator.Network{ID: "n1"}
	gw := configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "gw1"}
	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{gw},
	}

	actual := map[string]proto.Message{}
	expected := map[string]proto.Message{}
	err := builder.Build("n1", "gw1", graph, nw, actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// Network config exists
	nw.Configs = map[string]interface{}{
		cwf.CwfNetworkType: defaultConfig,
	}
	actual = map[string]proto.Message{}
	expected = map[string]proto.Message{
		"eap_aka": &fegmconfig.EapAkaConfig{LogLevel: 1,
			Timeout: &fegmconfig.EapAkaConfig_Timeouts{
				ChallengeMs:            20000,
				ErrorNotificationMs:    10000,
				SessionMs:              43200000,
				SessionAuthenticatedMs: 5000,
			},
			PlmnIds: []string{},
		},
		"aaa_server": &fegmconfig.AAAConfig{LogLevel: 1,
			IdleSessionTimeoutMs: 21600000,
			AccountingEnabled:    false,
			CreateSessionOnAuth:  false,
		},
		"pipelined": &ltemconfig.PipelineD{
			LogLevel:      protos.LogLevel_INFO,
			UeIpBlock:     "192.168.128.0/24", // Unused by CWF
			NatEnabled:    false,
			DefaultRuleId: "",
			RelayEnabled:  true,
			Services: []ltemconfig.PipelineD_NetworkServices{
				ltemconfig.PipelineD_ENFORCEMENT,
			},
		},
		"sessiond": &ltemconfig.SessionD{
			LogLevel:     protos.LogLevel_INFO,
			RelayEnabled: true,
		},
		"redirectd": &ltemconfig.RedirectD{
			LogLevel: protos.LogLevel_INFO,
		},
	}
	err = builder.Build("n1", "gw1", graph, nw, actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

var defaultConfig = &models.NetworkCarrierWifiConfigs{
	EapAka: &models.EapAka{
		Timeout: &models.EapAkaTimeout{
			ChallengeMs:            20000,
			ErrorNotificationMs:    10000,
			SessionMs:              43200000,
			SessionAuthenticatedMs: 5000,
		},
		PlmnIds: []string{},
	},
	AaaServer: &models.AaaServer{
		IDLESessionTimeoutMs: 21600000,
		AccountingEnabled:    false,
		CreateSessionOnAuth:  false,
	},
	NetworkServices: []string{"policy_enforcement"},
	DefaultRuleID:   swag.String(""),
}
