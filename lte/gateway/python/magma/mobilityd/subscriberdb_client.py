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

import ipaddress
from typing import Optional
from lte.protos.subscriberdb_pb2 import APNConfiguration
from magma.subscriberdb.sid import SIDUtils

import grpc
import logging


class StaticIPInfo:
    """
    Operator can configure Static GW IP and MAC.
    This would be used by AGW services to generate networking
    configuration.
    """

    def __init__(self, ip: str, gw_ip: str, gw_mac: str,
                 vlan: str):
        self.ip = ipaddress.ip_address(ip)
        self.gw_mac = gw_mac
        gw_ip_parsed = None
        try:
            gw_ip_parsed = ipaddress.ip_address(gw_ip)
        except ValueError:
            logging.debug("invalid internet gw ip: %s", gw_ip)
        self.gw_ip = gw_ip_parsed
        self.vlan = vlan

    def __str__(self):
        return "IP: {} GW-IP: {} GW-MAC: {} VLAN: {}".format(self.ip,
                                                             self.gw_ip,
                                                             self.gw_mac,
                                                             self.vlan)


class SubscriberDbClient:
    def __init__(self, subscriberdb_rpc_stub):
        self.subscriber_client = subscriberdb_rpc_stub

    def get_subscriber_ip(self, sid: str) -> Optional[StaticIPInfo]:
        """
        Make RPC call to 'GetSubscriberData' method of local SubscriberDB
        service to get assigned IP address if any.
        """
        if self.subscriber_client is None:
            return None

        try:
            apn_config = self._find_ip_and_apn_config(sid)
            logging.debug("ip: Got APN: %s", apn_config)
            if apn_config:
                return StaticIPInfo(ip=apn_config.assigned_static_ip,
                                    gw_ip=apn_config.resource.gateway_ip,
                                    gw_mac=apn_config.resource.gateway_mac,
                                    vlan=apn_config.resource.vlan_id)

        except ValueError:
            logging.warning("Invalid data for sid %s: ", sid)

        except grpc.RpcError as err:
            logging.error(
                "GetSubscriberData while reading static ip, error[%s] %s",
                err.code(),
                err.details())
        return None

    def get_subscriber_apn_vlan(self, sid: str) -> int:
        """
        Make RPC call to 'GetSubscriberData' method of local SubscriberDB
        service to get assigned IP address if any.
        TODO: Move this API to separate APN configuration service.
        """
        if self.subscriber_client is None:
            return 0

        try:
            apn_config = self._find_ip_and_apn_config(sid)
            logging.debug("vlan: Got APN: %s", apn_config)
            if apn_config:
                return apn_config.resource.vlan_id

        except ValueError:
            logging.warning("Invalid data for sid %s: ", sid)
            return 0

        except grpc.RpcError as err:
            logging.error(
                "GetSubscriberData while reading vlan-id error[%s] %s",
                err.code(),
                err.details())
        return 0

    # use same API to retrieve IP address and related config.
    def _find_ip_and_apn_config(self, sid: str) -> (Optional[APNConfiguration]):
        if '.' in sid:
            imsi, apn_name_part = sid.split('.', maxsplit=1)
            apn_name, _ = apn_name_part.split(',', maxsplit=1)
        else:
            imsi, _ = sid.split(',', maxsplit=1)
            apn_name = ''

        logging.debug("Find APN config for: %s", sid)
        data = self.subscriber_client.GetSubscriberData(SIDUtils.to_pb(imsi))
        if data and data.non_3gpp and data.non_3gpp.apn_config:
            selected_apn_conf = None
            for apn_config in data.non_3gpp.apn_config:
                logging.debug("APN config: %s", apn_config)
                try:
                    if apn_config.assigned_static_ip:
                        ipaddress.ip_address(apn_config.assigned_static_ip)
                except ValueError:
                    continue
                if apn_config.service_selection == '*':
                    selected_apn_conf = apn_config
                elif apn_config.service_selection == apn_name:
                    selected_apn_conf = apn_config
                    break

            return selected_apn_conf

        return None
