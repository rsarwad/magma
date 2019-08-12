/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streamer

import (
	"sort"

	"magma/lte/cloud/go/lte"
	protos2 "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
)

type SubscribersProvider struct{}

func (provider *SubscribersProvider) GetStreamName() string {
	return "subscriberdb"
}

func (provider *SubscribersProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	ent, err := configurator.LoadEntityForPhysicalID(gatewayId, configurator.EntityLoadCriteria{})
	if err != nil {
		return nil, err
	}

	subEnts, err := configurator.LoadAllEntitiesInNetwork(ent.NetworkID, lte.SubscriberEntityType, configurator.EntityLoadCriteria{LoadConfig: true})
	if err != nil {
		return nil, err
	}

	subProtos := make([]*protos2.SubscriberData, 0, len(subEnts))
	for _, sub := range subEnts {
		subdata := sub.Config.(*models.Subscriber)
		subProto := &protos2.SubscriberData{}
		err = subdata.ToMconfig(subProto)
		if err != nil {
			return nil, err
		}
		subProto.NetworkId = &protos.NetworkID{Id: ent.NetworkID}
		subProtos = append(subProtos, subProto)
	}
	return subscribersToUpdates(subProtos)
}

func subscribersToUpdates(subs []*protos2.SubscriberData) ([]*protos.DataUpdate, error) {
	ret := make([]*protos.DataUpdate, 0, len(subs))
	for _, sub := range subs {
		marshaledProto, err := proto.Marshal(sub)
		if err != nil {
			return nil, err
		}
		update := &protos.DataUpdate{Key: protos2.SidString(sub.Sid), Value: marshaledProto}
		ret = append(ret, update)
	}
	sort.Slice(ret, func(i, j int) bool { return ret[i].Key < ret[j].Key })
	return ret, nil
}
