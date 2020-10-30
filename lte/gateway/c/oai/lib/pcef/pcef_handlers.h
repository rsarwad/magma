/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#pragma once

#include <gmp.h>
#include <netinet/in.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "intertask_interface.h"
#include "common_types.h"
#include "ip_forward_messages_types.h"
#include "spgw_types.h"

struct pcef_create_session_data {
  char msisdn[MSISDN_LENGTH + 1];
  char imeisv[IMEISV_DIGITS_MAX + 1];
  uint8_t imeisv_exists;
  char mcc_mnc[7];
  char imsi_mcc_mnc[7];
  char apn[APN_MAX_LENGTH + 1];
  char sgw_ip[INET_ADDRSTRLEN];
  char uli[14];
  charging_characteristics_t charging_characteristics;
  uint8_t uli_exists;
  uint32_t msisdn_len;
  uint32_t mcc_mnc_len;
  uint32_t imsi_mcc_mnc_len;
  uint32_t ambr_dl;
  uint32_t ambr_ul;
  uint32_t pl;
  uint32_t pci;
  uint32_t pvi;
  uint32_t qci;
};

/**
 * pcef_create_session is an asynchronous call that initiates the UE session in
 * the PCEF and sends an S5 ITTI message to SGW when done.
 * This is a long process, so it needs to by asynchronous
 */
void pcef_create_session(
    spgw_state_t* state, const char* imsi, const char* ip,
    const struct pcef_create_session_data* session_data,
    itti_sgi_create_end_point_response_t sgi_response,
    s5_create_session_request_t bearer_request,
    s_plus_p_gw_eps_bearer_context_information_t* ctx_p);

/**
 * pcef_end_session is a *synchronous* call that ends the UE session in the
 * PCEF and returns true if successful.
 * This may turn asynchronous in the future if it's too long
 */
bool pcef_end_session(char* imsi, char* apn);

/**
 * pcef_send_policy2bearer_binding is an asynchronous call that binds policy
 * rule id to the newly created bearer id for a particular session that is
 * uniquely identified by imsi and default bearer id.
 */
void pcef_send_policy2bearer_binding(
    const char* imsi, uint8_t default_bearer_id, char* policy_rule_name,
    uint8_t eps_bearer_id);

void get_session_req_data(
    spgw_state_t* spgw_state,
    const itti_s11_create_session_request_t* saved_req,
    struct pcef_create_session_data* data);

#ifdef __cplusplus
}
#endif
