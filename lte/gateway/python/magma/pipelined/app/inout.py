"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
from collections import namedtuple
from ryu.ofproto.ofproto_v1_4 import OFPP_LOCAL

from .base import MagmaController
from magma.pipelined.openflow import flows
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import load_direction, Direction

# ingress and egress service names -- used by other controllers
INGRESS = "ingress"
EGRESS = "egress"


class InOutController(MagmaController):
    """
    A controller that sets up an openflow pipeline for Magma.

    The EPC controls table 0 which is the first table every packet touches.
    This controller owns the ingress and output portions of the pipeline, the
    first table a packet hits after the EPC controller's table 0 and the last
    table a packet hits before exiting the pipeline.
    """

    APP_NAME = "inout"

    InOutConfig = namedtuple(
        'InOutConfig',
        ['gtp_port', 'uplink_port_name'],
    )

    def __init__(self, *args, **kwargs):
        super(InOutController, self).__init__(*args, **kwargs)
        self.config = self._get_config(kwargs['config'])
        self._uplink_port = OFPP_LOCAL
        if (self.config.uplink_port_name):
            self._uplink_port = BridgeTools.get_ofport(self.config.uplink_port_name)

    def _get_config(self, config_dict):
        port_name = None
        if 'ovs_uplink_port_name' in config_dict:
            port_name = config_dict['ovs_uplink_port_name']

        return self.InOutConfig(
            gtp_port=config_dict['ovs_gtp_port_number'],
            uplink_port_name=port_name
        )

    def initialize_on_connect(self, datapath):
        self._clear_ingress_egress_tables(datapath)
        self._install_default_egress_flows(datapath)
        self._install_default_ingress_flows(datapath)

    def cleanup_on_disconnect(self, datapath):
        self._clear_ingress_egress_tables(datapath)

    def _clear_ingress_egress_tables(self, datapath):
        flows.delete_all_flows_from_table(datapath,
                                          self._service_manager.get_table_num(
                                              INGRESS))
        flows.delete_all_flows_from_table(datapath,
                                          self._service_manager.get_table_num(
                                              EGRESS))

    def _install_default_egress_flows(self, dp):
        """
        Egress table is the last table that a packet touches in the pipeline.
        Output downlink traffic to gtp port, uplink trafic to LOCAL

        Raises:
            MagmaOFError if any of the default flows fail to install.
        """
        downlink_match = MagmaMatch(direction=Direction.IN)
        flows.add_output_flow(dp, self._service_manager.get_table_num(EGRESS),
                              downlink_match, [],
                              output_port=self.config.gtp_port)

        uplink_match = MagmaMatch(direction=Direction.OUT)
        flows.add_output_flow(dp, self._service_manager.get_table_num(EGRESS),
                              uplink_match, [],
                              output_port=self._uplink_port)

    def _install_default_ingress_flows(self, dp):
        """
        Sets up the ingress table, the first step in the packet processing
        pipeline.

        This sets up flow rules to annotate packets with a metadata bit
        indicating the direction. Incoming packets are defined as packets
        originating from the LOCAL port, outgoing packets are defined as
        packets originating from the gtp port.

        All other packets bypass the pipeline.

        Note that the ingress rules do *not* install any flows that cause
        PacketIns (i.e., sends packets to the controller).

        Raises:
            MagmaOFError if any of the default flows fail to install.
        """
        parser = dp.ofproto_parser
        tbl_num = self._service_manager.get_table_num(INGRESS)
        next_table = self._service_manager.get_next_table_num(INGRESS)
        egress_table = self._service_manager.get_table_num(EGRESS)

        # set traffic direction bits
        # set a direction bit for outgoing (pn -> inet) traffic.
        match = MagmaMatch(in_port=self.config.gtp_port)
        actions = [load_direction(parser, Direction.OUT)]
        flows.add_resubmit_next_service_flow(dp, tbl_num, match,
                                             actions=actions,
                                             priority=flows.DEFAULT_PRIORITY,
                                             resubmit_table=next_table)

        # Allow passthrough pkts(skip pipeline and send to egress table)
        match = MagmaMatch(in_port=self.config.gtp_port,
                           direction=Direction.PASSTHROUGH)
        flows.add_resubmit_next_service_flow(dp, tbl_num, match,
                                             actions=actions,
                                             priority=flows.PASSTHROUGH_PRIORITY,
                                             resubmit_table=egress_table)

        # set a direction bit for incoming (internet -> UE) traffic.
        match = MagmaMatch(in_port=self._uplink_port)
        actions = [load_direction(parser, Direction.IN)]
        flows.add_resubmit_next_service_flow(dp, tbl_num, match,
                                             actions=actions,
                                             priority=flows.DEFAULT_PRIORITY,
                                             resubmit_table=next_table)

        # Allow passthrough pkts(skip pipeline and send to egress table)
        match = MagmaMatch(in_port=self._uplink_port,
                           direction=Direction.PASSTHROUGH)
        flows.add_resubmit_next_service_flow(dp, tbl_num, match,
                                             actions=actions,
                                             priority=flows.PASSTHROUGH_PRIORITY,
                                             resubmit_table=egress_table)
