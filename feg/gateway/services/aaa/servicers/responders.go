/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// package servcers implements WiFi AAA GRPC services
package servicers

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/aaa/events"
	"magma/feg/gateway/services/aaa/metrics"
	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/aaa/session_manager"
	"magma/gateway/directoryd"
	lteprotos "magma/lte/cloud/go/protos"
	orcprotos "magma/orc8r/lib/go/protos"
)

const (
	MinIMSILen = 10
	MaxIMSILen = 16
)

// AbortSession is a method of AbortSessionResponder service.
func (srv *accountingService) AbortSession(
	ctx context.Context, req *lteprotos.AbortSessionRequest) (*lteprotos.AbortSessionResult, error) {

	res := &lteprotos.AbortSessionResult{}
	if req == nil {
		return res, Errorf(codes.InvalidArgument, "Nil Request")
	}
	imsi := req.GetUserName()
	sctx := &protos.Context{
		Imsi:      req.GetUserName(),
		SessionId: req.GetSessionId(),
	}
	if len(imsi) < MinIMSILen || len(imsi) > MaxIMSILen {
		errMsg := fmt.Sprintf("Invalid IMSI: %s", imsi)
		if srv.config.GetEventLoggingEnabled() {
			events.LogSessionTerminationFailedEvent(sctx, events.AbortSession, errMsg)
		}
		return res, Errorf(codes.InvalidArgument, errMsg)
	}
	sid := srv.sessions.FindSession(imsi)
	if len(sid) == 0 {
		res.Code = lteprotos.AbortSessionResult_USER_NOT_FOUND
		res.ErrorMessage = fmt.Sprintf("Session for IMSI: %s is not found", imsi)
		glog.Error(res.ErrorMessage)
		if srv.config.GetEventLoggingEnabled() {
			events.LogSessionTerminationFailedEvent(sctx, events.AbortSession, res.ErrorMessage)
		}
		return res, nil
	}
	s := srv.sessions.GetSession(sid)
	if s == nil {
		res.Code = lteprotos.AbortSessionResult_SESSION_NOT_FOUND
		res.ErrorMessage = fmt.Sprintf("Session for Radius Session ID: %s and IMSI: %s is not found", sid, imsi)
		glog.Error(res.ErrorMessage)
		if srv.config.GetEventLoggingEnabled() {
			events.LogSessionTerminationFailedEvent(sctx, events.AbortSession, res.ErrorMessage)
		}
		return res, nil
	}
	s.Lock()
	sctx = proto.Clone(s.GetCtx()).(*protos.Context)
	asid := sctx.AcctSessionId
	s.Unlock()
	if len(req.GetSessionId()) > 0 &&
		len(asid) > 0 &&
		asid != req.GetSessionId() {

		res.Code = lteprotos.AbortSessionResult_SESSION_NOT_FOUND
		res.ErrorMessage = fmt.Sprintf(
			"Accounting Session ID Mismatch for RadSID %s and IMSI: %s. Requested: %s, recorded: %s",
			sid, imsi, req.GetSessionId(), asid)
		glog.Error(res.ErrorMessage)
		if srv.config.GetEventLoggingEnabled() {
			events.LogSessionTerminationFailedEvent(sctx, events.AbortSession, res.ErrorMessage)
		}
		return res, nil
	}
	if srv.config.GetAccountingEnabled() {
		// ? can potentially end a new, valid session
		req := &lteprotos.LocalEndSessionRequest{
			Sid: makeSID(imsi),
			Apn: sctx.GetApn(),
		}
		session_manager.EndSession(req)
		metrics.EndSession.WithLabelValues(sctx.GetApn(), metrics.DecorateIMSI(sctx.GetImsi())).Inc()
	} else {
		deleteRequest := &orcprotos.DeleteRecordRequest{
			Id: imsi,
		}
		directoryd.DeleteRecord(deleteRequest)
	}
	srv.sessions.RemoveSession(sid)

	err := srv.dae.Disconnect(sctx)
	if err != nil {
		res.Code = lteprotos.AbortSessionResult_RADIUS_SERVER_ERROR
		res.ErrorMessage = fmt.Sprintf(
			"Radius Disconnect Error: %v for IMSI: %s, Acct SID: %s, Radius SID: %s", err, imsi, asid, sid)
		glog.Error(res.ErrorMessage)
		if srv.config.GetEventLoggingEnabled() {
			events.LogSessionTerminationFailedEvent(sctx, events.AbortSession, res.ErrorMessage)
		}
		return res, nil
	} else if srv.config.GetEventLoggingEnabled() {
		events.LogSessionTerminationSucceededEvent(sctx, events.AbortSession)
	}
	return res, err
}

// TerminateRegistration is a method of SWx Gateway Service Responder service.
func (srv *accountingService) TerminateRegistration(
	ctx context.Context, req *fegprotos.RegistrationTerminationRequest) (*fegprotos.RegistrationAnswer, error) {

	res := &fegprotos.RegistrationAnswer{}
	if req == nil {
		return res, Errorf(codes.InvalidArgument, "Nil Request")
	}
	sctx := &protos.Context{
		Imsi:      req.GetUserName(),
		SessionId: req.GetSessionId()}
	imsi := req.GetUserName()
	if len(imsi) < MinIMSILen {
		errMsg := fmt.Sprintf("Invalid IMSI: %s", imsi)
		if srv.config.GetEventLoggingEnabled() {
			events.LogSessionTerminationFailedEvent(sctx, events.RegistrationTermination, errMsg)
		}
		return res, Errorf(codes.InvalidArgument, errMsg)
	}
	sid := srv.sessions.FindSession(imsi)
	if len(sid) == 0 {
		errMsg := fmt.Sprintf("Session for IMSI: %s is not found", imsi)
		if srv.config.GetEventLoggingEnabled() {
			events.LogSessionTerminationFailedEvent(sctx, events.RegistrationTermination, errMsg)
		}
		return res, Errorf(codes.NotFound, errMsg)
	}
	s := srv.sessions.GetSession(sid)
	if s == nil {
		errMsg := fmt.Sprintf("Session for RadSID: %s and IMSI: %s is not found", sid, imsi)
		if srv.config.GetEventLoggingEnabled() {
			events.LogSessionTerminationFailedEvent(sctx, events.RegistrationTermination, errMsg)
		}
		return res, Errorf(codes.Internal, errMsg)
	}
	s.Lock()
	sctx = proto.Clone(s.GetCtx()).(*protos.Context)
	authSid := sctx.AuthSessionId
	acctSid := sctx.AcctSessionId
	s.Unlock()
	if len(req.GetSessionId()) > 0 &&
		(len(authSid) > 0 || len(acctSid) > 0) &&
		(authSid != req.GetSessionId() && acctSid != req.GetSessionId()) {
		errMsg := fmt.Sprintf("Accounting Session ID Mismatch for RadSID %s and IMSI: %s. Requested: %s, recorded: auth: %s | acct: %s",
			sid, imsi, req.GetSessionId(), authSid, acctSid)
		if srv.config.GetEventLoggingEnabled() {
			events.LogSessionTerminationFailedEvent(sctx, events.RegistrationTermination, errMsg)
		}
		return res, Errorf(codes.FailedPrecondition, errMsg)

	}
	deleteRequest := &orcprotos.DeleteRecordRequest{
		Id: imsi,
	}
	directoryd.DeleteRecord(deleteRequest) // remove it from directoryd even if session manager will try to remove it again

	if srv.config.GetAccountingEnabled() {
		// ? can potentially end a new, valid session
		sid := makeSID(imsi)
		req := &lteprotos.LocalEndSessionRequest{
			Sid: sid,
			Apn: sctx.GetApn(),
		}
		session_manager.EndSession(req)
		metrics.EndSession.WithLabelValues(sctx.GetApn(), sid.Id).Inc()
	}

	srv.sessions.RemoveSession(sid)

	err := srv.dae.Disconnect(sctx)
	if err != nil {
		if srv.config.GetEventLoggingEnabled() {
			events.LogSessionTerminationFailedEvent(sctx, events.RegistrationTermination, err.Error())
		}
		err = Error(codes.Internal, err)
	} else if srv.config.GetEventLoggingEnabled() {
		events.LogSessionTerminationSucceededEvent(sctx, events.RegistrationTermination)
	}
	return res, err
}
