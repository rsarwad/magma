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

package mock_pgw

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/wmnsk/go-gtp/gtpv1"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/gtp"

	"github.com/wmnsk/go-gtp/gtpv2"
	"github.com/wmnsk/go-gtp/gtpv2/ie"
	"github.com/wmnsk/go-gtp/gtpv2/message"
)

const (
	dummyUserPlanePgwIP = "10.0.0.1"
	gtpTimeout          = 500 * time.Millisecond
)

// MockPgw is just a wrapper around gtp.Client
type MockPgw struct {
	*gtp.Client
	LastValues
	CreateSessionOptions CreateSessionOptions
}

type LastValues struct {
	LastTEIDu uint32
	LastTEIDc uint32
	LastQos   *protos.QosInformation
}

// CreateSessionOptions to control Create Session Response values to produce errors
type CreateSessionOptions struct {
	SgwTeidc  uint32
	PgwFTEIDc uint32
	PgwFTEIDu uint32
}

func NewStarted(ctx context.Context, pgwAddrsStr string) (*MockPgw, error) {
	mPgw := New()
	err := mPgw.Start(ctx, pgwAddrsStr)
	if err != nil {
		return nil, err
	}
	return mPgw, nil
}

func New() *MockPgw {
	return &MockPgw{CreateSessionOptions: CreateSessionOptions{}}
}

func (mPgw *MockPgw) Start(ctx context.Context, pgwAddrsStr string) error {
	var err error
	mPgw.Client, err = gtp.NewRunningClient(ctx, pgwAddrsStr, gtpv2.IFTypeS5S8PGWGTPC, gtpTimeout)
	if err != nil {
		return fmt.Errorf("Failed to get SGW IP: %s", err)
	}

	// register handlers for ALL the message you expect remote endpoint to send.
	mPgw.AddHandlers(map[uint8]gtpv2.HandlerFunc{
		message.MsgTypeCreateSessionRequest:       mPgw.getHandleCreateSessionRequest(),
		message.MsgTypeModifyAccessBearersRequest: mPgw.getHandleModifyBearerRequest(),
		message.MsgTypeDeleteSessionRequest:       mPgw.getHandleDeleteSessionRequest(),
		//message.MsgTypeEchoRequest: mPgw.getHandleEchoRequest(), // ONLY FOR DEBUGGING PURPOSES
	})
	return nil
}

func (mPgw *MockPgw) SetCreateSessionWithErrorCause() {
	mPgw.AddHandlers(map[uint8]gtpv2.HandlerFunc{
		message.MsgTypeCreateSessionRequest: mPgw.getHandleCreateSessionRequestWithDeniedService(),
	})

}

// ONLY FOR DEBUGGING PURPOSES
// getHandleEchoResponse is the same method as the one found in Go-GTP gtpv1.handleEchoResponse
func (mPgw *MockPgw) getHandleEchoRequest() gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, sgwAddr net.Addr, msg message.Message) error {
		if _, ok := msg.(*message.EchoRequest); !ok {
			return gtpv1.ErrUnexpectedType
		}
		// respond with EchoResponse.
		return c.RespondTo(sgwAddr, msg, message.NewEchoResponse(0, ie.NewRecovery(c.RestartCounter)))
	}
}
