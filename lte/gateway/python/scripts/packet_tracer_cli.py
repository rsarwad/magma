"""
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
"""

import fire
from lte.protos.pipelined_pb2 import SerializedRyuPacket
from lte.protos.pipelined_pb2_grpc import PipelinedStub
from magma.common.service_registry import ServiceRegistry
from ryu.lib.packet import ethernet, arp, ipv4, icmp, tcp
from ryu.lib.packet.ether_types import ETH_TYPE_ARP, ETH_TYPE_IP
from ryu.lib.packet.packet import Packet
from termcolor import colored


class PacketTracerCLI:
    """
    Packet tracer for magma OVS tables.
    Use to generate traffic from packets and send it through magma OVS tables.
    PacketTracer reports which OVS table caused a drop of the packet.
    """

    def raw(self, data, imsi='001010000000013'):
        """
        Send a packet constructed from raw bytes through the magma switch and
        display which tabled caused a drop
        (-1 if the packet wasn't dropped by any table)
        """
        data = bytes(data)
        pkt = Packet(data)

        # Send the packet to a grpc service
        chan = ServiceRegistry.get_rpc_channel('pipelined',
                                               ServiceRegistry.LOCAL)
        client = PipelinedStub(chan)

        print('Sending: {}'.format(pkt))
        table_id = client.TracePacket(SerializedRyuPacket(
            pkt=data,
            imsi=imsi,
        )).table_id

        if table_id == -1:
            print('Successfully passed through all the tables!')
        else:
            print('Dropped by table: {}'.format(table_id))

    def icmp(self, imsi='001010000000013',
             src_mac='00:00:00:00:00:00', src_ip='192.168.70.2',
             dst_mac='ff:ff:ff:ff:ff:ff', dst_ip='192.168.70.3'):
        """
        Send an ICMP packet through the magma switch and display which tabled
        caused a drop (-1 if the packet wasn't dropped by any table)
        """
        pkt = ethernet.ethernet(src=src_mac, dst=dst_mac) / \
              ipv4.ipv4(src=src_ip, dst=dst_ip, proto=1) / \
              icmp.icmp()
        pkt.serialize()
        self.raw(data=pkt.data, imsi=imsi)

    def arp(self, imsi='001010000000013',
            src_mac='00:00:00:00:00:00', src_ip='192.168.70.2',
            dst_mac='ff:ff:ff:ff:ff:ff', dst_ip='192.168.70.3'):
        """
        Send an ARP packet through the magma switch and display which tabled
        caused a drop (-1 if the packet wasn't dropped by any table)
        """
        pkt = ethernet.ethernet(ethertype=ETH_TYPE_ARP,
                                src=src_mac, dst=dst_mac) / \
              arp.arp(hwtype=arp.ARP_HW_TYPE_ETHERNET, proto=ETH_TYPE_IP,
                      hlen=6, plen=4,
                      opcode=arp.ARP_REQUEST,
                      src_mac=src_mac, src_ip=src_ip,
                      dst_mac=dst_mac, dst_ip=dst_ip)
        pkt.serialize()
        self.raw(data=pkt.data, imsi=imsi)

    def tcp(self, imsi='001010000000013',
            src_mac='00:00:00:00:00:00', src_ip='192.168.70.2', src_port=80,
            dst_mac='ff:ff:ff:ff:ff:ff', dst_ip='192.168.70.3', dst_port=80,
            bits=tcp.TCP_SYN, seq=0, ack=0):
        """
        Send a TCP packet through the magma switch and display which tabled
        caused a drop (-1 if the packet wasn't dropped by any table)
        """
        pkt = ethernet.ethernet(src=src_mac, dst=dst_mac) / \
              ipv4.ipv4(ttl=55, proto=6, src=src_ip, dst=dst_ip) / \
              tcp.tcp(src_port=src_port, dst_port=dst_port, bits=bits,
                      seq=seq, ack=ack)
        pkt.serialize()
        self.raw(data=pkt.data, imsi=imsi)

    def http(self, imsi='001010000000013',
             src_ip='192.168.70.2',
             dst_ip='8.8.8.8'):
        """
        Perform (mock) an HTTP handshake and send each of the 3 packets through
        the magma switch and display which tabled caused a drop
        (-1 if the packet wasn't dropped by any table)
        """
        self.tcp(imsi=imsi, src_ip=src_ip, dst_ip=dst_ip, bits=tcp.TCP_SYN)
        self.tcp(imsi=imsi, src_ip=dst_ip, dst_ip=src_ip, dst_port=20,
                 bits=(tcp.TCP_SYN | tcp.TCP_ACK),
                 seq=3833491143, ack=1)
        self.tcp(imsi=imsi, src_ip=src_ip, src_port=20, dst_ip=dst_ip,
                 bits=tcp.TCP_ACK,
                 seq=1, ack=3833491144)


if __name__ == '__main__':
    cli = PacketTracerCLI()
    try:
        fire.Fire(cli)
    except Exception as e:
        print(colored('Error', 'red'), e)
