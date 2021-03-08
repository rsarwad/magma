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
	"sync"
	"testing"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/gtp"
	"magma/feg/gateway/services/s8_proxy/servicers/mock_pgw"

	"github.com/stretchr/testify/assert"
)

const (
	GtpTimeoutForTest = gtp.DefaultGtpTimeout // use the same the default value defined in s8_proxy
	//port 0 means golang will choose the port. Selected port will be injected on getDefaultConfig
	s8proxyAddrs = ":0" // equivalent to sgwAddrs
	pgwAddrs     = "127.0.0.1:0"
	IMSI1        = "123456789012345"
	BEARER       = 5
	AGWTeidU     = uint32(10)
	AGWTeidC     = uint32(2)
)

func TestS8proxyCreateAndDeleteSession(t *testing.T) {
	// set up client ans server
	s8p, mockPgw := startSgwAndPgw(t, GtpTimeoutForTest)
	defer mockPgw.Close()

	// ------------------------
	// ---- Create Session ----
	csReq := getDefaultCreateSessionRequest(mockPgw.LocalAddr().String())

	// Send and receive Create Session Request
	csRes, err := s8p.CreateSession(context.Background(), csReq)
	assert.NoError(t, err)
	assert.NotEmpty(t, csRes)

	// check User Plane FTEID was received properly
	assert.Equal(t, mockPgw.LastTEIDu, csRes.BearerContext.UserPlaneFteid.Teid)
	assert.NotEmpty(t, csRes.BearerContext.UserPlaneFteid.Ipv4Address)
	assert.Empty(t, csRes.BearerContext.UserPlaneFteid.Ipv6Address)

	// check Agw control Plane TEID
	session, err := s8p.gtpClient.GetSessionByIMSI(IMSI1)
	assert.NoError(t, err)
	agwTeidC, err := session.GetTEID(gtp.SGWControlPlaneIfType)
	assert.Equal(t, AGWTeidC, agwTeidC)

	// check Pgw Control Plane TEID
	assert.NotEqual(t, 0, csRes.CPgwFteid)

	// check received QOS
	sentQos := csReq.BearerContext.Qos
	receivedQos := mockPgw.LastQos
	assert.Equal(t, sentQos.Gbr.BrDl, receivedQos.Gbr.BrDl)
	assert.Equal(t, sentQos.Gbr.BrUl, receivedQos.Gbr.BrUl)
	assert.Equal(t, sentQos.Mbr.BrDl, receivedQos.Mbr.BrDl)
	assert.Equal(t, sentQos.Mbr.BrUl, receivedQos.Mbr.BrUl)
	assert.Equal(t, sentQos.Qci, receivedQos.Qci)

	// ------------------------
	// ---- Delete Session ----
	cdReq := getDeleteSessionRequest(mockPgw.LocalAddr().String(), csRes.CPgwFteid)

	_, err = s8p.DeleteSession(context.Background(), cdReq)
	assert.NoError(t, err)
	// session shouldnt exist anymore
	_, err = s8p.gtpClient.GetSessionByIMSI(IMSI1)
	assert.Error(t, err)

	// Create again the same session
	csRes, err = s8p.CreateSession(context.Background(), csReq)
	assert.NoError(t, err)
	assert.NotEmpty(t, csRes)
}

func TestS8ProxyDeleteInexistentSession(t *testing.T) {
	s8p, mockPgw := startSgwAndPgw(t, 200*time.Millisecond)
	defer mockPgw.Close()

	// ------------------------
	// ---- Delete Session inexistent session ----
	cdReq := &protos.DeleteSessionRequestPgw{Imsi: "000000000000015"}
	cdReq = &protos.DeleteSessionRequestPgw{
		PgwAddrs: mockPgw.LocalAddr().String(),
		Imsi:     "000000000000015",
		BearerId: 4,
		CAgwTeid: 88,
		CPgwFteid: &protos.Fteid{
			Ipv4Address: pgwAddrs,
			Teid:        87,
		},
	}
	_, err := s8p.DeleteSession(context.Background(), cdReq)
	assert.Error(t, err)
	assert.Equal(t, mockPgw.LastTEIDc, uint32(87))

}

func TestS8proxyCreateSessionDeniedService(t *testing.T) {
	// set up client ans server
	s8p, mockPgw := startSgwAndPgw(t, GtpTimeoutForTest)
	defer mockPgw.Close()

	// ------------------------
	// ---- Create Session ----
	csReq := getDefaultCreateSessionRequest(mockPgw.LocalAddr().String())

	// PGW denies service
	mockPgw.SetCreateSessionWithErrorCause()
	csRes, err := s8p.CreateSession(context.Background(), csReq)
	assert.Error(t, err)
	assert.Empty(t, csRes)
}

func TestS8proxyManyCreateAndDeleteSession(t *testing.T) {
	// set up client ans server
	s8p, mockPgw := startSgwAndPgw(t, GtpTimeoutForTest)
	defer mockPgw.Close()

	// ------------------------
	// ---- Create Sessions ----
	nRequest := 100
	pgwActualAddrs := mockPgw.LocalAddr().String()
	csReqs := getMultipleCreateSessionRequest(nRequest, pgwActualAddrs)

	// routines will write on specific index
	errors := make([]error, nRequest)
	csResps := make([]*protos.CreateSessionResponsePgw, nRequest)
	var wg sync.WaitGroup
	// PGW denies service
	for i, csReq := range csReqs {
		wg.Add(1)
		csReqShadow := csReq
		index := i
		go func() {
			defer wg.Done()
			csResps[index], errors[index] = s8p.CreateSession(context.Background(), csReqShadow)
		}()
	}
	wg.Wait()

	// Check all sessions were created
	assert.Equal(t, nRequest, len(errors))
	for _, err := range errors {
		assert.NoError(t, err, "Some sessions return error: %s", err)
	}
	for _, csReq := range csReqs {
		_, err := s8p.gtpClient.GetSessionByIMSI(csReq.Imsi)
		assert.NoError(t, err)
	}

	// ---- Delete Sessions ----
	errors = make([]error, nRequest)
	for i, csReq := range csReqs {
		wg.Add(1)
		csReqShadow := csReq
		csResShadow := csResps[i]
		index := i
		go func() {
			defer wg.Done()
			cdReq := &protos.DeleteSessionRequestPgw{
				PgwAddrs:  pgwActualAddrs,
				Imsi:      csReqShadow.Imsi,
				BearerId:  csResShadow.BearerContext.Id,
				CAgwTeid:  csResShadow.CAgwTeid,
				CPgwFteid: csResShadow.CPgwFteid,
			}

			_, err := s8p.DeleteSession(context.Background(), cdReq)
			errors[index] = err
		}()
	}
	wg.Wait()

	assert.Equal(t, nRequest, len(errors))
	for _, err := range errors {
		assert.NoError(t, err)
	}

	// check sessions are deleted
	for _, csReq := range csReqs {
		_, err := s8p.gtpClient.GetSessionByIMSI(csReq.Imsi)
		assert.Error(t, err)
	}
}

// TestS8proxyCreateSessionWrongSgwTEIDcFromPgw creates the situation where the PGW responds to the
// proper sequence message but with wrong SgwTeidC
func TestS8proxyCreateSessionWrongSgwTEIDcFromPgw(t *testing.T) {
	// set up client ans server
	// this test will timeout, reducing  gtp timeout to prevent waiting too much
	s8p, mockPgw := startSgwAndPgw(t, 200*time.Millisecond)
	defer mockPgw.Close()

	// ------------------------
	// ---- Create Session ----
	csReq := getDefaultCreateSessionRequest(mockPgw.LocalAddr().String())

	// PGW denies service
	mockPgw.CreateSessionOptions.SgwTeidc = 99
	csRes, err := s8p.CreateSession(context.Background(), csReq)
	assert.Error(t, err)
	assert.Empty(t, csRes)
}

func TestS8proxyEcho(t *testing.T) {
	s8p, mockPgw := startSgwAndPgw(t, GtpTimeoutForTest)
	defer mockPgw.Close()

	//------------------------
	//---- Echo Request ----
	eReq := &protos.EchoRequest{PgwAddrs: mockPgw.LocalAddr().String()}
	_, err := s8p.SendEcho(context.Background(), eReq)
	assert.NoError(t, err)
}

// startSgwAndPgw starts s8_proxy and a mock pgw for testing
func startSgwAndPgw(t *testing.T, gtpTimeout time.Duration) (*S8Proxy, *mock_pgw.MockPgw) {
	// Create and run PGW
	mockPgw, err := mock_pgw.NewStarted(nil, pgwAddrs)
	if err != nil {
		t.Fatalf("Error creating mock PGW: +%s", err)
	}

	// in case pgwAddres has a 0 port, mock_pgw will chose the port. With this variable we make
	// sure we use the right address (this only happens in testing)
	actualPgwAddress := mockPgw.LocalAddr().String()
	fmt.Printf("Running PGW at %s\n", actualPgwAddress)

	// Run S8_proxy
	config := getDefaultConfig(mockPgw.LocalAddr().String(), gtpTimeout)
	s8p, err := NewS8Proxy(config)
	if err != nil {
		t.Fatalf("Error creating S8 proxy +%s", err)
	}
	return s8p, mockPgw
}

func getDefaultCreateSessionRequest(pgwAddrs string) *protos.CreateSessionRequestPgw {
	_, offset := time.Now().Zone()
	return &protos.CreateSessionRequestPgw{
		PgwAddrs: pgwAddrs,
		Imsi:     IMSI1,
		Msisdn:   "300000000000003",
		Mei:      "111",
		CAgwTeid: AGWTeidC,
		ServingNetwork: &protos.ServingNetwork{
			Mcc: "222",
			Mnc: "333",
		},
		RatType: protos.RATType_EUTRAN,
		BearerContext: &protos.BearerContext{
			Id: BEARER,
			UserPlaneFteid: &protos.Fteid{
				Ipv4Address: "127.0.0.10",
				Ipv6Address: "",
				Teid:        AGWTeidU,
			},
			Qos: &protos.QosInformation{
				Pci:                     0,
				PriorityLevel:           0,
				PreemptionCapability:    0,
				PreemptionVulnerability: 0,
				Qci:                     9,
				Gbr: &protos.Ambr{
					BrUl: 123,
					BrDl: 234,
				},
				Mbr: &protos.Ambr{
					BrUl: 567,
					BrDl: 890,
				},
			},
		},
		PdnType: protos.PDNType_IPV4,
		Paa: &protos.PdnAddressAllocation{
			Ipv4Address: "10.0.0.10",
			Ipv6Address: "",
			Ipv6Prefix:  0,
		},

		Apn:           "internet.com",
		SelectionMode: protos.SelectionModeType_APN_provided_subscription_verified,
		Ambr: &protos.Ambr{
			BrUl: 999,
			BrDl: 888,
		},
		Uli: &protos.UserLocationInformation{
			Lac:    1,
			Ci:     2,
			Sac:    3,
			Rac:    4,
			Tac:    5,
			Eci:    6,
			MeNbi:  7,
			EMeNbi: 8,
		},
		IndicationFlag: nil,
		TimeZone: &protos.TimeZone{
			DeltaSeconds:       int32(offset),
			DaylightSavingTime: 0,
		},
	}
}

func getMultipleCreateSessionRequest(nRequest int, pgwAddrs string) []*protos.CreateSessionRequestPgw {
	res := []*protos.CreateSessionRequestPgw{}
	_, offset := time.Now().Zone()
	for i := 0; i < nRequest; i++ {
		newReq := &protos.CreateSessionRequestPgw{
			PgwAddrs: pgwAddrs,
			Imsi:     fmt.Sprintf("%d", 100000000000000+i),
			Msisdn:   fmt.Sprintf("%d", 17730000000+i),
			Mei:      fmt.Sprintf("%d", 900000000000000+i),
			CAgwTeid: AGWTeidC + uint32(i),
			ServingNetwork: &protos.ServingNetwork{
				Mcc: "222",
				Mnc: "333",
			},
			RatType: protos.RATType_EUTRAN,
			BearerContext: &protos.BearerContext{
				Id: BEARER,
				UserPlaneFteid: &protos.Fteid{
					Ipv4Address: "127.0.0.10",
					Ipv6Address: "",
					Teid:        AGWTeidU + uint32(i),
				},
				Qos: &protos.QosInformation{
					Pci:                     0,
					PriorityLevel:           0,
					PreemptionCapability:    0,
					PreemptionVulnerability: 0,
					Qci:                     9,
					Gbr: &protos.Ambr{
						BrUl: 123,
						BrDl: 234,
					},
					Mbr: &protos.Ambr{
						BrUl: 567,
						BrDl: 890,
					},
				},
			},
			PdnType: protos.PDNType_IPV4,
			Paa: &protos.PdnAddressAllocation{
				Ipv4Address: "10.0.0.10",
				Ipv6Address: "",
				Ipv6Prefix:  0,
			},

			Apn:           "internet.com",
			SelectionMode: protos.SelectionModeType_APN_provided_subscription_verified,
			Ambr: &protos.Ambr{
				BrUl: 999,
				BrDl: 888,
			},
			Uli: &protos.UserLocationInformation{
				Lac:    1,
				Ci:     2,
				Sac:    3,
				Rac:    4,
				Tac:    5,
				Eci:    6,
				MeNbi:  7,
				EMeNbi: 8,
			},
			IndicationFlag: nil,
			TimeZone: &protos.TimeZone{
				DeltaSeconds:       int32(offset),
				DaylightSavingTime: 0,
			},
		}
		res = append(res, newReq)
	}
	return res
}

func getDeleteSessionRequest(pgwAddrs string, cPgwFteid *protos.Fteid) *protos.DeleteSessionRequestPgw {
	res := &protos.DeleteSessionRequestPgw{
		PgwAddrs:  pgwAddrs,
		Imsi:      IMSI1,
		BearerId:  BEARER,
		CAgwTeid:  AGWTeidC,
		CPgwFteid: cPgwFteid,
	}
	return res
}

func getDefaultConfig(pgwActualAddrs string, gtpTimeout time.Duration) *S8ProxyConfig {
	return &S8ProxyConfig{
		GtpTimeout: gtpTimeout,
		ClientAddr: s8proxyAddrs,
	}
}
