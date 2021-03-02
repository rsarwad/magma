/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#include "EventdClient.h"

#include "ServiceRegistrySingleton.h"

using grpc::ClientContext;
using grpc::Status;
using magma::orc8r::Event;
using magma::orc8r::EventService;
using magma::orc8r::Void;

namespace magma {

AsyncEventdClient& AsyncEventdClient::getInstance() {
  static AsyncEventdClient instance;
  return instance;
}

AsyncEventdClient::AsyncEventdClient() {
  std::shared_ptr<Channel> channel;
  channel = ServiceRegistrySingleton::Instance()->GetGrpcChannel(
      "eventd", ServiceRegistrySingleton::LOCAL);
  stub_ = EventService::NewStub(channel);
}

void AsyncEventdClient::log_event(
    const Event& request, std::function<void(Status status, Void)> callback) {
  auto local_response =
      new AsyncLocalResponse<Void>(std::move(callback), RESPONSE_TIMEOUT);
  local_response->set_response_reader(std::move(
      stub_->AsyncLogEvent(local_response->get_context(), request, &queue_)));
}

}  // namespace magma
