/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package mock_ocs

import (
	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/session_proxy/credit_control"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/golang/glog"
)

type ccrMessage struct {
	SessionID        datatype.UTF8String       `avp:"Session-Id"`
	OriginHost       datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm      datatype.DiameterIdentity `avp:"Origin-Realm"`
	DestinationRealm datatype.DiameterIdentity `avp:"Destination-Realm"`
	DestinationHost  datatype.DiameterIdentity `avp:"Destination-Host"`
	RequestType      datatype.Enumerated       `avp:"CC-Request-Type"`
	RequestNumber    datatype.Unsigned32       `avp:"CC-Request-Number"`
	MSCC             []*ccrCredit              `avp:"Multiple-Services-Credit-Control"`
	SubscriptionIDs  []*subscriptionID         `avp:"Subscription-Id"`
}

type subscriptionID struct {
	IDType credit_control.SubscriptionIDType `avp:"Subscription-Id-Type"`
	IDData string                            `avp:"Subscription-Id-Data"`
}

type usedServiceUnit struct {
	InputOctets  uint64 `avp:"CC-Input-Octets"`
	OutputOctets uint64 `avp:"CC-Output-Octets"`
	TotalOctets  uint64 `avp:"CC-Total-Octets"`
}

type ccrCredit struct {
	RatingGroup     uint32           `avp:"Rating-Group"`
	UsedServiceUnit *usedServiceUnit `avp:"Used-Service-Unit"`
}

// getCCRHandler returns a handler to be called when the server receives a CCR
func getCCRHandler(srv *OCSDiamServer) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		glog.V(2).Infof("Received CCR from %s\n", c.RemoteAddr())
		var ccr ccrMessage
		if err := m.Unmarshal(&ccr); err != nil {
			glog.Errorf("Failed to unmarshal CCR %s", err)
			return
		}
		imsi := getIMSI(ccr)
		if len(imsi) == 0 {
			glog.Errorf("Could not find IMSI in CCR")
			sendAnswer(ccr, c, m, diam.AuthenticationRejected)
			return
		}
		account, found := srv.accounts[imsi]
		if !found {
			sendAnswer(ccr, c, m, diam.AuthenticationRejected)
			return
		}
		account.CurrentState = &SubscriberSessionState{
			Connection: c,
			SessionID:  string(ccr.SessionID),
		}

		if credit_control.CreditRequestType(ccr.RequestType) == credit_control.CRTTerminate {
			sendAnswer(ccr, c, m, diam.Success)
			return
		}

		creditAnswers := make([]*diam.AVP, 0, len(ccr.MSCC))
		for _, mscc := range ccr.MSCC {
			if mscc.UsedServiceUnit != nil {
				decrementUsedCredit(
					account.ChargingCredit[mscc.RatingGroup],
					mscc.UsedServiceUnit,
				)
			}
			returnOctets, final := getQuotaGrant(srv, account.ChargingCredit[mscc.RatingGroup])
			if returnOctets.GetTotalOctets() <= 0 {
				sendAnswer(ccr, c, m, DiameterCreditLimitReached)
				return
			}
			creditAnswers = append(creditAnswers, toGrantedUnitsAVP(
				mscc.RatingGroup,
				srv.ocsConfig.ValidityTime,
				returnOctets,
				final,
			))
		}

		sendAnswer(ccr, c, m, diam.Success, creditAnswers...)
	}
}

func decrementUsedCredit(credit *CreditBucket, usage *usedServiceUnit) {
	credit.Volume.TotalOctets = decrementOrZero(credit.Volume.GetTotalOctets(), usage.TotalOctets)
	credit.Volume.OutputOctets = decrementOrZero(credit.Volume.GetOutputOctets(), usage.OutputOctets)
	credit.Volume.InputOctets = decrementOrZero(credit.Volume.GetInputOctets(), usage.InputOctets)
}

func decrementOrZero(first, second uint64) uint64 {
	result := first - second
	if result < 0 {
		return 0
	}
	return result
}

// sendAnswer sends a CCA to the connection given
func sendAnswer(
	ccr ccrMessage,
	conn diam.Conn,
	message *diam.Message,
	statusCode uint32,
	additionalAVPs ...*diam.AVP,
) {
	a := message.Answer(statusCode)
	a.NewAVP(avp.OriginHost, avp.Mbit, 0, ccr.DestinationHost)
	a.NewAVP(avp.OriginRealm, avp.Mbit, 0, ccr.DestinationRealm)
	a.NewAVP(avp.DestinationRealm, avp.Mbit, 0, ccr.OriginRealm)
	a.NewAVP(avp.DestinationHost, avp.Mbit, 0, ccr.OriginHost)
	a.NewAVP(avp.CCRequestType, avp.Mbit, 0, ccr.RequestType)
	a.NewAVP(avp.CCRequestNumber, avp.Mbit, 0, ccr.RequestNumber)
	a.NewAVP(avp.SessionID, avp.Mbit, 0, ccr.SessionID)
	for _, avp := range additionalAVPs {
		a.InsertAVP(avp)
	}
	// SessionID must be the first AVP
	a.InsertAVP(diam.NewAVP(avp.SessionID, avp.Mbit, 0, ccr.SessionID))

	_, err := a.WriteTo(conn)
	if err != nil {
		glog.Errorf("Failed to write message to %s: %s\n%s\n",
			conn.RemoteAddr(), err, a)
		return
	}
	glog.V(2).Infof("Sent CCA to %s:\n", conn.RemoteAddr())
}

// getIMSI finds the account IMSI in a CCR message
func getIMSI(message ccrMessage) string {
	for _, subID := range message.SubscriptionIDs {
		if subID.IDType == credit_control.EndUserIMSI {
			return subID.IDData
		}
	}
	return ""
}

// getQuotaGrant gets how much credit to return in a CCA-update, which is the
// minimum between the max usage and how much credit is in the account
// Returns credits to return and true if these are the final bytes
func getQuotaGrant(srv *OCSDiamServer, bucket *CreditBucket) (*protos.Octets, bool) {
	var grant *protos.Octets
	var maxTotalUsage uint64

	switch bucket.Unit {
	case protos.CreditInfo_Bytes:
		maxUsage := srv.ocsConfig.MaxUsageOctets
		maxTotalUsage = maxUsage.GetTotalOctets()
		perRequest := bucket.Volume
		grant = &protos.Octets{
			TotalOctets:  getMin(maxUsage.GetTotalOctets(), perRequest.GetTotalOctets()),
			InputOctets:  getMin(maxUsage.GetInputOctets(), perRequest.GetInputOctets()),
			OutputOctets: getMin(maxUsage.GetOutputOctets(), perRequest.GetOutputOctets())}

	case protos.CreditInfo_Time:
		maxTotalUsage = uint64(srv.ocsConfig.MaxUsageTime)
		grant = &protos.Octets{TotalOctets: getMin(uint64(srv.ocsConfig.MaxUsageTime), bucket.Volume.GetTotalOctets())}
	}
	if grant.GetTotalOctets() <= maxTotalUsage {
		return grant, true
	}
	return grant, false
}

func getMin(first, second uint64) uint64 {
	if first > second {
		return second
	}
	return first
}

func toGrantedUnitsAVP(ratingGroup uint32, validityTime uint32, quotaGrant *protos.Octets, isFinalUnit bool) *diam.AVP {
	creditGroup := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.GrantedServiceUnit, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.CCTotalOctets, avp.Mbit, 0, datatype.Unsigned64(quotaGrant.GetTotalOctets())),
					diam.NewAVP(avp.CCInputOctets, avp.Mbit, 0, datatype.Unsigned64(quotaGrant.GetInputOctets())),
					diam.NewAVP(avp.CCOutputOctets, avp.Mbit, 0, datatype.Unsigned64(quotaGrant.GetOutputOctets())),
				},
			}),
			diam.NewAVP(avp.ValidityTime, avp.Mbit, 0, datatype.Unsigned32(validityTime)),
			diam.NewAVP(avp.RatingGroup, avp.Mbit, 0, datatype.Unsigned32(ratingGroup)),
		},
	}
	if isFinalUnit {
		creditGroup.AddAVP(
			diam.NewAVP(avp.FinalUnitIndication, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					// TODO support other final unit actions
					diam.NewAVP(avp.FinalUnitAction, avp.Mbit, 0, datatype.Enumerated(TerminateAction)),
				},
			}),
		)
	}
	return diam.NewAVP(avp.MultipleServicesCreditControl, avp.Mbit, 0, creditGroup)
}
