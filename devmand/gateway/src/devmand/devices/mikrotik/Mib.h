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

#pragma once

#include <folly/futures/Future.h>

#include <devmand/channels/snmp/Channel.h>

namespace devmand {
namespace devices {
namespace mikrotik {

class Mib {
 public:
  Mib() = delete;
  ~Mib() = delete;
  Mib(const Mib&) = delete;
  Mib& operator=(const Mib&) = delete;
  Mib(Mib&&) = delete;
  Mib& operator=(Mib&&) = delete;

 public:
  static folly::Future<std::string> getBaseMac(
      channels::snmp::Channel& channel);

  static folly::Future<std::string> getSerialNumber(
      channels::snmp::Channel& channel);

  static folly::Future<std::string> getFirmwareVersion(
      channels::snmp::Channel& channel);

  static folly::Future<std::string> getModel(channels::snmp::Channel& channel);

  // TODO move into DISMAN-EVENT-MIB
  static folly::Future<std::string> getUpTime(channels::snmp::Channel& channel);

  static folly::Future<std::string> getLongtitude(
      channels::snmp::Channel& channel);

  static folly::Future<std::string> getLatitude(
      channels::snmp::Channel& channel);

  static folly::Future<std::string> getAltitude(
      channels::snmp::Channel& channel);

  static folly::Future<std::string> getIpv4Address(
      channels::snmp::Channel& channel);

  static folly::Future<std::string> getIpv6Address(
      channels::snmp::Channel& channel);
};

} // namespace mikrotik
} // namespace devices
} // namespace devmand
