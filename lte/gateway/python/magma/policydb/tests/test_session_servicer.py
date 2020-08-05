"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import unittest
from typing import List
from lte.protos.mconfig import mconfigs_pb2
from lte.protos.policydb_pb2 import RatingGroup
from lte.protos.session_manager_pb2 import CreateSessionRequest, \
    UpdateSessionRequest, CreditUsageUpdate, SessionTerminateRequest, \
    CreditUpdateResponse, CreditLimitType
from lte.protos.subscriberdb_pb2 import SubscriberData, LTESubscription
from orc8r.protos.common_pb2 import NetworkID
from magma.policydb.servicers.session_servicer import SessionRpcServicer


CSR_STATIC_RULES = '[rule_id: "redirect"]'


USR = '''
  responses {
    success: true
    sid: "abc"
    charging_key: 1
    limit_type: INFINITE_UNMETERED
  }
  responses {
    success: true
    sid: "abc"
    charging_key: 2
    limit_type: INFINITE_METERED
  }'''


class MockSubscriberDBStub:
    def __init__(self):
        pass

    def ListSubscribers(self) -> List[str]:
        return ["IMSI001010000000001", "IMSI001010000000002"]

    def GetSubscriberData(self, _: NetworkID) -> SubscriberData:
        return SubscriberData(
            lte=LTESubscription(
                assigned_policies=["redirect"],
            )
        )


class SessionRpcServicerTest(unittest.TestCase):
    def setUp(self):
        rating_groups_by_id = {
            1: RatingGroup(
                id=1,
                limit_type=RatingGroup.INFINITE_UNMETERED,
            ),
            2: RatingGroup(
                id=2,
                limit_type=RatingGroup.INFINITE_METERED,
            ),
        }
        self.servicer = SessionRpcServicer(self._get_mconfig(),
                                           rating_groups_by_id,
                                           MockSubscriberDBStub())

    def _get_mconfig(self) -> mconfigs_pb2.PolicyDB:
        return mconfigs_pb2.PolicyDB(
            log_level=1,
        )

    def tearDown(self):
        # TODO: not sure if this is actually needed
        self.servicer = None

    def test_CreateSession(self):
        """
        Create a session

        Assert:
            There is a high volume of granted credits
            A static rule is installed for redirection to the captive portal
            server
            Two dynamic rules are installed for traffic from UE to captive
            portal and vice versa
        """
        msg = CreateSessionRequest()
        msg.session_id = '1234'
        msg.imsi_plmn_id = '00101'
        resp = self.servicer.CreateSession(msg, None)

        # There should be a static rule installed for the redirection
        static_rules = self._rm_whitespace(str(resp.static_rules))
        expected = self._rm_whitespace(CSR_STATIC_RULES)
        self.assertEqual(static_rules, expected, 'There should be one static '
                         'rule installed for redirection.')

        # Credit granted should be unlimited and un-metered
        credit_limit_type = resp.credits[0].limit_type
        expected = CreditLimitType.Value("INFINITE_UNMETERED")
        self.assertEqual(credit_limit_type, expected, 'There should be an '
                         'infinite, unmetered credit grant')

    def test_UpdateSession(self):
        """
        Update a session

        Assert:
            An arbitrarily large amount of credit should be granted.
        """
        msg = UpdateSessionRequest()
        credit_update = CreditUsageUpdate()
        credit_update.sid = 'abc'
        msg.updates.extend([credit_update])
        resp = self.servicer.UpdateSession(msg, None)
        session_response = self._rm_whitespace(str(resp))
        expected = self._rm_whitespace(USR)
        self.assertEqual(session_response, expected, 'There should be an '
                         'infinite, unmetered credit grant and an infinite, '
                         'metered credit grant')

    def test_TerminateSession(self):
        """
        Terminate a session

        Assert:
            The session can be terminated successfully. Will always succeed.
        """
        msg = SessionTerminateRequest()
        msg.sid = 'abc'
        msg.session_id = 'session_id_123'
        resp = self.servicer.TerminateSession(msg, None)
        self.assertEqual(resp.sid, 'abc', 'SID should be same as request')
        self.assertEqual(resp.session_id, 'session_id_123', 'session ID should '
                         'be same as request')

    def _rm_whitespace(self, inp: str) -> str:
        return inp.replace(' ', '').replace('\n', '')
