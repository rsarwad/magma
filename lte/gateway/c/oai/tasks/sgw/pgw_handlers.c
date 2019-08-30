/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

/*! \file pgw_handlers.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#define PGW
#define S5_HANDLERS_C

#include <arpa/inet.h>
#include <netinet/in.h>
#include <stdint.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

#include "assertions.h"
#include "intertask_interface.h"
#include "log.h"
#include "sgw.h"
#include "spgw_config.h"
#include "pgw_pco.h"
#include "pgw_ue_ip_address_alloc.h"
#include "pgw_handlers.h"
#include "pcef_handlers.h"
#include "common_defs.h"
#include "3gpp_23.003.h"
#include "3gpp_23.401.h"
#include "3gpp_24.008.h"
#include "3gpp_29.274.h"
#include "common_types.h"
#include "hashtable.h"
#include "intertask_interface_types.h"
#include "ip_forward_messages_types.h"
#include "itti_types.h"
#include "pgw_config.h"
#include "s11_messages_types.h"
#include "service303.h"
#include "sgw_context_manager.h"
#include "sgw_ie_defs.h"

static void get_session_req_data(
  const itti_s11_create_session_request_t *saved_req,
  struct pcef_create_session_data *data);
extern sgw_app_t sgw_app;
extern spgw_config_t spgw_config;
extern uint32_t sgw_get_new_s1u_teid(void);
extern void print_bearer_ids_helper(ebi_t[], uint32_t);
//--------------------------------------------------------------------------------

int pgw_handle_create_bearer_request(
  const itti_s5_create_bearer_request_t *const bearer_req_p)
{
  // assign the IP here and just send back a S5_CREATE_BEARER_RESPONSE
  s_plus_p_gw_eps_bearer_context_information_t *new_bearer_ctxt_info_p = NULL;
  MessageDef *message_p = NULL;
  sgw_eps_bearer_ctxt_t *eps_bearer_entry_p = NULL;
  hashtable_rc_t hash_rc = HASH_TABLE_OK;
  itti_sgi_create_end_point_response_t sgi_create_endpoint_resp = {0};
  struct in_addr inaddr;
  char *imsi = NULL;
  OAILOG_FUNC_IN(LOG_PGW_APP);

  OAILOG_DEBUG(
    LOG_PGW_APP,
    "Rx S5_CREATE_BEARER_REQUEST, Context S-GW S11 teid %u, EPS bearer id %u\n",
    bearer_req_p->context_teid,
    bearer_req_p->eps_bearer_id);
  hash_rc = hashtable_ts_get(
    sgw_app.s11_bearer_context_information_hashtable,
    bearer_req_p->context_teid,
    (void **) &new_bearer_ctxt_info_p);

  if (HASH_TABLE_OK == hash_rc) {
    /*hash_rc = hashtable_ts_get
    (new_bearer_ctxt_info_p->sgw_eps_bearer_context_information.pdn_connection.sgw_eps_bearers_array, bearer_req_p->eps_bearer_id, (void **)&eps_bearer_entry_p);
    DevAssert (HASH_TABLE_OK == hash_rc);*/
    eps_bearer_entry_p =
      new_bearer_ctxt_info_p->sgw_eps_bearer_context_information.pdn_connection
        .sgw_eps_bearers_array[EBI_TO_INDEX(bearer_req_p->eps_bearer_id)];
    OAILOG_DEBUG(
      LOG_PGW_APP,
      "Updated eps_bearer_entry_p eps_b_id %u with SGW S1U teid %u\n",
      bearer_req_p->eps_bearer_id,
      bearer_req_p->S1u_teid);
    eps_bearer_entry_p->s_gw_teid_S1u_S12_S4_up = bearer_req_p->S1u_teid;
    memset(
      &sgi_create_endpoint_resp,
      0,
      sizeof(itti_sgi_create_end_point_response_t));

    //--------------------------------------------------------------------------
    // PCO processing
    //--------------------------------------------------------------------------
    protocol_configuration_options_t *pco_req =
      &new_bearer_ctxt_info_p->sgw_eps_bearer_context_information.saved_message
         .pco;
    protocol_configuration_options_t pco_resp = {0};
    protocol_configuration_options_ids_t pco_ids;
    memset(&pco_ids, 0, sizeof pco_ids);

    // TODO: perhaps change to a nonfatal assert?
    AssertFatal(
      0 == pgw_process_pco_request(pco_req, &pco_resp, &pco_ids),
      "Error in processing PCO in request");
    copy_protocol_configuration_options(
      &sgi_create_endpoint_resp.pco, &pco_resp);
    clear_protocol_configuration_options(&pco_resp);

    //--------------------------------------------------------------------------
    // IP forward will forward packets to this teid
    sgi_create_endpoint_resp.context_teid = bearer_req_p->context_teid;
    sgi_create_endpoint_resp.sgw_S1u_teid = bearer_req_p->S1u_teid;
    sgi_create_endpoint_resp.eps_bearer_id = bearer_req_p->eps_bearer_id;
    // TO DO NOW
    sgi_create_endpoint_resp.paa.pdn_type =
      new_bearer_ctxt_info_p->sgw_eps_bearer_context_information.saved_message
        .pdn_type;

    imsi =
      (char *)
        new_bearer_ctxt_info_p->sgw_eps_bearer_context_information.imsi.digit;

    switch (sgi_create_endpoint_resp.paa.pdn_type) {
      case IPv4:
        // Use NAS by default if no preference is set.
        //
        // For context, the protocol configuration options (PCO) section of the
        // packet from the UE is optional, which means that it is perfectly valid
        // for a UE to send no PCO preferences at all. The previous logic only
        // allocates an IPv4 address if the UE has explicitly set the PCO
        // parameter for allocating IPv4 via NAS signaling (as opposed to via
        // DHCPv4). This means that, in the absence of either parameter being set,
        // the does not know what to do, so we need a default option as well.
        //
        // Since we only support the NAS signaling option right now, we will
        // default to using NAS signaling UNLESS we see a preference for DHCPv4.
        // This means that all IPv4 addresses are now allocated via NAS signaling
        // unless specified otherwise.
        //
        // In the long run, we will want to evolve the logic to use whatever
        // information we have to choose the ``best" allocation method. This means
        // adding new bitfields to pco_ids in pgw_pco.h, setting them in pgw_pco.c
        // and using them here in conditional logic. We will also want to
        // implement different logic between the PDN types.
        if (!pco_ids.ci_ipv4_address_allocation_via_dhcpv4) {
          if (0 == allocate_ue_ipv4_address(imsi, &inaddr)) {
            increment_counter(
              "ue_pdn_connection",
              1,
              2,
              "pdn_type",
              "ipv4",
              "result",
              "success");
            sgi_create_endpoint_resp.paa.ipv4_address = inaddr;
            OAILOG_DEBUG(
              LOG_PGW_APP, "Allocated IPv4 address for imsi <%s>\n", imsi);
            sgi_create_endpoint_resp.status = SGI_STATUS_OK;
          } else {
            increment_counter(
              "ue_pdn_connection",
              1,
              2,
              "pdn_type",
              "ipv4",
              "result",
              "failure");
            OAILOG_ERROR(
              LOG_PGW_APP, "Failed to allocate IPv4 PAA for PDN type IPv4\n");
            sgi_create_endpoint_resp.status =
              SGI_STATUS_ERROR_ALL_DYNAMIC_ADDRESSES_OCCUPIED;
          }
        }

        break;

      case IPv6:
        increment_counter(
          "ue_pdn_connection", 1, 2, "pdn_type", "ipv4v6", "result", "failure");
        OAILOG_ERROR(LOG_PGW_APP, "IPV6 PDN type NOT Supported\n");
        sgi_create_endpoint_resp.status =
          SGI_STATUS_ERROR_SERVICE_NOT_SUPPORTED;

        break;

      case IPv4_AND_v6:
        if (0 == allocate_ue_ipv4_address(imsi, &inaddr)) {
          increment_counter(
            "ue_pdn_connection",
            1,
            2,
            "pdn_type",
            "ipv4v6",
            "result",
            "success");
          sgi_create_endpoint_resp.paa.ipv4_address = inaddr;
          OAILOG_DEBUG(LOG_PGW_APP, "Allocated IPv4 address\n");
          sgi_create_endpoint_resp.status = SGI_STATUS_OK;
          sgi_create_endpoint_resp.paa.pdn_type = IPv4;
        } else {
          increment_counter(
            "ue_pdn_connection",
            1,
            2,
            "pdn_type",
            "ipv4v6",
            "result",
            "failure");
          OAILOG_ERROR(
            LOG_PGW_APP,
            "Failed to allocate IPv4 PAA for PDN type IPv4_AND_v6\n");
          sgi_create_endpoint_resp.status =
            SGI_STATUS_ERROR_ALL_DYNAMIC_ADDRESSES_OCCUPIED;
        }

        break;

      default:
        AssertFatal(
          0, "BAD paa.pdn_type %d", sgi_create_endpoint_resp.paa.pdn_type);
        break;
    }

  } else { // if (HASH_TABLE_OK == hash_rc)
    OAILOG_DEBUG(
      LOG_PGW_APP,
      "Rx S11_S1U_ENDPOINT_CREATED, Context: teid %u NOT FOUND\n",
      bearer_req_p->context_teid);
    sgi_create_endpoint_resp.status = SGI_STATUS_ERROR_CONTEXT_NOT_FOUND;
  }
  if (
    spgw_config.pgw_config.relay_enabled &&
    sgi_create_endpoint_resp.status == SGI_STATUS_OK) {
    // create session in PCEF and return
    char ip_str[INET_ADDRSTRLEN];
    inet_ntop(AF_INET, &(inaddr.s_addr), ip_str, INET_ADDRSTRLEN);
    struct pcef_create_session_data session_data;
    get_session_req_data(
      &new_bearer_ctxt_info_p->sgw_eps_bearer_context_information.saved_message,
      &session_data);
    pcef_create_session(
      imsi, ip_str, &session_data, sgi_create_endpoint_resp, *bearer_req_p);
    OAILOG_FUNC_RETURN(LOG_PGW_APP, RETURNok);
  }
  message_p = itti_alloc_new_message(TASK_SPGW_APP, S5_CREATE_BEARER_RESPONSE);
  itti_s5_create_bearer_response_t *s5_response =
    &message_p->ittiMsg.s5_create_bearer_response;
  memset(s5_response, 0, sizeof(itti_s5_create_bearer_response_t));
  s5_response->context_teid = bearer_req_p->context_teid;
  s5_response->S1u_teid = bearer_req_p->S1u_teid;
  s5_response->eps_bearer_id = bearer_req_p->eps_bearer_id;
  s5_response->sgi_create_endpoint_resp = sgi_create_endpoint_resp;
  s5_response->failure_cause = S5_OK;
  OAILOG_DEBUG(
    LOG_PGW_APP,
    "Sending S5 Create Bearer Response to SPGW APP: Context teid %u, S1U-teid = %u,"
    "EPS Bearer Id = %u\n",
    s5_response->context_teid,
    s5_response->S1u_teid,
    s5_response->eps_bearer_id);
  itti_send_msg_to_task(TASK_SPGW_APP, INSTANCE_DEFAULT, message_p);
  OAILOG_FUNC_RETURN(LOG_PGW_APP, RETURNok);
}

static int get_imeisv_from_session_req(
  const itti_s11_create_session_request_t *saved_req,
  char *imeisv)
{
  if (saved_req->mei.present & MEI_IMEISV) {
    // IMEISV as defined in 3GPP TS 23.003 MEI_IMEISV
    imeisv[0] = saved_req->mei.choice.imeisv.u.num.tac8;
    imeisv[1] = saved_req->mei.choice.imeisv.u.num.tac7;
    imeisv[2] = saved_req->mei.choice.imeisv.u.num.tac6;
    imeisv[3] = saved_req->mei.choice.imeisv.u.num.tac5;
    imeisv[4] = saved_req->mei.choice.imeisv.u.num.tac4;
    imeisv[5] = saved_req->mei.choice.imeisv.u.num.tac3;
    imeisv[6] = saved_req->mei.choice.imeisv.u.num.tac2;
    imeisv[7] = saved_req->mei.choice.imeisv.u.num.tac1;
    imeisv[8] = saved_req->mei.choice.imeisv.u.num.snr6;
    imeisv[9] = saved_req->mei.choice.imeisv.u.num.snr5;
    imeisv[10] = saved_req->mei.choice.imeisv.u.num.snr4;
    imeisv[11] = saved_req->mei.choice.imeisv.u.num.snr3;
    imeisv[12] = saved_req->mei.choice.imeisv.u.num.snr2;
    imeisv[13] = saved_req->mei.choice.imeisv.u.num.snr1;
    imeisv[14] = saved_req->mei.choice.imeisv.u.num.svn2;
    imeisv[15] = saved_req->mei.choice.imeisv.u.num.svn1;
    imeisv[IMEISV_DIGITS_MAX] = '\0';

    return 1;
  }
  return 0;
}

static void get_plmn_from_session_req(
  const itti_s11_create_session_request_t *saved_req,
  char *mcc_mnc)
{
  mcc_mnc[0] = saved_req->serving_network.mcc[0];
  mcc_mnc[1] = saved_req->serving_network.mcc[1];
  mcc_mnc[2] = saved_req->serving_network.mcc[2];
  mcc_mnc[3] = saved_req->serving_network.mnc[0];
  mcc_mnc[4] = saved_req->serving_network.mnc[1];
  if (saved_req->serving_network.mnc[2] != 0xff) {
    mcc_mnc[5] = saved_req->serving_network.mnc[2];
    mcc_mnc[6] = '\0';
  } else {
    mcc_mnc[5] = '\0';
  }
}

static void get_imsi_plmn_from_session_req(
  const itti_s11_create_session_request_t *saved_req,
  char *mcc_mnc)
{
  mcc_mnc[0] = saved_req->imsi.digit[0];
  mcc_mnc[1] = saved_req->imsi.digit[1];
  mcc_mnc[2] = saved_req->imsi.digit[2];
  mcc_mnc[3] = saved_req->imsi.digit[3];
  mcc_mnc[4] = saved_req->imsi.digit[4];
  if ((saved_req->imsi.digit[5] & 0xf) != 0xf) {
    mcc_mnc[5] = saved_req->imsi.digit[5];
    mcc_mnc[6] = '\0';
  } else {
    mcc_mnc[5] = '\0';
  }
}

static int get_uli_from_session_req(
  const itti_s11_create_session_request_t *saved_req,
  char *uli)
{
  if (!saved_req->uli.present) {
    return 0;
  }

  uli[0] = 130; // TAI and ECGI - defined in 29.061

  // TAI as defined in 29.274 8.21.4
  uli[1] = ((saved_req->uli.s.tai.mcc[1] & 0xf) << 4) |
           ((saved_req->uli.s.tai.mcc[0] & 0xf));
  uli[2] = ((saved_req->uli.s.tai.mnc[2] & 0xf) << 4) |
           ((saved_req->uli.s.tai.mcc[2] & 0xf));
  uli[3] = ((saved_req->uli.s.tai.mnc[1] & 0xf) << 4) |
           ((saved_req->uli.s.tai.mnc[0] & 0xf));
  uli[4] = (saved_req->uli.s.tai.tac >> 8) & 0xff;
  uli[5] = saved_req->uli.s.tai.tac & 0xff;

  // ECGI as defined in 29.274 8.21.5
  uli[6] = ((saved_req->uli.s.ecgi.mcc[1] & 0xf) << 4) |
           ((saved_req->uli.s.ecgi.mcc[0] & 0xf));
  uli[7] = ((saved_req->uli.s.ecgi.mnc[2] & 0xf) << 4) |
           ((saved_req->uli.s.ecgi.mcc[2] & 0xf));
  uli[8] = ((saved_req->uli.s.ecgi.mnc[1] & 0xf) << 4) |
           ((saved_req->uli.s.ecgi.mnc[0] & 0xf));
  uli[9] = (saved_req->uli.s.ecgi.eci >> 24) & 0xf;
  uli[10] = (saved_req->uli.s.ecgi.eci >> 16) & 0xff;
  uli[11] = (saved_req->uli.s.ecgi.eci >> 8) & 0xff;
  uli[12] = saved_req->uli.s.ecgi.eci & 0xff;
  uli[13] = '\0';
  return 1;
}

static int get_msisdn_from_session_req(
  const itti_s11_create_session_request_t *saved_req,
  char *msisdn)
{
  int len = saved_req->msisdn.length;
  int i, j;

  for (i = 0; i < len; ++i) {
    j = i << 1;
    msisdn[j] = (saved_req->msisdn.digit[i] & 0xf) + '0';
    msisdn[j + 1] = ((saved_req->msisdn.digit[i] >> 4) & 0xf) + '0';
  }
  if ((saved_req->msisdn.digit[len - 1] & 0xf0) == 0xf0) {
    len = (len << 1) - 1;
  } else {
    len = len << 1;
  }
  return len;
}

static void get_session_req_data(
  const itti_s11_create_session_request_t *saved_req,
  struct pcef_create_session_data *data)
{
  const bearer_qos_t *qos;

  data->msisdn_len = get_msisdn_from_session_req(saved_req, data->msisdn);

  data->imeisv_exists = get_imeisv_from_session_req(saved_req, data->imeisv);
  data->uli_exists = get_uli_from_session_req(saved_req, data->uli);
  get_plmn_from_session_req(saved_req, data->mcc_mnc);
  get_imsi_plmn_from_session_req(saved_req, data->imsi_mcc_mnc);

  memcpy(data->apn, saved_req->apn, APN_MAX_LENGTH + 1);

  inet_ntop(AF_INET, &sgw_app.sgw_ip_address_S1u_S12_S4_up, data->sgw_ip, INET_ADDRSTRLEN);

  // QoS Info
  data->ambr_dl = saved_req->ambr.br_dl;
  data->ambr_ul = saved_req->ambr.br_ul;
  qos = &saved_req->bearer_contexts_to_be_created.bearer_contexts[0]
           .bearer_level_qos;
  data->pl = qos->pl;
  data->pci = qos->pci;
  data->pvi = qos->pvi;
  data->qci = qos->qci;
}

//-----------------------------------------------------------------------------

uint32_t pgw_handle_nw_initiated_bearer_actv_req(
  Imsi_t *imsi,
  ebi_t lbi,
  traffic_flow_template_t *ul_tft,
  traffic_flow_template_t *dl_tft,
  bearer_qos_t *eps_bearer_qos)
{
  OAILOG_FUNC_IN(LOG_PGW_APP);
  MessageDef *message_p = NULL;
  uint32_t i = 0;
  uint32_t rc = RETURNok;
  hash_table_ts_t *hashtblP = NULL;
  uint32_t num_elements = 0;
  s_plus_p_gw_eps_bearer_context_information_t *spgw_ctxt_p = NULL;
  hash_node_t *node = NULL;
  itti_s5_nw_init_actv_bearer_request_t
    *itti_s5_actv_bearer_req = NULL;

  OAILOG_INFO(
    LOG_PGW_APP,
    "Received Create Bearer Req from PCRF with IMSI %s\n",
    imsi->digit);

  message_p =
    itti_alloc_new_message(TASK_SPGW_APP,
      S5_NW_INITIATED_ACTIVATE_BEARER_REQ);
  if (message_p == NULL) {
    OAILOG_ERROR(
    LOG_PGW_APP,
    "itti_alloc_new_message failed for"
    "S5_NW_INITIATED_ACTIVATE_DEDICATED_BEARER_REQ\n");
    OAILOG_FUNC_RETURN(LOG_PGW_APP, RETURNerror);
  }
  itti_s5_actv_bearer_req =
    &message_p->ittiMsg.s5_nw_init_actv_bearer_request;
  //Send ITTI message to SGW
  memset(
    itti_s5_actv_bearer_req,
    0,
    sizeof(itti_s5_nw_init_actv_bearer_request_t));

  //Copy Bearer QoS
  memcpy(
    &itti_s5_actv_bearer_req->eps_bearer_qos,
    eps_bearer_qos,
    sizeof(bearer_qos_t));
  //Copy TFT
  memcpy(
    &itti_s5_actv_bearer_req->tft,
    ul_tft,
    sizeof(traffic_flow_template_t));
  //Assign LBI
  hashtblP = sgw_app.s11_bearer_context_information_hashtable;
  if (!hashtblP) {
    OAILOG_ERROR(LOG_PGW_APP, "There is no UE Context in the SGW context \n");
    OAILOG_FUNC_RETURN(LOG_PGW_APP, RETURNerror);
  }

  //Fetch S11 MME TEID using IMSI and LBI
  while ((num_elements < hashtblP->num_elements) && (i < hashtblP->size)) {
    pthread_mutex_lock(&hashtblP->lock_nodes[i]);
    if (hashtblP->nodes[i] != NULL) {
      node = hashtblP->nodes[i];
    }
    pthread_mutex_unlock(&hashtblP->lock_nodes[i]);
    while (node) {
      num_elements++;
      hashtable_ts_get(
        hashtblP, (const hash_key_t) node->key, (void **) &spgw_ctxt_p);
      if (spgw_ctxt_p != NULL) {
        if (!strncmp((const char *)spgw_ctxt_p
          ->sgw_eps_bearer_context_information.imsi.digit,
          (const char *)imsi->digit, strlen((const char *)imsi->digit))) {
          if (spgw_ctxt_p->sgw_eps_bearer_context_information.
            pdn_connection.default_bearer == lbi) {
            itti_s5_actv_bearer_req->lbi = lbi;
            itti_s5_actv_bearer_req->mme_teid_S11 =
              spgw_ctxt_p->sgw_eps_bearer_context_information.mme_teid_S11;
            break;
          }
        }
      }
      node = node->next;
    }
    i++;
  }
  if (i >= hashtblP->size) {
    OAILOG_ERROR(LOG_PGW_APP, "Could not find LBI/IMSI in SPGW context\n");
    //TODO-Send Rsp to PCRF with cause = REJECTED
    /*rc = send_dedicated_bearer_actv_rsp(lbi,
    REQUEST_REJECTED);*/
    OAILOG_FUNC_RETURN(LOG_PGW_APP, rc);
  }
  //Send S5_ACTIVATE_DEDICATED_BEARER_REQ to SGW APP
  OAILOG_INFO(LOG_PGW_APP, "LBI for the received Create Bearer Req %d\n",
    itti_s5_actv_bearer_req->lbi);
  OAILOG_INFO(LOG_PGW_APP,
    "Sending S5_ACTIVATE_DEDICATED_BEARER_REQ to SGW with MME TEID %d\n",
    itti_s5_actv_bearer_req->mme_teid_S11);
  rc = itti_send_msg_to_task(TASK_SPGW_APP, INSTANCE_DEFAULT, message_p);
  OAILOG_FUNC_RETURN(LOG_PGW_APP, rc);
}

//------------------------------------------------------------------------------

uint32_t pgw_handle_nw_initiated_bearer_deactv_req(
  Imsi_t *imsi,
  uint32_t no_of_bearers,
  ebi_t ebi[])
{
  uint32_t rc = RETURNok;
  OAILOG_FUNC_IN(LOG_PGW_APP);
  MessageDef *message_p = NULL;
  uint32_t i = 0;
  uint32_t j = 0;
  hash_table_ts_t *hashtblP = NULL;
  uint32_t num_elements = 0;
  s_plus_p_gw_eps_bearer_context_information_t *spgw_ctxt_p = NULL;
  hash_node_t *node = NULL;
  itti_s5_nw_init_deactv_bearer_request_t
    *itti_s5_deactv_ded_bearer_req = NULL;
  bool found = false;

  OAILOG_INFO(
    LOG_PGW_APP,
    "Received nw_initiated_deactv_bearer_req from NW\n");
  print_bearer_ids_helper(ebi, no_of_bearers);
  message_p =
    itti_alloc_new_message(TASK_SPGW_APP,
      S5_NW_INITIATED_DEACTIVATE_BEARER_REQ);
  if (message_p == NULL) {
    OAILOG_ERROR(
      LOG_PGW_APP,
      "itti_alloc_new_message failed for nw_initiated_deactv_bearer_req\n");
    OAILOG_FUNC_RETURN(LOG_PGW_APP, RETURNerror);
  }

  itti_s5_deactv_ded_bearer_req =
    &message_p->ittiMsg.s5_nw_init_deactv_bearer_request;
  //Send ITTI message to SGW
  memset(
    itti_s5_deactv_ded_bearer_req,
    0,
    sizeof(itti_s5_nw_init_deactv_bearer_request_t));
  itti_s5_deactv_ded_bearer_req->delete_default_bearer = false;
  itti_s5_deactv_ded_bearer_req->no_of_bearers = no_of_bearers;
  memcpy(
    &itti_s5_deactv_ded_bearer_req->ebi,
    ebi,
    (sizeof(ebi_t)) * no_of_bearers);
  hashtblP = sgw_app.s11_bearer_context_information_hashtable;
  if (hashtblP == NULL) {
    OAILOG_ERROR(
      LOG_PGW_APP,
      "hashtblP is NULL for nw_initiated_deactv_bearer_req\n");
    OAILOG_FUNC_RETURN(LOG_PGW_APP, RETURNerror);
  }

  //Check if EBI recvd == LBI to know if default bearer has to be deactivated
  while ((num_elements < hashtblP->num_elements) && (i < hashtblP->size)
    && (!found)) {
    pthread_mutex_lock(&hashtblP->lock_nodes[i]);
    if (hashtblP->nodes[i] != NULL) {
      node = hashtblP->nodes[i];
      spgw_ctxt_p = node->data;
      pthread_mutex_unlock(&hashtblP->lock_nodes[i]);
      num_elements++;
      if (spgw_ctxt_p != NULL) {
        if (!strcmp((const char *)spgw_ctxt_p->
          sgw_eps_bearer_context_information.imsi.digit,
          (const char *)imsi->digit)) {
          itti_s5_deactv_ded_bearer_req->s11_mme_teid =
            spgw_ctxt_p->sgw_eps_bearer_context_information.mme_teid_S11;
          for (j = 0; j < no_of_bearers; j++) {
            if (ebi[j] == spgw_ctxt_p->sgw_eps_bearer_context_information.
              pdn_connection.default_bearer) {
              itti_s5_deactv_ded_bearer_req->delete_default_bearer = true;
              found = true;
              break;
            }
          }
        }
      }
    }
    i++;
  }
  OAILOG_INFO(
    LOG_PGW_APP,
    "Sending nw_initiated_deactv_bearer_req to SGW"
    "with delete_default_bearer flag set to %d\n",
    itti_s5_deactv_ded_bearer_req->delete_default_bearer);
  rc = itti_send_msg_to_task(TASK_SPGW_APP, INSTANCE_DEFAULT, message_p);

  OAILOG_FUNC_RETURN(LOG_PGW_APP, rc);
}

//------------------------------------------------------------------------------

uint32_t pgw_handle_nw_init_activate_bearer_rsp(
  const itti_s5_nw_init_actv_bearer_rsp_t *const act_ded_bearer_rsp)
{
  uint32_t rc = RETURNok;
  OAILOG_FUNC_IN(LOG_PGW_APP);

  OAILOG_INFO(
    LOG_PGW_APP,
    "Sending Create Bearer Rsp to PCRF with EBI %d\n",
    act_ded_bearer_rsp->ebi);
  //Send Create Bearer Rsp to PCRF
  //TODO-Uncomment once implemented at PCRF
  /*rc = send_dedicated_bearer_actv_rsp(act_ded_bearer_rsp->ebi,
    act_ded_bearer_rsp->cause);*/
  OAILOG_FUNC_RETURN(LOG_PGW_APP, rc);
}

//------------------------------------------------------------------------------

uint32_t pgw_handle_nw_init_deactivate_bearer_rsp(
  const itti_s5_nw_init_deactv_bearer_rsp_t *const deact_ded_bearer_rsp)
{
  uint32_t rc = RETURNok;
  OAILOG_FUNC_IN(LOG_PGW_APP);
  ebi_t ebi[BEARERS_PER_UE];

  memcpy(ebi,deact_ded_bearer_rsp->ebi,deact_ded_bearer_rsp->no_of_bearers);
  print_bearer_ids_helper(ebi, deact_ded_bearer_rsp->no_of_bearers);
  //Send Delete Bearer Rsp to PCRF
  //TODO-Uncomment once implemented at PCRF
  //rc = send_dedicated_bearer_deactv_rsp(deact_ded_bearer_rsp->ebi);
  OAILOG_FUNC_RETURN(LOG_PGW_APP, rc);
}
