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

package servicers

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/golang/glog"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/gtp"
)

type echoResponse struct {
	error
}

type S8Proxy struct {
	config      *S8ProxyConfig
	gtpClient   *gtp.Client
	echoChannel chan (error)
}

type S8ProxyConfig struct {
	GtpTimeout time.Duration
	ClientAddr string
	ServerAddr *net.UDPAddr
}

// NewS8Proxy creates an s8 proxy, but does not checks the PGW is alive
func NewS8Proxy(config *S8ProxyConfig) (*S8Proxy, error) {
	gtpCli, err := gtp.NewRunningClient(
		context.Background(), config.ClientAddr,
		gtp.SGWControlPlaneIfType, config.GtpTimeout)
	if err != nil {
		return nil, fmt.Errorf("Error creating S8_Proxy: %s", err)
	}
	return newS8ProxyImp(gtpCli, config)
}

// NewS8ProxyWithEcho creates an s8 proxy already connected to a server (checks with echo if PGW is alive)
// Used mainly for testing with s8_cli
func NewS8ProxyWithEcho(config *S8ProxyConfig) (*S8Proxy, error) {
	gtpCli, err := gtp.NewConnectedAutoClient(
		context.Background(), config.ServerAddr.String(),
		gtp.SGWControlPlaneIfType, config.GtpTimeout)
	if err != nil {
		return nil, fmt.Errorf("Error creating S8_Proxy: %s", err)
	}
	return newS8ProxyImp(gtpCli, config)
}

func newS8ProxyImp(cli *gtp.Client, config *S8ProxyConfig) (*S8Proxy, error) {
	// TODO: validate config
	s8p := &S8Proxy{
		config:      config,
		gtpClient:   cli,
		echoChannel: make(chan error),
	}
	addS8GtpHandlers(s8p)
	return s8p, nil
}

func (s *S8Proxy) CreateSession(ctx context.Context, req *protos.CreateSessionRequestPgw) (*protos.CreateSessionResponsePgw, error) {
	cPgwUDPAddr, err := s.configOrRequestedPgwAddress(req.PgwAddrs)
	if err != nil {
		err = fmt.Errorf("Create Session Request failed due to missing server address: %s", err)
		glog.Error(err)
		return nil, err
	}

	// build csReq IE message
	csReqMsg, err := buildCreateSessionRequestMsg(cPgwUDPAddr, req)
	if err != nil {
		return nil, err
	}

	// send, register and receive create session (session is created on the gtp client during this process too)
	csRes, err := s.sendAndReceiveCreateSession(req, cPgwUDPAddr, csReqMsg)
	if err != nil {
		err = fmt.Errorf("Create Session Request failed: %s", err)
		glog.Error(err)
		return nil, err
	}
	return csRes, nil
}

func (s *S8Proxy) DeleteSession(ctx context.Context, req *protos.DeleteSessionRequestPgw) (*protos.DeleteSessionResponsePgw, error) {
	// TODO make this stateless once MME has all the requiered information
	cPgwUDPAddr, err := s.configOrRequestedPgwAddress(req.PgwAddrs)
	if err != nil {
		err = fmt.Errorf("Delete Session failed due to missing server address: %s", err)
		glog.Error(err)
		return nil, err
	}

	dsReqMsg := buildDeleteSessionRequestMsg(req)
	cdRes, err := s.sendAndReceiveDeleteSession(req, cPgwUDPAddr, dsReqMsg)
	if err != nil {
		glog.Errorf("Couldnt delete session for IMSI %s:, %s", req.Imsi, err)
		return nil, err
	}
	// remove session from the s8_proxy client
	s.gtpClient.RemoveSessionByIMSI(req.Imsi)
	return cdRes, nil
}

func (s *S8Proxy) SendEcho(ctx context.Context, req *protos.EchoRequest) (*protos.EchoResponse, error) {
	cPgwUDPAddr, err := s.configOrRequestedPgwAddress(req.PgwAddrs)
	if err != nil {
		err = fmt.Errorf("SendEcho to %s failed: %s", cPgwUDPAddr, err)
		glog.Error(err)
		return nil, err
	}
	err = s.sendAndReceiveEchoRequest(cPgwUDPAddr)
	if err != nil {
		return nil, err
	}
	return &protos.EchoResponse{}, nil
}

// configOrRequestedPgwAddress returns an UDPAddrs if the passed string corresponds to a valid ip,
// otherwise it uses the server address configured on s8_proxy
func (s *S8Proxy) configOrRequestedPgwAddress(pgwAddrsFromRequest string) (*net.UDPAddr, error) {
	addrs := ParseAddress(pgwAddrsFromRequest)
	if addrs != nil {
		// address comming from string has precednece
		return addrs, nil
	}
	if s.config.ServerAddr != nil {
		return s.config.ServerAddr, nil
	}
	return nil, fmt.Errorf("Neither the request nor s8_proxy has a valid server (pgw) address")
}
