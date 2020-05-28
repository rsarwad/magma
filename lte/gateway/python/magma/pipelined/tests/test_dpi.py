"""
Copyright (c) 2019-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest
import warnings
from concurrent.futures import Future

from lte.protos.mconfig.mconfigs_pb2 import PipelineD
from lte.protos.policydb_pb2 import FlowMatch
from lte.protos.pipelined_pb2 import FlowRequest


from magma.pipelined.tests.app.start_pipelined import (
    TestSetup,
    PipelinedController,
)
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.pipelined_test_util import (
    start_ryu_app_thread,
    stop_ryu_app_thread,
    create_service_manager,
    SnapshotVerifier,
)


class DPITest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'

    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(DPITest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([PipelineD.DPI], [])

        dpi_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[PipelinedController.DPI,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.DPI:
                    dpi_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': cls.BRIDGE_IP,
                'ovs_gtp_port_number': 32768,
                'clean_restart': True,
                'setup_type': 'LTE',
                'dpi': {
                    'enabled': False,
                    'mon_port': 'mon1',
                    'mon_port_number': 32769,
                    'idle_timeout': 42,
                },
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)

        cls.thread = start_ryu_app_thread(test_setup)
        cls.dpi_controller = dpi_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def test_add_app_rules(self):
        """
        Test DPI classifier flows are properly added

        Assert:
            1 FLOW_CREATED -> no rule added as its not classified yet
            1 App not tracked -> no rule installed(`notanAPP`)
            3 App types are matched on:
                facebook other
                google_docs other
                viber audio
        """
        MAC_DEST = "5e:cc:cc:b1:49:4b"
        flow_match1 = FlowMatch(
            ip_proto=FlowMatch.IPPROTO_TCP, ipv4_dst='45.10.0.8',
            ipv4_src='1.2.3.4', tcp_dst=80, tcp_src=51115,
            direction=FlowMatch.UPLINK
        )
        flow_match2 = FlowMatch(
            ip_proto=FlowMatch.IPPROTO_TCP, ipv4_dst='1.10.0.1',
            ipv4_src='6.2.3.1', tcp_dst=111, tcp_src=222,
            direction=FlowMatch.UPLINK
        )
        flow_match3 = FlowMatch(
            ip_proto=FlowMatch.IPPROTO_UDP, ipv4_dst='22.2.2.24',
            ipv4_src='15.22.32.2', udp_src=111, udp_dst=222,
            direction=FlowMatch.UPLINK
        )
        flow_match_for_no_proto = FlowMatch(
            ip_proto=FlowMatch.IPPROTO_UDP, ipv4_dst='1.1.1.1'
        )
        flow_match_not_added = FlowMatch(
            ip_proto=FlowMatch.IPPROTO_UDP, ipv4_src='22.22.22.22'
        )
        self.dpi_controller.add_classify_flow(
            flow_match_not_added, FlowRequest.FLOW_CREATED,
            'nickproto', 'bestproto', MAC_DEST, MAC_DEST)
        self.dpi_controller.add_classify_flow(
            flow_match_for_no_proto, FlowRequest.FLOW_PARTIAL_CLASSIFICATION,
            'notanAPP', 'null', MAC_DEST, MAC_DEST)
        self.dpi_controller.add_classify_flow(
            flow_match1, FlowRequest.FLOW_PARTIAL_CLASSIFICATION,
            'base.ip.http.facebook', 'NotReal', MAC_DEST, MAC_DEST)
        self.dpi_controller.add_classify_flow(
            flow_match2, FlowRequest.FLOW_PARTIAL_CLASSIFICATION,
            'base.ip.https.google_gen.google_docs', 'MAGMA',
            MAC_DEST, MAC_DEST)
        self.dpi_controller.add_classify_flow(
            flow_match3, FlowRequest.FLOW_PARTIAL_CLASSIFICATION,
            'base.ip.udp.viber', 'AudioTransfer Receiving',
            MAC_DEST, MAC_DEST)

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)
        with snapshot_verifier:
            pass

    def test_remove_app_rules(self):
        """
        Test DPI classifier flows are properly removed

        Assert:
            Remove the facebook match flow
        """
        MAC_DEST = "5e:cc:cc:b1:49:4b"
        flow_match1 = FlowMatch(
            ip_proto=FlowMatch.IPPROTO_TCP, ipv4_dst='45.10.0.8',
            ipv4_src='1.2.3.4', tcp_dst=80, tcp_src=51115,
            direction=FlowMatch.UPLINK
        )
        self.dpi_controller.remove_classify_flow(flow_match1, MAC_DEST, MAC_DEST)

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)
        with snapshot_verifier:
            pass


if __name__ == "__main__":
    unittest.main()
