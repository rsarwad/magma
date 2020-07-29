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

#include <fcntl.h>
#include <sys/inotify.h>
#include <sys/stat.h>
#include <iostream>
#include <list>
#include <map>
#include <stdexcept>
#include <string>

#include <folly/io/async/EventBase.h>
#include <folly/io/async/EventHandler.h>

#include <devmand/utils/FileUtils.h>

namespace devmand {

enum class FileEvent {
  Access = IN_ACCESS,
  Attrib = IN_ATTRIB,
  CloseNoWrite = IN_CLOSE_NOWRITE,
  CloseWrite = IN_CLOSE_WRITE,
  Create = IN_CREATE,
  Delete = IN_DELETE,
  DeleteSelf = IN_DELETE_SELF,
  Ignored = IN_IGNORED,
  IsDir = IN_ISDIR,
  Modify = IN_MODIFY,
  MoveSelf = IN_MOVE_SELF,
  MoveFrom = IN_MOVED_FROM,
  MoveTo = IN_MOVED_TO,
  Open = IN_OPEN,
  QOverflow = IN_Q_OVERFLOW,
  UnMount = IN_UNMOUNT
};

struct FileWatchEvent {
  FileEvent event;
  std::string filename;
};

using ActionCallback = std::function<void(FileWatchEvent event)>;

struct FileWatch {
  int fd;
  ActionCallback action;
  std::string filename;
};

class FileWatcher final : public folly::EventHandler {
 public:
  FileWatcher(folly::EventBase& _eventBase);
  FileWatcher() = delete;
  ~FileWatcher();
  FileWatcher(const FileWatcher&) = delete;
  FileWatcher& operator=(const FileWatcher&) = delete;
  FileWatcher(FileWatcher&&) = delete;
  FileWatcher& operator=(FileWatcher&&) = delete;

 public:
  bool addWatch(
      const std::string& filename,
      ActionCallback action = [](FileWatchEvent) {},
      bool shouldReadInitial = false,
      uint32_t eventMask = IN_ALL_EVENTS); // TODO make mask formable from enum

  virtual void handlerReady(uint16_t events) noexcept override;

 private:
  void removeWatch(int watchFd);

 private:
  folly::EventBase& eventBase;
  int inotifyFd{-1};
  std::map<int, FileWatch> watches;
};

} // namespace devmand
