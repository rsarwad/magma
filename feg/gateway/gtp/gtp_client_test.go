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

package gtp

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wmnsk/go-gtp/gtpv2"
	"github.com/wmnsk/go-gtp/gtpv2/ie"
	"github.com/wmnsk/go-gtp/gtpv2/message"

	"magma/feg/cloud/go/protos"
)

const (
	gtpTimeout = 200 * time.Millisecond
	sgwAddrs   = ":0"
	pgwAddrs   = "127.0.0.1:0" //port 0 means go will choose. Selected port will be injected on getDefaultConfig
	IMSI1      = "123456789012345"
	bearerId1  = uint8(5)
	bearerId2  = uint8(6)
	cSgwTeid   = uint32(10)
	cPgwTeid   = uint32(22)
	uPgwTeid   = uint32(33)

	qci1 = uint8(5)
	qci2 = uint8(6)
)

func TestEcho(t *testing.T) {
	// run GTP server (PGW)
	pgwCli := startGTPServer(t)
	actualServerIPAndPort := pgwCli.LocalAddr().String()

	// run GTP client (SGW) and send echo message to check if server is available
	_, err := NewConnectedAutoClient(context.Background(), actualServerIPAndPort, gtpv2.IFTypeS5S8SGWGTPC, gtpTimeout)

	// if no error service was started and echo was received properly
	assert.NoError(t, err)
}

// TestGtpClient tests gtp client-server interaction using dummy handlers
func TestGtpClient(t *testing.T) {
	// run GTP server (PGW)
	gtpServer := startGTPServer(t)
	actualServerIPAndPort := gtpServer.LocalAddr().String()

	// run GTP client (SGW) but do not attach to any server
	gtpClient, err := NewRunningClient(context.Background(), sgwAddrs, gtpv2.IFTypeS5S8SGWGTPC, gtpTimeout)
	assert.NoError(t, err)

	remoteAddr, err := net.ResolveUDPAddr("udp", actualServerIPAndPort)
	assert.NoError(t, err)

	// find out the local interface to be used (because it is not specified for testing in sgwAddrs)
	localIP, err := GetLocalOutboundIP(remoteAddr)
	assert.NoError(t, err)

	// add a dummy handler at the server for create session request
	gtpServer.AddHandlers(map[uint8]gtpv2.HandlerFunc{
		message.MsgTypeCreateSessionRequest: getHandleCreateSessionRequest(actualServerIPAndPort, bearerId1),
	})

	// add a dummy handler at tlient client for create session response
	gtpClient.AddHandlers(map[uint8]gtpv2.HandlerFunc{
		message.MsgTypeCreateSessionResponse: getHandleCreateSessionResponse(gtpClient),
	})

	csr := getCreateSessionRequest(t, gtpClient, localIP, actualServerIPAndPort, bearerId1, qci1)
	msg := message.NewCreateSessionRequest(0, 0, csr...)
	resMsg, err := gtpClient.SendMessageAndExtractGrpc(IMSI1, cSgwTeid, remoteAddr, msg)
	assert.NoError(t, err)
	csRes := resMsg.(*protos.CreateSessionResponsePgw)
	assert.NotEmpty(t, csRes)
	assert.Equal(t, cPgwTeid, csRes.CPgwFteid.Teid)

	// create same session with differnt QCI
	csr = getCreateSessionRequest(t, gtpClient, localIP, actualServerIPAndPort, bearerId1, qci2)

	msg = message.NewCreateSessionRequest(0, 0, csr...)
	resMsg, err = gtpClient.SendMessageAndExtractGrpc(IMSI1, cSgwTeid, remoteAddr, msg)
	assert.NoError(t, err)
	csRes = resMsg.(*protos.CreateSessionResponsePgw)
	assert.NotEmpty(t, csRes)
	assert.Equal(t, cPgwTeid, csRes.CPgwFteid.Teid)
}

func startGTPServer(t *testing.T) *Client {
	pgwConn, err := NewRunningClient(context.Background(), pgwAddrs, gtpv2.IFTypeS5S8PGWGTPC, gtpTimeout)
	assert.NoError(t, err)
	return pgwConn
}

func getCreateSessionRequest(t *testing.T, cli *Client, laddrs net.IP, raddrs string, bearerId, qci uint8) []*ie.IE {
	// SGW control plane teid
	cSgwFTeid := ie.NewFullyQualifiedTEID(gtpv2.IFTypeS5S8SGWGTPC, cSgwTeid, raddrs, "").WithInstance(0)

	// SGW user plane teid
	uSgwFTeid := ie.NewFullyQualifiedTEID(gtpv2.IFTypeS5S8SGWGTPU, 11, raddrs, "").WithInstance(2)

	return []*ie.IE{
		ie.NewIMSI(IMSI1),
		ie.NewMSISDN("8130900000005"),
		ie.NewMobileEquipmentIdentity("123456780000015"),
		ie.NewUserLocationInformation(
			0, 0, 0, 1, 1, 0, 0, 0,
			"123", "456", 0, 0, 0, 0, 1, 1, 0, 0,
		),
		ie.NewRATType(gtpv2.RATTypeEUTRAN),
		ie.NewIndicationFromOctets(0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00),
		cSgwFTeid,
		ie.NewAccessPointName("internet"),
		ie.NewSelectionMode(gtpv2.SelectionModeMSorNetworkProvidedAPNSubscribedVerified),
		ie.NewPDNType(gtpv2.PDNTypeIPv4),
		ie.NewAPNRestriction(gtpv2.APNRestrictionNoExistingContextsorRestriction),
		ie.NewAggregateMaximumBitRate(0, 0),
		ie.NewBearerContext(
			ie.NewEPSBearerID(bearerId),
			ie.NewBearerQoS(1, 2, 3, qci, 40000, 4001, 4000, 4001),
			uSgwFTeid,
		),
	}
}

// getHandleCreateSessionRequest dummy create sesson request handler
func getHandleCreateSessionRequest(pgwAddrs string, bearerId uint8) gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, sgwAddr net.Addr, msg message.Message) error {

		cPgwFTeid := ie.NewFullyQualifiedTEID(gtpv2.IFTypeS5S8PGWGTPC, cPgwTeid, pgwAddrs, "").WithInstance(1)
		uPgwFTeid := ie.NewFullyQualifiedTEID(gtpv2.IFTypeS5S8PGWGTPU, uPgwTeid, pgwAddrs, "").WithInstance(2)

		csRspFromPGW := message.NewCreateSessionResponse(
			cSgwTeid, msg.Sequence(),
			ie.NewCause(gtpv2.CauseRequestAccepted, 0, 0, 0, nil),
			cPgwFTeid,
			ie.NewAPNRestriction(gtpv2.APNRestrictionPublic2),
			ie.NewBearerContext(
				ie.NewCause(gtpv2.CauseRequestAccepted, 0, 0, 0, nil),
				ie.NewEPSBearerID(bearerId),
				uPgwFTeid,
			))
		if err := c.RespondTo(sgwAddr, msg, csRspFromPGW); err != nil {
			return err
		}
		return nil
	}
}

// getHandleCreateSessionResponse dummy create session response
func getHandleCreateSessionResponse(cli *Client) gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, pgwAddr net.Addr, msg message.Message) error {

		csResGtp := msg.(*message.CreateSessionResponse)
		csRes := &protos.CreateSessionResponsePgw{}
		if pgwCFteidIE := csResGtp.PGWS5S8FTEIDC; pgwCFteidIE != nil {
			teid, err := pgwCFteidIE.TEID()
			if err != nil {
				return err
			}
			csRes.CPgwFteid = &protos.Fteid{
				Ipv4Address: "",
				Ipv6Address: "",
				Teid:        teid,
			}
		} else {
			return &gtpv2.RequiredIEMissingError{Type: ie.FullyQualifiedTEID}
		}

		return cli.PassMessage(msg.TEID(), pgwAddr, csResGtp, csRes, nil)

	}
}
