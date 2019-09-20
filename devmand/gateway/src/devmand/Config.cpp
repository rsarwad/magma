// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/Config.h>

namespace devmand {

DEFINE_string(listen_interface, "eth0", "The interface to listen on.");
DEFINE_string(
    device_configuration_file,
    "/etc/devmand/devices.yml",
    "Accepts .yml or .mconfig files. Inotify watches the file, and applies necessary changes.");
DEFINE_uint64(poll_interval, 10, "The polling interval in seconds.");

} // namespace devmand
