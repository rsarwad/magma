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
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include "3gpp_29.274.h"
#include "log.h"
#include "common_defs.h"
#include "intertask_interface.h"
#include "sgw_context_manager.h"
#include "spgw_types.h"
#include "sgw_s8_state.h"
#include "sgw_s8_s11_handlers.h"
#include "s8_client_api.h"

extern task_zmq_ctx_t spgw_app_task_zmq_ctx;

uint32_t sgw_get_new_s1u_teid(sgw_state_t* state) {
  if (state->s1u_teid == 0) {
    state->s1u_teid = INITIAL_SGW_S8_S1U_TEID;
  }
  __sync_fetch_and_add(&state->s1u_teid, 1);
  return state->s1u_teid;
}

uint32_t sgw_get_new_s5s8u_teid(sgw_state_t* state) {
  __sync_fetch_and_add(&state->s5s8u_teid, 1);
  return (state->s5s8u_teid);
}

sgw_eps_bearer_context_information_t* sgw_get_sgw_eps_bearer_context(
    teid_t teid) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  sgw_eps_bearer_context_information_t* sgw_bearer_context_info = NULL;
  hash_table_ts_t* state_imsi_ht = get_spgw_ue_state();

  hashtable_ts_get(
      state_imsi_ht, (const hash_key_t) teid,
      (void**) &sgw_bearer_context_info);
  OAILOG_FUNC_RETURN(LOG_SGW_S8, sgw_bearer_context_info);
}

spgw_ue_context_t* sgw_create_or_get_ue_context(
    sgw_state_t* sgw_state, imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  spgw_ue_context_t* ue_context_p = NULL;
  hashtable_ts_get(
      sgw_state->imsi_ue_context_htbl, (const hash_key_t) imsi64,
      (void**) &ue_context_p);
  if (!ue_context_p) {
    ue_context_p = (spgw_ue_context_t*) calloc(1, sizeof(spgw_ue_context_t));
    if (ue_context_p) {
      LIST_INIT(&ue_context_p->sgw_s11_teid_list);
      hashtable_ts_insert(
          sgw_state->imsi_ue_context_htbl, (const hash_key_t) imsi64,
          (void*) ue_context_p);
    } else {
      OAILOG_ERROR_UE(
          LOG_SGW_S8, imsi64, "Failed to allocate memory for UE context \n");
    }
  }
  OAILOG_FUNC_RETURN(LOG_SGW_S8, ue_context_p);
}

int sgw_update_teid_in_ue_context(
    sgw_state_t* sgw_state, imsi64_t imsi64, teid_t teid) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  spgw_ue_context_t* ue_context_p =
      sgw_create_or_get_ue_context(sgw_state, imsi64);
  if (!ue_context_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Failed to get UE context for sgw_s11_teid " TEID_FMT "\n", teid);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }

  sgw_s11_teid_t* sgw_s11_teid_p =
      (sgw_s11_teid_t*) calloc(1, sizeof(sgw_s11_teid_t));
  if (!sgw_s11_teid_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Failed to allocate memory for sgw_s11_teid:" TEID_FMT "\n", teid);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }

  sgw_s11_teid_p->sgw_s11_teid = teid;
  LIST_INSERT_HEAD(&ue_context_p->sgw_s11_teid_list, sgw_s11_teid_p, entries);
  OAILOG_DEBUG(
      LOG_SGW_S8,
      "Inserted sgw_s11_teid to list of teids of UE context" TEID_FMT "\n",
      teid);
  OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNok);
}

sgw_eps_bearer_context_information_t*
sgw_create_bearer_context_information_in_collection(teid_t teid) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  sgw_eps_bearer_context_information_t* new_sgw_bearer_context_information =
      calloc(1, sizeof(sgw_eps_bearer_context_information_t));

  if (new_sgw_bearer_context_information == NULL) {
    OAILOG_ERROR(
        LOG_SGW_S8,
        "Failed to create new sgw bearer context information object for "
        "sgw_s11_teid " TEID_FMT "\n",
        teid);
    return NULL;
  }
  // Insert the new tunnel with sgw_s11_teid into the hash list.
  hash_table_ts_t* state_imsi_ht = get_sgw_ue_state();
  hashtable_ts_insert(
      state_imsi_ht, (const hash_key_t) teid,
      new_sgw_bearer_context_information);

  OAILOG_DEBUG(
      LOG_SGW_S8,
      "Inserted new sgw eps bearer context into hash list,state_imsi_ht with "
      "key as sgw_s11_teid " TEID_FMT "\n ",
      teid);
  return new_sgw_bearer_context_information;
}

void sgw_s8_handle_s11_create_session_request(
    sgw_state_t* sgw_state,
    const itti_s11_create_session_request_t* const session_req_pP,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p = NULL;
  OAILOG_INFO_UE(
      LOG_SGW_S8, imsi64, "Received S11 CREATE SESSION REQUEST from MME_APP\n");
  sgw_eps_bearer_context_information_t* new_sgw_eps_context = NULL;
  mme_sgw_tunnel_t sgw_s11_tunnel                           = {0};

  increment_counter("sgw_s8_create_session", 1, NO_LABELS);
  if (session_req_pP->rat_type != RAT_EUTRAN) {
    OAILOG_WARNING_UE(
        LOG_SGW_S8, imsi64,
        "Received session request with RAT != RAT_TYPE_EUTRAN: type %d\n",
        session_req_pP->rat_type);
  }
  /*
   * As we are abstracting GTP-C transport, FTeid ip address is useless.
   * We just use the teid to identify MME tunnel. Normally we received either:
   * - ipv4 address if ipv4 flag is set
   * - ipv6 address if ipv6 flag is set
   * - ipv4 and ipv6 if both flags are set
   * Communication between MME and S-GW involves S11 interface so we are
   * expecting S11_MME_GTP_C (11) as interface_type.
   */
  if ((session_req_pP->sender_fteid_for_cp.teid == 0) &&
      (session_req_pP->sender_fteid_for_cp.interface_type != S11_MME_GTP_C)) {
    // MME sent request with teid = 0. This is not valid...
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64, "Received invalid teid \n");
    increment_counter(
        "sgw_s8_create_session", 1, 2, "result", "failure", "cause",
        "sender_fteid_incorrect_parameters");
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }

  sgw_get_new_S11_tunnel_id(&sgw_state->tunnel_id);
  sgw_s11_tunnel.local_teid  = sgw_state->tunnel_id;
  sgw_s11_tunnel.remote_teid = session_req_pP->sender_fteid_for_cp.teid;
  OAILOG_DEBUG_UE(
      LOG_SGW_S8, imsi64,
      "Rx CREATE-SESSION-REQUEST MME S11 teid " TEID_FMT
      "SGW S11 teid " TEID_FMT " APN %s EPS bearer Id %d\n",
      sgw_s11_tunnel.remote_teid, sgw_s11_tunnel.local_teid,
      session_req_pP->apn,
      session_req_pP->bearer_contexts_to_be_created.bearer_contexts[0]
          .eps_bearer_id);

  if (sgw_update_teid_in_ue_context(
          sgw_state, imsi64, sgw_s11_tunnel.local_teid) == RETURNerror) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Failed to update sgw_s11_teid" TEID_FMT " in UE context \n",
        sgw_s11_tunnel.local_teid);
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }

  new_sgw_eps_context = sgw_create_bearer_context_information_in_collection(
      sgw_s11_tunnel.local_teid);
  if (!new_sgw_eps_context) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Could not create new sgw context for create session req message for "
        "mme_s11_teid " TEID_FMT "\n",
        session_req_pP->sender_fteid_for_cp.teid);
    increment_counter(
        "sgw_s8_create_session", 1, 2, "result", "failure", "cause",
        "internal_software_error");
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  memcpy(
      new_sgw_eps_context->imsi.digit, session_req_pP->imsi.digit,
      IMSI_BCD_DIGITS_MAX);
  new_sgw_eps_context->imsi64 = imsi64;

  new_sgw_eps_context->imsi_unauthenticated_indicator = 1;
  new_sgw_eps_context->mme_teid_S11 = session_req_pP->sender_fteid_for_cp.teid;
  new_sgw_eps_context->s_gw_teid_S11_S4 = sgw_s11_tunnel.local_teid;
  new_sgw_eps_context->trxn             = session_req_pP->trxn;

  // Update PDN details
  if (session_req_pP->apn) {
    new_sgw_eps_context->pdn_connection.apn_in_use =
        strdup(session_req_pP->apn);
  } else {
    new_sgw_eps_context->pdn_connection.apn_in_use = strdup("NO APN");
  }
  new_sgw_eps_context->pdn_connection.s_gw_teid_S5_S8_cp =
      sgw_s11_tunnel.local_teid;
  bearer_context_to_be_created_t csr_bearer_context =
      session_req_pP->bearer_contexts_to_be_created.bearer_contexts[0];
  new_sgw_eps_context->pdn_connection.default_bearer =
      csr_bearer_context.eps_bearer_id;

  /* creating an eps bearer entry
   * copy informations from create session request to bearer context information
   */

  eps_bearer_ctxt_p = sgw_cm_create_eps_bearer_ctxt_in_collection(
      &new_sgw_eps_context->pdn_connection, csr_bearer_context.eps_bearer_id);
  if (eps_bearer_ctxt_p == NULL) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64, "Failed to create new EPS bearer entry\n");
    increment_counter(
        "sgw_s8_create_session", 1, 2, "result", "failure", "cause",
        "internal_software_error");
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  eps_bearer_ctxt_p->eps_bearer_qos = csr_bearer_context.bearer_level_qos;
  eps_bearer_ctxt_p->s_gw_teid_S1u_S12_S4_up = sgw_get_new_s1u_teid(sgw_state);
  eps_bearer_ctxt_p->s_gw_teid_S5_S8_up = sgw_get_new_s5s8u_teid(sgw_state);
  csr_bearer_context.s5_s8_u_sgw_fteid.teid =
      eps_bearer_ctxt_p->s_gw_teid_S5_S8_up;
  csr_bearer_context.s5_s8_u_sgw_fteid.ipv4 = 1;
  csr_bearer_context.s5_s8_u_sgw_fteid.ipv4_address =
      sgw_state->sgw_ip_address_S5S8_up;

  send_s8_create_session_request(
      sgw_s11_tunnel.local_teid, session_req_pP, imsi64);
  sgw_display_s11_bearer_context_information(new_sgw_eps_context);
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static int update_bearer_context_info(
    sgw_eps_bearer_context_information_t* sgw_context_p,
    const s5s8_create_session_response_t* const session_rsp_p) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  sgw_eps_bearer_ctxt_t* default_bearer_ctx_p = sgw_cm_get_eps_bearer_entry(
      &sgw_context_p->pdn_connection, session_rsp_p->eps_bearer_id);
  if (!default_bearer_ctx_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, sgw_context_p->imsi64,
        "Failed to get default eps bearer context for context teid " TEID_FMT
        "and bearer_id :%d \n",
        session_rsp_p->context_teid, session_rsp_p->eps_bearer_id);
    /*Rashmi send CSRep */
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  memcpy(&default_bearer_ctx_p->paa, &session_rsp_p->paa, sizeof(paa_t));
  for (uint8_t indx = 0; indx < BEARERS_PER_UE; indx++) {
    s5s8_bearer_context_t s5s8_bearer_context =
        session_rsp_p->bearer_context[indx];
    sgw_eps_bearer_ctxt_t* eps_bearer_ctx_p = sgw_cm_get_eps_bearer_entry(
        &sgw_context_p->pdn_connection, s5s8_bearer_context.eps_bearer_id);
    if (!eps_bearer_ctx_p) {
      OAILOG_ERROR_UE(
          LOG_SGW_S8, sgw_context_p->imsi64,
          "Failed to get eps bearer context for context teid " TEID_FMT
          "and bearer_id :%d \n",
          session_rsp_p->context_teid, s5s8_bearer_context.eps_bearer_id);
      /*Rashmi send CSRep */
      OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
    }
    eps_bearer_ctx_p->p_gw_address_in_use_up.address.ipv4_address =
        s5s8_bearer_context.pgw_s8_up.ipv4_address;
    memcpy(
        &eps_bearer_ctx_p->p_gw_address_in_use_up.address.ipv6_address,
        &s5s8_bearer_context.pgw_s8_up.ipv6_address, sizeof(struct in6_addr));
    eps_bearer_ctx_p->p_gw_teid_S5_S8_up = s5s8_bearer_context.pgw_s8_up.teid;

    memcpy(
        &eps_bearer_ctx_p->eps_bearer_qos, &s5s8_bearer_context.qos,
        sizeof(bearer_qos_t));
  }
  OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNok);
}

static int sgw_s8_send_create_session_response(
    sgw_state_t* sgw_state, sgw_eps_bearer_context_information_t* sgw_context_p,
    const s5s8_create_session_response_t* const session_rsp_p) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  MessageDef* message_p                                         = NULL;
  itti_s11_create_session_response_t* create_session_response_p = NULL;

  message_p = itti_alloc_new_message(TASK_SGW_S8, S11_CREATE_SESSION_RESPONSE);
  if (message_p == NULL) {
    OAILOG_CRITICAL_UE(
        LOG_SGW_S8, sgw_context_p->imsi64,
        "Failed to allocate memory for S11_create_session_response \n");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  if (!sgw_context_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, sgw_context_p->imsi64, "sgw_context_p is NULL \n");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }

  create_session_response_p = &message_p->ittiMsg.s11_create_session_response;
  if (!create_session_response_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, sgw_context_p->imsi64,
        "create_session_response_p is NULL \n");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  create_session_response_p->teid              = sgw_context_p->mme_teid_S11;
  create_session_response_p->cause.cause_value = session_rsp_p->cause;
  create_session_response_p->s11_sgw_fteid.teid =
      sgw_context_p->s_gw_teid_S11_S4;
  create_session_response_p->s11_sgw_fteid.interface_type = S11_SGW_GTP_C;
  create_session_response_p->s11_sgw_fteid.ipv4           = 1;
  create_session_response_p->s11_sgw_fteid.ipv4_address.s_addr =
      spgw_config.sgw_config.ipv4.S11.s_addr;

  if (session_rsp_p->cause == REQUEST_ACCEPTED) {
    sgw_eps_bearer_ctxt_t* default_bearer_ctxt_p = sgw_cm_get_eps_bearer_entry(
        &sgw_context_p->pdn_connection, session_rsp_p->eps_bearer_id);
    if (!default_bearer_ctxt_p) {
      OAILOG_ERROR_UE(
          LOG_SGW_S8, sgw_context_p->imsi64,
          "Failed to get default eps bearer context for sgw_s11_teid " TEID_FMT
          "and bearer_id :%d \n",
          session_rsp_p->context_teid, session_rsp_p->eps_bearer_id);
      OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
    }
    memcpy(
        &create_session_response_p->paa, &default_bearer_ctxt_p->paa,
        sizeof(paa_t));
    create_session_response_p->bearer_contexts_created.num_bearer_context = 1;
    bearer_context_created_t bearer_context =
        create_session_response_p->bearer_contexts_created.bearer_contexts[0];

    bearer_context.s1u_sgw_fteid.teid =
        default_bearer_ctxt_p->s_gw_teid_S1u_S12_S4_up;
    bearer_context.s1u_sgw_fteid.interface_type = S1_U_SGW_GTP_U;
    bearer_context.s1u_sgw_fteid.ipv4           = 1;
    bearer_context.s1u_sgw_fteid.ipv4_address.s_addr =
        sgw_state->sgw_ip_address_S1u_S12_S4_up.s_addr;
    bearer_context.eps_bearer_id = session_rsp_p->eps_bearer_id;
    /*
     * Set the Cause information from bearer context created.
     * "Request accepted" is returned when the GTPv2 entity has accepted a
     * control plane request.
     */
    create_session_response_p->bearer_contexts_created.bearer_contexts[0]
        .cause.cause_value = session_rsp_p->cause;
  } else {
    create_session_response_p->bearer_contexts_marked_for_removal
        .num_bearer_context = 1;
    bearer_context_marked_for_removal_t bearer_context =
        create_session_response_p->bearer_contexts_marked_for_removal
            .bearer_contexts[0];
    bearer_context.cause.cause_value = session_rsp_p->cause;
    bearer_context.eps_bearer_id     = session_rsp_p->eps_bearer_id;
    create_session_response_p->trxn  = sgw_context_p->trxn;
  }
  message_p->ittiMsgHeader.imsi = sgw_context_p->imsi64;
  OAILOG_DEBUG_UE(
      LOG_SGW_S8, sgw_context_p->imsi64,
      "Sending S11 Create Session Response to mme for mme_s11_teid: " TEID_FMT
      "\n",
      create_session_response_p->teid);

  send_msg_to_task(&spgw_app_task_zmq_ctx, TASK_MME_APP, message_p);

  OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNok);
}

void sgw_s8_handle_create_session_response(
    sgw_state_t* sgw_state,
    const s5s8_create_session_response_t* const session_rsp_p,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  if (!session_rsp_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Received null create session response from s8_proxy \n");
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  OAILOG_INFO_UE(
      LOG_SGW_S8, imsi64,
      " Rx S5S8_CREATE_SESSION_RSP for context_teid " TEID_FMT "\n",
      session_rsp_p->context_teid);

  sgw_eps_bearer_context_information_t* sgw_context_p =
      sgw_get_sgw_eps_bearer_context(session_rsp_p->context_teid);
  if (!sgw_context_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Failed to fetch sgw_eps_bearer_context_info from "
        "context_teid " TEID_FMT " \n",
        session_rsp_p->context_teid);
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }

  if (session_rsp_p->cause != REQUEST_ACCEPTED) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Received failed create session response with cause: %d for "
        "context_id: " TEID_FMT "\n",
        session_rsp_p->cause, session_rsp_p->context_teid);
    /*Rashmi send CSRep */
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  // update pdn details received from PGW
  sgw_context_p->pdn_connection.p_gw_teid_S5_S8_cp =
      session_rsp_p->pgw_s8_cp_teid.teid;
  // update bearer context details received from PGW
  update_bearer_context_info(sgw_context_p, session_rsp_p);

  // send Create session response to mme
  sgw_s8_send_create_session_response(sgw_state, sgw_context_p, session_rsp_p);
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

