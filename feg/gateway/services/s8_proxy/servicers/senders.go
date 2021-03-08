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

// senders contains the calls that will be run when a GTP command is sent.
// Those functions will also return the result of the call

package servicers

import (
	"fmt"
	"net"
	"time"

	"magma/feg/cloud/go/protos"

	"github.com/golang/glog"
	"github.com/wmnsk/go-gtp/gtpv2/message"
)

// sendAndReceiveCreateSession creates a session in the gtp client, sends the create session request
// to PGW and waits for its answers.
// Returns a GRPC message translaged from the GTP-U create session response
func (s *S8Proxy) sendAndReceiveCreateSession(
	csReq *protos.CreateSessionRequestPgw,
	cPgwUDPAddr *net.UDPAddr,
	csReqMsg message.Message) (*protos.CreateSessionResponsePgw, error) {
	glog.V(2).Infof("Send Create Session Request (gtp) to %s:\n%s",
		cPgwUDPAddr.String(), csReq.String())

	grpcMessage, err := s.gtpClient.SendMessageAndExtractGrpc(csReq.Imsi, csReq.CAgwTeid, cPgwUDPAddr, csReqMsg)
	/*
		session, seq, err := s.gtpClient.CreateSession(cPgwUDPAddr, csReqIEs...)



		if err != nil {
			return nil, fmt.Errorf("failed to send create session at %s: %s", cPgwUDPAddr.String(), err)
		}

		// add TEID to session and register session
		session.AddTEID(sessionTeids.uAgwFTeid.MustInterfaceType(), sessionTeids.uAgwFTeid.MustTEID())
		s.gtpClient.RegisterSession(sessionTeids.cFegFTeid.MustTEID(), session)

		grpcMessage, err := waitMessageAndExtractGrpc(session, seq)
	*/
	if err != nil {
		return nil, fmt.Errorf("no response message to CreateSessionRequest: %s", err)
	}

	// check if message is proper
	csRes, ok := grpcMessage.(*protos.CreateSessionResponsePgw)
	if !ok {
		s.gtpClient.RemoveSessionByIMSI(csReq.Imsi)
		//s.gtpClient.RemoveSession(session)
		return nil, fmt.Errorf("Wrong response type (no CreateSessionResponse), maybe received out of order response message: %s", err)
	}
	// TODO : Delete
	glog.V(2).Infof("Create Session Response (grpc):\n%s", csRes.String())
	return csRes, nil
}

// sendAndReceiveDeleteSession  sends delete session request GTP-U message to PGW and
// waits for its answers.
// Returns a GRPC message translaged from the GTP-U create session response
func (s *S8Proxy) sendAndReceiveDeleteSession(req *protos.DeleteSessionRequestPgw,
	cPgwUDPAddr *net.UDPAddr,
	dsReqMsg message.Message) (*protos.DeleteSessionResponsePgw, error) {

	//seq, err := s.gtpClient.DeleteSession(teid, session, dsReqIEs...)
	glog.V(2).Infof("Send Delete Session Request (gtp) to %s:\n%s", cPgwUDPAddr,
		dsReqMsg)
	grpcMessage, err := s.gtpClient.SendMessageAndExtractGrpc(req.Imsi, req.CAgwTeid, cPgwUDPAddr, dsReqMsg)
	//if err != nil {
	//	return nil, err
	//}
	//grpcMessage, err := waitMessageAndExtractGrpc(session, seq)
	if err != nil {
		return nil, fmt.Errorf("no response message to DeleteSessionRequest: %s", err)
	}
	dsRes, ok := grpcMessage.(*protos.DeleteSessionResponsePgw)
	if !ok {
		return nil, fmt.Errorf("Wrong response type (no DeleteSessionResponse), maybe received out of order response message: %s", err)
	}
	glog.V(2).Infof("Delete Session Response (grpc):\n%s", dsRes.String())
	return dsRes, err
}

func (s *S8Proxy) sendAndReceiveEchoRequest(cPgwUDPAddr *net.UDPAddr) error {
	_, err := s.gtpClient.Conn.EchoRequest(cPgwUDPAddr)
	if err != nil {
		return err
	}
	select {
	case res := <-s.echoChannel:
		return res
	case <-time.After(s.gtpClient.GtpTimeout):
		return fmt.Errorf("waitEchoResponse timeout")
	}
}
