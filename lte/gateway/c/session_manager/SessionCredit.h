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
#pragma once

#include <lte/protos/session_manager.grpc.pb.h>

#include "StoredState.h"

namespace magma {
/**
 * SessionCredit tracks all the credit volumes associated with a charging key
 * for a user. It can receive used credit, add allowed credit, and check if
 * there is an update (quota exhausted, etc)
 */
class SessionCredit {
 public:
  struct Usage {
    uint64_t bytes_tx;
    uint64_t bytes_rx;
  };

  SessionCredit();

  SessionCredit(ServiceState start_state);

  SessionCredit(ServiceState start_state, CreditLimitType limit_type);

  SessionCredit(const StoredSessionCredit &marshaled);

  StoredSessionCredit marshal();

  /**
   * get_update_criteria constructs a SessionCreditUpdateCriteria with default
   * values.
   */
  SessionCreditUpdateCriteria get_update_criteria();

  /**
   * add_used_credit increments USED_TX and USED_RX
   * as being recently updated
   */
  void add_used_credit(
      uint64_t used_tx, uint64_t used_rx,
      SessionCreditUpdateCriteria& update_criteria);

  /**
   * reset_reporting_credit resets the REPORTING_* to 0
   * Also marks the session as not in reporting.
   */
  void reset_reporting_credit(SessionCreditUpdateCriteria* uc);

  /**
   * Credit update has failed to the OCS, so mark this credit as failed so it
   * can be cut off accordingly
   */
  void mark_failure(uint32_t code, SessionCreditUpdateCriteria* uc);
  /**
   * receive_credit increments ALLOWED* and moves the REPORTING_* credit to
   * the REPORTED_* credit
   */
  void receive_credit(const GrantedUnits& gsu, SessionCreditUpdateCriteria* uc);

  /**
   * get_update returns a filled-in CreditUsage if an update exists, and a blank
   * one if no update exists. Check has_update before calling.
   * This method also sets the REPORTING_* credit buckets
   */
  SessionCredit::Usage get_usage_for_reporting(
      SessionCreditUpdateCriteria& update_criteria);

  SessionCredit::Usage get_all_unreported_usage_for_reporting(
      SessionCreditUpdateCriteria& update_criteria);

  /**
   * Returns true if either of REPORTING_* buckets are more than 0
   */
  bool is_reporting() const;

  /**
   * Helper function to get the credit in a particular bucket
   */
  uint64_t get_credit(Bucket bucket) const;

  void set_grant_tracking_type(
      GrantTrackingType g_type, SessionCreditUpdateCriteria& uc);

  /**
   * Add credit to the specified bucket. This does not necessarily correspond
   * to allowed or used credit.
   * NOTE: Use only for merging updates into SessionStore
   * @param credit
   * @param bucket
   */
  void add_credit(
      uint64_t credit, Bucket bucket,
      SessionCreditUpdateCriteria& update_criteria);
  /**
   * is_quota_exhausted checks if any of the remaining quota (Allowed - Used)
   * on tx, rx, or tx+rx amounts are under a specific threshold, and depending
   * on the grant_tracking_type_ (which selects which of those thresholds
   * matters), it decides if quota is exhausted. The threshold which those three
   * usages are compare against it is a percentage of the amount of last
   * received grant. So basically the algorithm is: if quota remaining is under
   * a percentage of the last received grant, we mark it as exhausted for that
   * leg (rx/tx/total). If percentage is 1 (100%) then that leg will be marked
   * as exhausted when it gets to the top of its corresponding grant.
   *
   * Quota usage is measured by reporting from pipelined since the last
   * SessionUpdate.
   * Check if the session has exhausted its quota granted since the last report.
   *
   * @param usage_reporting_threshold
   * @return true if quota is exhausted for the session
   */
  bool is_quota_exhausted(float usage_reporting_threshold) const;

  /**
   * A threshold represented as a ratio for triggering usage update before
   * an user completely used up the quota
   * Session manager will send usage update when
   * (available bytes since last update) * USAGE_REPORTING_THRESHOLD >=
   * (used bytes since last update)
   */
  static float USAGE_REPORTING_THRESHOLD;

  /**
   * Set to true to terminate service when the quota of a session is exhausted.
   * An user can still use up to the extra margin.
   * Set to false to allow users to use without any constraint.
   */
  static bool TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED;

 private:
  uint64_t buckets_[MAX_VALUES];
  bool reporting_;
  CreditLimitType credit_limit_type_;
  GrantTrackingType grant_tracking_type_;

 private:
  void log_quota_and_usage() const;

  SessionCredit::Usage get_unreported_usage() const;

  void log_usage_report(SessionCredit::Usage) const;

  GrantTrackingType determine_grant_tracking_type(const GrantedUnits& grant);

  bool compute_quota_exhausted(
      const uint64_t allowed, const uint64_t used, float threshold_ratio,
      const uint64_t grantedUnits) const;

  uint64_t compute_reporting_limit(
      const uint64_t allowed, const uint64_t reported) const;

  void apply_reporting_limits(SessionCredit::Usage& usage);

  uint64_t calculate_allowed_floor(CreditUnit cu, Bucket allowed, Bucket floor);
};

}  // namespace magma
