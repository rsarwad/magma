"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import asyncio
from collections import defaultdict
from enum import Enum
from magma.pipelined.qos.qos_tc_impl import TCManager
from magma.pipelined.qos.qos_meter_impl import MeterManager
from magma.pipelined.qos.types import (QosInfo, get_json, get_key,
                                       get_subscriber_key)
from magma.pipelined.qos.utils import QosStore

from lte.protos.policydb_pb2 import FlowMatch
import logging

LOG = logging.getLogger('pipelined.qos.common')


class QosImplType(Enum):
    LINUX_TC = "linux_tc"
    OVS_METER = "ovs_meter"

    @staticmethod
    def list():
        return list(map(lambda t: t.value, QosImplType))


class QosManager(object):
    @staticmethod
    def getqos_impl(datapath, loop, config):
        try:
            qos_impl_type = QosImplType(config["qos"]["impl"])
        except ValueError:
            print("{} is not a valid qos implementation type({})".format(
                qos_impl_type, QosImplType.list()))
            raise

        if qos_impl_type == QosImplType.OVS_METER:
            return MeterManager(datapath, loop, config)
        else:
            return TCManager(datapath, loop, config)

    def redisAvailable(self):
        try:
            self._qos_store.client.ping()
        except ConnectionError:
            return False
        return True

    def __init__(self, datapath, loop, config):
        # pylint: disable=unnecessary-lambda
        self._enable_qos = config['qos']['enable']
        if not self._enable_qos:
            return
        self.qos_impl = QosManager.getqos_impl(datapath, loop, config)
        self._loop = loop
        self._subscriber_map = defaultdict(lambda: defaultdict())
        self._clean_restart = config['clean_restart']
        self._qos_store = QosStore(self.__class__.__name__)
        self._initialized = False
        self._redis_conn_retry_secs = 1

    def setup(self):
        if not self._enable_qos:
            return

        if self.redisAvailable():
            return self._setupInternal()
        else:
            LOG.info("failed to connect to redis..retrying in %d secs",
                     self._redis_conn_retry_secs)
            self._loop.call_later(self._redis_conn_retry_secs, self.setup)

    def _setupInternal(self, ):
        LOG.info("Qos Setup")
        if self._clean_restart:
            LOG.info("clean start, wiping out existing state")
            self.qos_impl.destroy()
            self._qos_store.clear()
            self.qos_impl.setup()
            self._initialized = True
        else:
            # read existing state from qos_impl
            LOG.info("recovering existing state")

            def callback(fut):
                LOG.debug("read_all_state complete => \n%s", fut.result())
                qos_state = fut.result()
                try:
                    # populate state from db
                    in_store_qid = set()
                    purge_store_set = set()
                    for k, v in self._qos_store.items():
                        if v not in qos_state:
                            purge_store_set.add(k)
                            continue
                        in_store_qid.add(v)
                        _, imsi, rule_num, d = get_key(k)
                        self._subscriber_map[imsi][rule_num] = (v, d)

                    # purge entries from qos_store
                    for k in purge_store_set:
                        LOG.debug("purging qos_store entry %s qos_handle", k)
                        del self._qos_store[k]

                    # purge unreferenced qos configs from system
                    for qos_handle, d in qos_state.items():
                        if qos_handle not in in_store_qid:
                            LOG.debug("removing qos_handle %d", qos_handle)
                            self.qos_impl.remove_qos(qos_handle, d,
                                                     recovery_mode=True)

                    self._initialized = True
                    LOG.info("init complete with state recovered successfully")
                except Exception as e:  # pylint: disable=broad-except
                    # in case of any exception start clean slate
                    LOG.error("error %s. restarting clean", str(e))
                    self._clean_restart = True
                    self.setup()
            asyncio.ensure_future(self.qos_impl.read_all_state(),
                                  loop=self._loop).add_done_callback(callback)

    def add_subscriber_qos(self,
                           imsi: str,
                           rule_num: int,
                           direction: FlowMatch.Direction,
                           qos_info: QosInfo):
        if not self._enable_qos or not self._initialized:
            LOG.error("add_subscriber_qos failed imsi %s rule_num %d \
                      direction %d failed qos not enabled or uninitialized",
                      imsi, rule_num, direction)
            return (None, None)

        LOG.debug("adding qos for imsi %s rule_num %d", imsi, rule_num)
        k = get_subscriber_key(imsi, rule_num, direction)
        qos_handle = self._qos_store.get(get_json(k))
        if qos_handle:
            LOG.debug("qos handle already exists for %s", k)
            return self.qos_impl.get_action_instruction(qos_handle)

        qos_handle = self.qos_impl.add_qos(direction, qos_info)
        self._subscriber_map[imsi][rule_num] = (qos_handle, direction)
        self._qos_store[get_json(k)] = qos_handle
        return self.qos_impl.get_action_instruction(qos_handle)

    def remove_subscriber_qos(self, imsi: str = "", rule_num: int = -1):
        if not self._enable_qos or not self._initialized:
            LOG.error("remove_subscriber_qos failed imsi %s rule_num %d \
                      failed qos not enabled or uninitialized", imsi,
                      rule_num)
            return

        LOG.debug("removing Qos for imsi %s rule_num %d", imsi, rule_num)
        if imsi:
            if imsi not in self._subscriber_map:
                LOG.debug("unable to find imsi %s", imsi)
                return

            if rule_num != -1:
                # delete queue associated with this rule
                if rule_num not in self._subscriber_map[imsi]:
                    LOG.error("unable to find rule_num %d for imsi %s",
                              rule_num, imsi)
                    return

                v = self._subscriber_map[imsi][rule_num]
                self.qos_impl.remove_qos(*v)
                if len(self._subscriber_map[imsi]) == 1:
                    del self._subscriber_map[imsi]
                else:
                    del self._subscriber_map[imsi][rule_num]

                # delete from qos store
                k = get_subscriber_key(imsi, rule_num, v[1])
                del self._qos_store[get_json(k)]
            else:
                # delete all queues associated with this subscriber
                for rule_num, v in self._subscriber_map[imsi].items():
                    self.qos_impl.remove_qos(*v)
                    del self._qos_store[get_json(get_subscriber_key(imsi,
                                        rule_num, v[1]))]

                self._subscriber_map.pop(imsi)
        else:
            # delete Qos queues associated with all subscribers
            LOG.info("removing Qos for all subscribers")
            for imsi, rule_map in self._subscriber_map.items():
                for rule_num, v in rule_map.items():
                    self.qos_impl.remove_qos(*v)
                    # delete from qos store
                    del self._qos_store[get_json(get_subscriber_key(imsi,
                                                 rule_num, v[1]))]

            self._subscriber_map.clear()
