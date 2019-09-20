// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/Application.h>

#include <chrono>
#include <future>
#include <iomanip>
#include <iostream>
#include <thread>

#include <folly/GLog.h>
#include <folly/executors/GlobalExecutor.h>
#include <folly/executors/IOExecutor.h>
#include <folly/json.h>

#include <devmand/Config.h>
#include <devmand/ErrorHandler.h>
#include <devmand/devices/Device.h>

using namespace std::chrono_literals;

namespace devmand {

Application::Application()
    : deviceFactory(*this),
      cartographer(
          [this](const cartography::DeviceConfig& deviceConfig) {
            add(deviceConfig);
          },
          [this](const cartography::DeviceConfig& deviceConfig) {
            del(deviceConfig);
          }) {
  ErrorHandler::executeWithCatch(
      [this]() -> void {
        snmpEngine = addEngine<channels::snmp::Engine>(name);
      },
      [this]() { this->statusCode = EXIT_FAILURE; });
}

std::string Application::getName() const {
  return name;
}

std::string Application::getVersion() const {
  return version;
}

// TODO move to device
void Application::pollDevices() {
  for (auto& device : devices) {
    device.second->updateSharedView(unifiedView);
  }
}

UnifiedView Application::getUnifiedView() {
  UnifiedView cpy;
  unifiedView.withULock([&cpy](auto& map) { cpy = map; });
  return cpy;
}

void Application::scheduleEvery(
    std::function<void()> event,
    const std::chrono::seconds& seconds) {
  eventBase.runInEventBaseThread([this, event, seconds]() {
    ErrorHandler::executeWithCatch([this, &event]() { event(); });

    std::function<void()> recurse = [this, event, seconds]() {
      this->scheduleEvery(event, seconds);
    };

    eventBase.scheduleAt(recurse, eventBase.now() + seconds);
  });
}

void Application::scheduleIn(
    std::function<void()> event,
    const std::chrono::seconds& seconds) {
  eventBase.runInEventBaseThread([this, event, seconds]() {
    std::function<void()> recurse = [this, event]() {
      ErrorHandler::executeWithCatch([this, event]() { event(); });
    };

    eventBase.scheduleAt(recurse, eventBase.now() + seconds);
  });
}

void Application::run() {
  if (statusCode != EXIT_SUCCESS) {
    LOG(ERROR) << "Not running application " << name << ", error on setup.";
    return;
  }

  LOG(INFO) << "Starting " << name << ".";

  ErrorHandler::executeWithCatch([this]() {
    // TODO so this part is not perm. Eventually we need to tie snmp into a
    // single epoll loop but until then just give it its own select thread
    auto snmpEngineThread = std::async(std::launch::async, [this] {
      assert(snmpEngine != nullptr);
      snmpEngine->run();
    });

    for (auto& service : services) {
      service->start();
    }

    // TODO move this to devices
    scheduleEvery(
        [this]() { pollDevices(); }, std::chrono::seconds(FLAGS_poll_interval));

    setGauge("devmand_running", 1);

    eventBase.loopForever();

    setGauge("devmand_running", 0);

    for (auto& service : services) {
      service->stop();
    }

    for (auto& service : services) {
      service->wait();
    }

    assert(snmpEngine != nullptr);
    snmpEngine->stopEventually();

    snmpEngineThread.wait();
  });

  LOG(INFO) << "Stopping " << name << ".";
}

int Application::status() const {
  return statusCode;
}

void Application::add(const cartography::DeviceConfig& deviceConfig) {
  ErrorHandler::executeWithCatch([this, &deviceConfig]() {
    add(deviceFactory.createDevice(deviceConfig));
    devices[deviceConfig.id]->applyConfig(deviceConfig.yangConfig);
  });
}

void Application::del(const cartography::DeviceConfig& deviceConfig) {
  if (devices.erase(deviceConfig.id) != 1) {
    LOG(ERROR) << "Failed to delete device " << deviceConfig.id;
  }
}

void Application::add(std::unique_ptr<devices::Device>&& device) {
  ErrorHandler::executeWithCatch(
      [this, &device]() {
        devices.emplace(
            device->getId(),
            std::forward<std::unique_ptr<devices::Device>>(device));
      },
      [this]() { this->statusCode = EXIT_FAILURE; });
}

void Application::setGauge(const std::string& key, double value) {
  setGauge(key, value, "", "");
}

void Application::setGauge(
    const std::string& key,
    double value,
    const std::string& label_name,
    const std::string& label_value) {
  for (auto& service : services) {
    service->setGauge(key, value, label_name, label_value);
  }
}

void Application::add(std::unique_ptr<Service>&& service) {
  ErrorHandler::executeWithCatch(
      [this, &service]() {
        services.emplace_back(std::forward<std::unique_ptr<Service>>(service));
      },
      [this]() { this->statusCode = EXIT_FAILURE; });
}

void Application::addPlatform(
    const std::string& platform,
    devices::Factory::PlatformBuilder platformBuilder) {
  deviceFactory.addPlatform(platform, platformBuilder);
}

void Application::setDefaultPlatform(
    devices::Factory::PlatformBuilder platformBuilder) {
  deviceFactory.setDefaultPlatform(platformBuilder);
}

void Application::addDeviceDiscoveryMethod(
    const std::shared_ptr<cartography::Method>& method) {
  assert(method != nullptr);
  cartographer.addDeviceDiscoveryMethod(method);
}

folly::EventBase& Application::getEventBase() {
  return eventBase;
}

DhcpdConfig& Application::getDhcpdConfig() {
  return dhcpdConfig;
}

} // namespace devmand
