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

import logging
from lte.protos.mconfig import mconfigs_pb2
from lte.protos.policydb_pb2_grpc import PolicyAssignmentControllerStub
from lte.protos.session_manager_pb2_grpc import SessionProxyResponderStub
from lte.protos.subscriberdb_pb2_grpc import SubscriberDBStub
from magma.common.service import MagmaService
from magma.common.service_registry import ServiceRegistry
from magma.common.streamer import StreamerClient
from magma.policydb.basename_store import BaseNameDict
from magma.policydb.rating_group_store import RatingGroupsDict
from magma.policydb.reauth_handler import ReAuthHandler
from magma.policydb.rule_map_store import RuleAssignmentsDict
from magma.policydb.servicers.policy_servicer import PolicyRpcServicer
from magma.policydb.servicers.session_servicer import SessionRpcServicer
from .streamer_callback import PolicyDBStreamerCallback, \
    RuleMappingsStreamerCallback, RatingGroupsStreamerCallback


def main():
    """ main() for subscriberdb """
    service = MagmaService('policydb', mconfigs_pb2.PolicyDB())

    assignments_dict = RuleAssignmentsDict()
    basenames_dict = BaseNameDict()
    rating_groups_dict = RatingGroupsDict()
    sessiond_chan = ServiceRegistry.get_rpc_channel('sessiond',
                                                    ServiceRegistry.LOCAL)
    sessiond_stub = SessionProxyResponderStub(sessiond_chan)
    reauth_handler = ReAuthHandler(assignments_dict, sessiond_stub)

    # Add all servicers to the server
    chan = ServiceRegistry.get_rpc_channel('subscriberdb',
                                           ServiceRegistry.LOCAL)
    subscriberdb_stub = SubscriberDBStub(chan)
    session_servicer = SessionRpcServicer(service.mconfig,
                                          rating_groups_dict,
                                          subscriberdb_stub)
    session_servicer.add_to_server(service.rpc_server)

    orc8r_chan = ServiceRegistry.get_rpc_channel('policydb',
                                                 ServiceRegistry.CLOUD)
    policy_stub = PolicyAssignmentControllerStub(orc8r_chan)
    policy_servicer = PolicyRpcServicer(reauth_handler, basenames_dict,
                                        policy_stub)
    policy_servicer.add_to_server(service.rpc_server)

    # Start a background thread to stream updates from the cloud
    if service.config['enable_streaming']:
        stream = StreamerClient(
            {
                'policydb': PolicyDBStreamerCallback(),
                'rule_mappings': RuleMappingsStreamerCallback(reauth_handler,
                                                              basenames_dict,
                                                              assignments_dict),
                'rating_groups': RatingGroupsStreamerCallback(
                    rating_groups_dict),

            },
            service.loop,
        )
        stream.start()
    else:
        logging.info('enable_streaming set to False. Streamer disabled!')

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


if __name__ == "__main__":
    main()
