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

#include "SessionState.h"
#include "SessionStore.h"
#include "StoredState.h"
#include "magma_logging.h"

namespace magma {
namespace lte {

SessionStore::SessionStore(std::shared_ptr<StaticRuleStore> rule_store)
    : rule_store_(rule_store),
      store_client_(std::make_shared<MemoryStoreClient>(rule_store)),
      metering_reporter_(std::make_shared<MeteringReporter>()) {}

SessionStore::SessionStore(
    std::shared_ptr<StaticRuleStore> rule_store,
    std::shared_ptr<RedisStoreClient> store_client)
    : rule_store_(rule_store),
      store_client_(store_client),
      metering_reporter_(std::make_shared<MeteringReporter>()) {}

SessionMap SessionStore::read_sessions(const SessionRead& req) {
  return store_client_->read_sessions(req);
}

SessionMap SessionStore::read_all_sessions() {
  return store_client_->read_all_sessions();
}

void SessionStore::set_and_save_reporting_flag(
    bool value, const UpdateSessionRequest& update_session_request,
    SessionUpdate& session_uc) {
  MLOG(MDEBUG) << "saving flag is_reporting = " << value << " on session store";
  auto session_map = store_client_->read_all_sessions();

  for (const CreditUsageUpdate& credit_update :
       update_session_request.updates()) {
    const std::string imsi       = credit_update.common_context().sid().id();
    const std::string session_id = credit_update.session_id();
    const CreditKey& ckey        = credit_update.usage().charging_key();
    const std::string mkey       = credit_update.usage().monitoring_key();

    SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);
    auto session_it = find_session(session_map, criteria);
    if (!session_it) {
      MLOG(MERROR) << session_id
                   << " not found when setting set_and_save_reporting_flag";
      continue;
    }

    auto& session   = **session_it;
    auto& credit_uc = session_uc[imsi][session_id];

    if (!session->set_credit_reporting(ckey, value, &credit_uc)) {
      MLOG(MDEBUG)
          << session_id
          << " set_and_save_reporting_flag couldn't set reporting for ckey "
          << ckey;
    }
  }

  for (const UsageMonitoringUpdateRequest& monitor_update :
       update_session_request.usage_monitors()) {
    const std::string imsi       = monitor_update.sid();
    const std::string session_id = monitor_update.session_id();
    const auto mkey              = monitor_update.update().monitoring_key();

    SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);
    auto session_it = find_session(session_map, criteria);
    if (!session_it) {
      MLOG(MERROR) << session_id
                   << " not found when setting set_and_save_reporting_flag";
      continue;
    }
    auto& session   = **session_it;
    auto& credit_uc = session_uc[imsi][session_id];

    if (!session->set_monitor_reporting(mkey, value, &credit_uc)) {
      MLOG(MDEBUG)
          << session_id
          << " set_and_save_reporting_flag couldn't set monitors for mkey:"
          << mkey;
    }
  }

  store_client_->write_sessions(std::move(session_map));
}

void SessionStore::sync_request_numbers(const SessionUpdate& update_criteria) {
  // Read the current stored state
  auto subscriber_ids = std::set<std::string>{};
  for (const auto& it : update_criteria) {
    subscriber_ids.insert(it.first);
  }
  auto session_map = store_client_->read_sessions(subscriber_ids);

  // Sync stored state so that subsequent reads have the right request_number
  MLOG(MDEBUG) << "Syncing request numbers into existing sessions";
  for (auto& it : session_map) {
    auto imsi = it.first;
    auto it2  = it.second.begin();
    while (it2 != it.second.end()) {
      auto updates    = update_criteria.find(it.first)->second;
      auto session_id = (*it2)->get_session_id();
      if (updates.find(session_id) != updates.end()) {
        (*it2)->increment_request_number(
            updates[session_id].request_number_increment);
      }
      ++it2;
    }
  }
  MLOG(MDEBUG) << "sync_request_numbers: Writing into session store";
  store_client_->write_sessions(std::move(session_map));
}

SessionMap SessionStore::read_sessions_for_deletion(const SessionRead& req) {
  auto session_map   = store_client_->read_sessions(req);
  auto session_map_2 = store_client_->read_sessions(req);
  // For all sessions of the subscriber, increment the request numbers
  for (const std::string& imsi : req) {
    for (auto& session : session_map_2[imsi]) {
      session->increment_request_number(1);
    }
  }
  store_client_->write_sessions(std::move(session_map_2));
  return session_map;
}

bool SessionStore::create_sessions(
    const std::string& subscriber_id, SessionVector sessions) {
  auto session_map           = SessionMap{};
  session_map[subscriber_id] = std::move(sessions);
  store_client_->write_sessions(std::move(session_map));
  return true;
}

bool SessionStore::update_sessions(const SessionUpdate& update_criteria) {
  // Read the current state
  auto subscriber_ids = std::set<std::string>{};
  for (const auto& it : update_criteria) {
    subscriber_ids.insert(it.first);
  }
  auto session_map = store_client_->read_sessions(subscriber_ids);
  // Now attempt to modify the state
  for (auto& it : session_map) {
    auto imsi = it.first;
    auto it2  = it.second.begin();
    while (it2 != it.second.end()) {
      auto updates    = update_criteria.find(it.first)->second;
      auto session_id = (*it2)->get_session_id();
      if (updates.find(session_id) != updates.end()) {
        auto update = updates[session_id];
        if (!(*it2)->apply_update_criteria(update)) {
          return false;
        }
        metering_reporter_->report_usage(imsi, session_id, update);

        if (update.is_session_ended) {
          // TODO: Instead of deleting from session_map, mark as ended and
          //       no longer mark on read
          it2 = it.second.erase(it2);
          continue;
        }
      }
      ++it2;
    }
  }
  return store_client_->write_sessions(std::move(session_map));
}

optional<SessionVector::iterator> SessionStore::find_session(
    SessionMap& session_map, SessionSearchCriteria criteria) {
  auto sm_it = session_map.find(criteria.imsi);
  if (sm_it == session_map.end()) {
    return {};
  }
  auto& sessions = sm_it->second;
  for (auto it = sessions.begin(); it != sessions.end(); ++it) {
    switch (criteria.search_type) {
      case IMSI_AND_SESSION_ID:
        if ((*it)->get_session_id() == criteria.secondary_key) {
          return it;
        }
        break;
      case IMSI_AND_TEID:
        if ((*it)->get_local_teid() == criteria.secondary_key_unit32) {
          return it;
        }
        break;
      case IMSI_AND_APN:
        if ((*it)->get_config().common_context.apn() ==
            criteria.secondary_key) {
          return it;
        }
        break;
      case IMSI_AND_UE_IPV4:
        if ((*it)->get_config().common_context.ue_ipv4() ==
            criteria.secondary_key) {
          return it;
        }
        break;
      case IMSI_AND_UE_IPV4_OR_IPV6:
        // cwag case (cwag doesn't store ip)
        if ((*it)->get_config().common_context.rat_type() ==
            RATType::TGPP_WLAN) {
          return it;
        }
        // other case(lte,5g)
        // lte case
        if ((*it)->get_config().common_context.ue_ipv4() ==
                criteria.secondary_key ||
            (*it)->get_config().common_context.ue_ipv6() ==
                criteria.secondary_key) {
          return it;
        }
        break;
      case IMSI_AND_BEARER:

        if ((*it)->get_config().common_context.rat_type() ==
            RATType::TGPP_LTE) {
          // lte case
          if ((*it)
                  ->get_config()
                  .rat_specific_context.lte_context()
                  .bearer_id() == criteria.secondary_key_unit32) {
            return it;
          }
        } else {
          // 5g and cwag
          MLOG(MERROR) << "find_session by bearer is not implemented "
                          "for cwg or 5g. Couldnt find session for "
                       << criteria.imsi;
          return it;
        }
        break;
    }
    continue;
  }
  return {};
}

SessionUpdate SessionStore::get_default_session_update(
    SessionMap& session_map) {
  SessionUpdate update = {};
  for (const auto& session_pair : session_map) {
    for (const auto& session : session_pair.second) {
      update[session_pair.first][session->get_session_id()] =
          get_default_update_criteria();
    }
  }
  return update;
}

}  // namespace lte
}  // namespace magma
