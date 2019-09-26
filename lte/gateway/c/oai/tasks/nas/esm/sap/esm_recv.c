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
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#include <stdbool.h>
#include <stdlib.h>

#include "log.h"
#include "dynamic_memory_check.h"
#include "common_types.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"
#include "mme_app_ue_context.h"
#include "esm_recv.h"
#include "esm_pt.h"
#include "esm_ebr.h"
#include "esm_proc.h"
#include "esm_cause.h"
#include "mme_config.h"
#include "3gpp_24.301.h"
#include "3gpp_36.401.h"
#include "NasRequestType.h"
#include "PdnType.h"
#include "common_defs.h"
#include "esm_data.h"
#include "mme_api.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/*
   --------------------------------------------------------------------------
   Functions executed by both the UE and the MME upon receiving ESM messages
   --------------------------------------------------------------------------
*/
/****************************************************************************
 **                                                                        **
 ** Name:    esm_recv_status()                                         **
 **                                                                        **
 ** Description: Processes ESM status message                              **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **      pti:       Procedure transaction identity             **
 **      ebi:       EPS bearer identity                        **
 **      msg:       The received ESM message                   **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    ESM cause code whenever the processing of  **
 **             the ESM message fails                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/

esm_cause_t esm_recv_status(
  emm_context_t *emm_context,
  proc_tid_t pti,
  ebi_t ebi,
  const esm_status_msg *msg)
{
  esm_cause_t esm_cause = ESM_CAUSE_SUCCESS;
  int rc = RETURNerror;

  OAILOG_FUNC_IN(LOG_NAS_ESM);
  OAILOG_INFO(
    LOG_NAS_ESM,
    "ESM-SAP   - Received ESM status message (pti=%d, ebi=%d)\n",
    pti,
    ebi);
  /*
   * Message processing
   */
  /*
   * Get the ESM cause
   */
  esm_cause = msg->esmcause;
  /*
   * Execute the ESM status procedure
   */
  rc = esm_proc_status_ind(emm_context, pti, ebi, &esm_cause);

  if (rc != RETURNerror) {
    esm_cause = ESM_CAUSE_SUCCESS;
  }

  /*
   * Return the ESM cause value
   */
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, esm_cause);
}

/*
   --------------------------------------------------------------------------
   Functions executed by the MME upon receiving ESM message from the UE
   --------------------------------------------------------------------------
*/
/****************************************************************************
 **                                                                        **
 ** Name:    esm_recv_pdn_connectivity_request()                       **
 **                                                                        **
 ** Description: Processes PDN connectivity request message                **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **      pti:       Procedure transaction identity             **
 **      ebi:       EPS bearer identity                        **
 **      msg:       The received ESM message                   **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     new_ebi:   New assigned EPS bearer identity           **
 **      data:      PDN connection and EPS bearer context data **
 **      Return:    ESM cause code whenever the processing of  **
 **             the ESM message fails                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
esm_cause_t esm_recv_pdn_connectivity_request(
  emm_context_t *emm_context,
  proc_tid_t pti,
  ebi_t ebi,
  const pdn_connectivity_request_msg *msg,
  ebi_t *new_ebi)
{
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int esm_cause = ESM_CAUSE_SUCCESS;
  mme_ue_s1ap_id_t ue_id =
    PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
      ->mme_ue_s1ap_id;

  OAILOG_INFO(
    LOG_NAS_ESM,
    "ESM-SAP   - Received PDN Connectivity Request message "
    "(ue_id=%u, pti=%d, ebi=%d)\n",
    ue_id,
    pti,
    ebi);

  /*
   * Procedure transaction identity checking
   */
  if ((pti == ESM_PT_UNASSIGNED) || esm_pt_is_reserved(pti)) {
    /*
     * 3GPP TS 24.301, section 7.3.1, case a
     * * * * Reserved or unassigned PTI value
     */
    OAILOG_ERROR(
      LOG_NAS_ESM, "ESM-SAP   - Invalid PTI value (pti=%d) for (ue_id = %u) \n",
      pti,
      ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_PTI_VALUE);
  }
  /*
   * EPS bearer identity checking
   */
  else if (ebi != ESM_EBI_UNASSIGNED) {
    /*
     * 3GPP TS 24.301, section 7.3.2, case a
     * * * * Reserved or assigned EPS bearer identity value
     */
    OAILOG_ERROR(
      LOG_NAS_ESM, "ESM-SAP   - Invalid EPS bearer identity (ebi=%d) for (ue_id = %u)\n",
      ebi,
      ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_EPS_BEARER_IDENTITY);
  }

  /*
   * Message processing
   */
  /*
   * Get PDN connection and EPS bearer context data structure to setup
   */
  if (!emm_context->esm_ctx.esm_proc_data) {
    emm_context->esm_ctx.esm_proc_data = (esm_proc_data_t *) calloc(
      1, sizeof(*emm_context->esm_ctx.esm_proc_data));
  }

  struct esm_proc_data_s *esm_data = emm_context->esm_ctx.esm_proc_data;

  esm_data->pti = pti;
  /*
   * Get the PDN connectivity request type
   */
  OAILOG_DEBUG(
    LOG_NAS_ESM,
    "ESM-SAP   - PDN Connectivity Request Type = (%d) for (ue_id = %u)\n ",
    msg->requesttype,
    ue_id);

  if (msg->requesttype == REQUEST_TYPE_INITIAL_REQUEST) {
    esm_data->request_type = ESM_PDN_REQUEST_INITIAL;
  } else if (msg->requesttype == REQUEST_TYPE_HANDOVER) {
    esm_data->request_type = ESM_PDN_REQUEST_HANDOVER;
  } else if (msg->requesttype == REQUEST_TYPE_EMERGENCY) {
    esm_data->request_type = ESM_PDN_REQUEST_EMERGENCY;
  } else {
    /*
     * Unkown PDN request type
     */
    esm_data->request_type = -1;
    OAILOG_ERROR(
      LOG_NAS_ESM,
      "ESM-SAP   - Invalid PDN request type (INITIAL/HANDOVER/EMERGENCY) for (ue_id = %u)\n",
      ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_MANDATORY_INFO);
  }

  /*
   * Get the value of the PDN type indicator
   */
  OAILOG_DEBUG(
    LOG_NAS_ESM,
    "ESM-SAP   - PDN Type = (%d) for (ue_id = %u)\n ", msg->pdntype, ue_id);
  if (msg->pdntype == PDN_TYPE_IPV4) {
    esm_data->pdn_type = ESM_PDN_TYPE_IPV4;
  } else if (msg->pdntype == PDN_TYPE_IPV6) {
    esm_data->pdn_type = ESM_PDN_TYPE_IPV6;
  } else if (msg->pdntype == PDN_TYPE_IPV4V6) {
    esm_data->pdn_type = ESM_PDN_TYPE_IPV4V6;
  } else {
    /*
     * Unkown PDN type
     */
    OAILOG_ERROR(LOG_NAS_ESM, "ESM-SAP   - Invalid PDN type for (ue_id = %u)\n", ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_UNKNOWN_PDN_TYPE);
  }

  /*
   * Get the Access Point Name, if provided
   */
  if (msg->presencemask & PDN_CONNECTIVITY_REQUEST_ACCESS_POINT_NAME_PRESENT) {
    if (esm_data->apn) bdestroy_wrapper(&esm_data->apn);
    esm_data->apn = msg->accesspointname;
  }

  if (
    msg->presencemask &
    PDN_CONNECTIVITY_REQUEST_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT) {
    if (esm_data->pco.num_protocol_or_container_id)
      clear_protocol_configuration_options(&esm_data->pco);
    copy_protocol_configuration_options(
      &esm_data->pco, &msg->protocolconfigurationoptions);
  }
  /*
   * Get the ESM information transfer flag
   */
  if (
    msg->presencemask &
    PDN_CONNECTIVITY_REQUEST_ESM_INFORMATION_TRANSFER_FLAG_PRESENT) {
    /*
     * 3GPP TS 24.301, sections 6.5.1.2, 6.5.1.3
     * * * * ESM information, i.e. protocol configuration options, APN, or both,
     * * * * has to be sent after the NAS signalling security has been activated
     * * * * between the UE and the MME.
     * * * *
     * * * * The MME then at a later stage in the PDN connectivity procedure
     * * * * initiates the ESM information request procedure in which the UE
     * * * * can provide the MME with protocol configuration options or APN
     * * * * or both.
     * * * * The MME waits for completion of the ESM information request
     * * * * procedure before proceeding with the PDN connectivity procedure.
     */
    if (!mme_config.nas_config.disable_esm_information) {
      esm_proc_esm_information_request(emm_context, pti);
      OAILOG_FUNC_RETURN(LOG_NAS_ESM, esm_cause);
    }
  }

  OAILOG_DEBUG(
    LOG_NAS_ESM,
    "ESM-PROC  - _esm_data.conf.features %08x, esm pdn type = %d\n",
    _esm_data.conf.features,
    esm_data->pdn_type);
  emm_context->emm_cause = ESM_CAUSE_SUCCESS;
  switch (_esm_data.conf.features & (MME_API_IPV4 | MME_API_IPV6)) {
    case (MME_API_IPV4 | MME_API_IPV6):
      /*
     * The network supports both IPv4 and IPv6 connection
     */
      if (
        (esm_data->pdn_type == ESM_PDN_TYPE_IPV4V6) &&
        (_esm_data.conf.features & MME_API_SINGLE_ADDR_BEARERS)) {
        /*
         * The network supports single IP version bearers only
         */
        emm_context->emm_cause = ESM_CAUSE_SINGLE_ADDRESS_BEARERS_ONLY_ALLOWED;
      }
      break;

    case MME_API_IPV6:
      /*
       * The network supports connection to IPv6 only
       */
      if (esm_data->pdn_type != ESM_PDN_TYPE_IPV6) {
        emm_context->emm_cause = ESM_CAUSE_PDN_TYPE_IPV6_ONLY_ALLOWED;
      }
      break;

    case MME_API_IPV4:
      /*
       * The network supports connection to IPv4 only
       */
      if (esm_data->pdn_type != ESM_PDN_TYPE_IPV4) {
        emm_context->emm_cause = ESM_CAUSE_PDN_TYPE_IPV4_ONLY_ALLOWED;
      }
      break;

    default:
      OAILOG_ERROR(
        LOG_NAS_ESM,
        "ESM-PROC  - _esm_data.conf.features incorrect value (no IPV4 or IPV6 "
        ") %X for (ue_id = %u)\n",
        _esm_data.conf.features,
        ue_id);
  }

#if ORIGINAL_CODE
  /*
   * Execute the PDN connectivity procedure requested by the UE
   */
  int pid = esm_proc_pdn_connectivity_request(
    emm_context,
    pti,
    request_type,
    &esm_data->apn,
    esm_data->pdn_type,
    &esm_data->pdn_addr,
    &esm_data->qos,
    &esm_cause);

  if (pid != RETURNerror) {
    /*
     * Create local default EPS bearer context
     */
    int rc = esm_proc_default_eps_bearer_context(
      emm_ctx, pid, new_ebi, &esm_data->qos, &esm_cause);

    if (rc != RETURNerror) {
      esm_cause = ESM_CAUSE_SUCCESS;
    }
  }
#else
    mme_app_send_s6a_update_location_req(
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context));
  esm_cause = ESM_CAUSE_SUCCESS;
#endif
  /*
   * Return the ESM cause value
   */
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, esm_cause);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_recv_pdn_disconnect_request()                         **
 **                                                                        **
 ** Description: Processes PDN disconnect request message                  **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **      pti:       Procedure transaction identity             **
 **      ebi:       EPS bearer identity                        **
 **      msg:       The received ESM message                   **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     linked_ebi:    Linked EPS bearer identity of the default  **
 **             bearer associated with the PDN to discon-  **
 **             nect from                                  **
 **      Return:    ESM cause code whenever the processing of  **
 **             the ESM message fails                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
esm_cause_t esm_recv_pdn_disconnect_request(
  emm_context_t *emm_context,
  proc_tid_t pti,
  ebi_t ebi,
  const pdn_disconnect_request_msg *msg,
  ebi_t *linked_ebi)
{
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  esm_cause_t esm_cause = ESM_CAUSE_SUCCESS;
  mme_ue_s1ap_id_t ue_id =
    PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
      ->mme_ue_s1ap_id;

  OAILOG_INFO(
    LOG_NAS_ESM,
    "ESM-SAP   - Received PDN Disconnect Request message "
    "(ue_id=%d, pti=%d, ebi=%d)\n",
    ue_id,
    pti,
    ebi);

  /*
   * Procedure transaction identity checking
   */
  if ((pti == ESM_PT_UNASSIGNED) || esm_pt_is_reserved(pti)) {
    /*
     * 3GPP TS 24.301, section 7.3.1, case b
     * * * * Reserved or unassigned PTI value
     */
    OAILOG_WARNING(
      LOG_NAS_ESM, "ESM-SAP   - Invalid PTI value (pti=%d)\n", pti);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_PTI_VALUE);
  }
  /*
   * EPS bearer identity checking
   */
  else if (ebi != ESM_EBI_UNASSIGNED) {
    /*
     * 3GPP TS 24.301, section 7.3.2, case b
     * * * * Reserved or assigned EPS bearer identity value
     */
    OAILOG_WARNING(
      LOG_NAS_ESM, "ESM-SAP   - Invalid EPS bearer identity (ebi=%d)\n", ebi);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_EPS_BEARER_IDENTITY);
  }

  /*
   * Message processing
   */
  /*
   * Execute the PDN disconnect procedure requested by the UE
   */
  int pid = esm_proc_pdn_disconnect_request(emm_context, pti, &esm_cause);

  if (pid != RETURNerror) {
    /*
     * Get the identity of the default EPS bearer context assigned to
     * * * * the PDN connection to disconnect from
     */
    *linked_ebi = msg->linkedepsbeareridentity;
    /*
     * Release the associated default EPS bearer context
     */
    int bid = 0;
    int rc = esm_proc_eps_bearer_context_deactivate(
      emm_context, false, *linked_ebi, &pid, &bid, &esm_cause);

    if (rc != RETURNerror) {
      esm_cause = ESM_CAUSE_SUCCESS;
    }
  }

  /*
   * Return the ESM cause value
   */
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, esm_cause);
}

//------------------------------------------------------------------------------
esm_cause_t esm_recv_information_response(
  emm_context_t *emm_context,
  proc_tid_t pti,
  ebi_t ebi,
  const esm_information_response_msg *msg)
{
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  esm_cause_t esm_cause = ESM_CAUSE_SUCCESS;
  mme_ue_s1ap_id_t ue_id =
    PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
      ->mme_ue_s1ap_id;

  OAILOG_INFO(
    LOG_NAS_ESM,
    "ESM-SAP   - Received ESM Information response message "
    "(ue_id=%d, pti=%d, ebi=%d)\n",
    ue_id,
    pti,
    ebi);

  /*
   * Procedure transaction identity checking
   */
  if ((pti == ESM_PT_UNASSIGNED) || esm_pt_is_reserved(pti)) {
    /*
     * 3GPP TS 24.301, section 7.3.1, case b
     * * * * Reserved or unassigned PTI value
     */
    OAILOG_WARNING(
      LOG_NAS_ESM, "ESM-SAP   - Invalid PTI value (pti=%d)\n", pti);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_PTI_VALUE);
  }
  /*
   * EPS bearer identity checking
   */
  else if (ebi != ESM_EBI_UNASSIGNED) {
    /*
     * 3GPP TS 24.301, section 7.3.2, case b
     * * * * Reserved or assigned EPS bearer identity value
     */
    OAILOG_WARNING(
      LOG_NAS_ESM, "ESM-SAP   - Invalid EPS bearer identity (ebi=%d)\n", ebi);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_EPS_BEARER_IDENTITY);
  }

  /*
   * Message processing
   */
  /*
   * Execute the PDN disconnect procedure requested by the UE
   */
  int pid = esm_proc_esm_information_response(
    emm_context,
    pti,
    msg->accesspointname,
    &msg->protocolconfigurationoptions,
    &esm_cause);

  if (pid != RETURNerror) {
    // Continue with pdn connectivity request
    mme_app_send_s6a_update_location_req(
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context));
    esm_cause = ESM_CAUSE_SUCCESS;
  }

  /*
   * Return the ESM cause value
   */
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, esm_cause);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_recv_activate_default_eps_bearer_context_accept()     **
 **                                                                        **
 ** Description: Processes Activate Default EPS Bearer Context Accept      **
 **      message                                                   **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **          pti:       Procedure transaction identity             **
 **      ebi:       EPS bearer identity                        **
 **      msg:       The received ESM message                   **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    ESM cause code whenever the processing of  **
 **             the ESM message fails                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
esm_cause_t esm_recv_activate_default_eps_bearer_context_accept(
  emm_context_t *emm_context,
  proc_tid_t pti,
  ebi_t ebi,
  const activate_default_eps_bearer_context_accept_msg *msg)
{
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  esm_cause_t esm_cause = ESM_CAUSE_SUCCESS;
  mme_ue_s1ap_id_t ue_id =
    PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
      ->mme_ue_s1ap_id;

  OAILOG_INFO(
    LOG_NAS_ESM,
    "ESM-SAP   - Received Activate Default EPS Bearer Context "
    "Accept message (ue_id=%d, pti=%d, ebi=%d)\n",
    ue_id,
    pti,
    ebi);

  /*
   * Procedure transaction identity checking
   */
  if (esm_pt_is_reserved(pti)) {
    /*
     * 3GPP TS 24.301, section 7.3.1, case f
     * * * * Reserved PTI value
     */
    OAILOG_WARNING(
      LOG_NAS_ESM, "ESM-SAP   - Invalid PTI value (pti=%d)\n", pti);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_PTI_VALUE);
  }
  /*
   * EPS bearer identity checking
   */
  else if (
    esm_ebr_is_reserved(ebi) || esm_ebr_is_not_in_use(emm_context, ebi)) {
    /*
     * 3GPP TS 24.301, section 7.3.2, case f
     * * * * Reserved or assigned value that does not match an existing EPS
     * * * * bearer context
     */
    OAILOG_WARNING(
      LOG_NAS_ESM, "ESM-SAP   - Invalid EPS bearer identity (ebi=%d)\n", ebi);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_EPS_BEARER_IDENTITY);
  }

  /*
   * Message processing
   */
  /*
   * Execute the default EPS bearer context activation procedure accepted
   * * * * by the UE
   */
  int rc =
    esm_proc_default_eps_bearer_context_accept(emm_context, ebi, &esm_cause);

  if (rc != RETURNerror) {
    esm_cause = ESM_CAUSE_SUCCESS;
  }

  /*
   * Return the ESM cause value
   */
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, esm_cause);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_recv_activate_default_eps_bearer_context_reject()     **
 **                                                                        **
 ** Description: Processes Activate Default EPS Bearer Context Reject      **
 **      message                                                   **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **          pti:       Procedure transaction identity             **
 **      ebi:       EPS bearer identity                        **
 **      msg:       The received ESM message                   **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    ESM cause code whenever the processing of  **
 **             the ESM message fail                       **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
esm_cause_t esm_recv_activate_default_eps_bearer_context_reject(
  emm_context_t *emm_context,
  proc_tid_t pti,
  ebi_t ebi,
  const activate_default_eps_bearer_context_reject_msg *msg)
{
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  esm_cause_t esm_cause = ESM_CAUSE_SUCCESS;
  mme_ue_s1ap_id_t ue_id =
    PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
      ->mme_ue_s1ap_id;

  OAILOG_INFO(
    LOG_NAS_ESM,
    "ESM-SAP   - Received Activate Default EPS Bearer Context "
    "Reject message (ue_id=%d, pti=%d, ebi=%d)\n",
    ue_id,
    pti,
    ebi);

  /*
   * Procedure transaction identity checking
   */
  if (esm_pt_is_reserved(pti)) {
    /*
     * 3GPP TS 24.301, section 7.3.1, case f
     * * * * Reserved PTI value
     */
    OAILOG_WARNING(
      LOG_NAS_ESM, "ESM-SAP   - Invalid PTI value (pti=%d)\n", pti);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_PTI_VALUE);
  }
  /*
   * EPS bearer identity checking
   */
  else if (
    esm_ebr_is_reserved(ebi) || esm_ebr_is_not_in_use(emm_context, ebi)) {
    /*
     * 3GPP TS 24.301, section 7.3.2, case f
     * * * * Reserved or assigned value that does not match an existing EPS
     * * * * bearer context
     */
    OAILOG_WARNING(
      LOG_NAS_ESM, "ESM-SAP   - Invalid EPS bearer identity (ebi=%d)", ebi);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_EPS_BEARER_IDENTITY);
  }

  /*
   * Message processing
   */
  /*
   * Execute the default EPS bearer context activation procedure not accepted
   * * * * by the UE
   */
  int rc =
    esm_proc_default_eps_bearer_context_reject(emm_context, ebi, &esm_cause);

  if (rc != RETURNerror) {
    esm_cause = ESM_CAUSE_SUCCESS;
  }

  /*
   * Return the ESM cause value
   */
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, esm_cause);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_recv_activate_dedicated_eps_bearer_context_accept()   **
 **                                                                        **
 ** Description: Processes Activate Dedicated EPS Bearer Context Accept    **
 **      message                                                   **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **          pti:       Procedure transaction identity             **
 **      ebi:       EPS bearer identity                        **
 **      msg:       The received ESM message                   **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    ESM cause code whenever the processing of  **
 **             the ESM message fails                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
esm_cause_t esm_recv_activate_dedicated_eps_bearer_context_accept(
  emm_context_t *emm_context,
  proc_tid_t pti,
  ebi_t ebi,
  const activate_dedicated_eps_bearer_context_accept_msg *msg)
{
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  esm_cause_t esm_cause = ESM_CAUSE_SUCCESS;
  mme_ue_s1ap_id_t ue_id =
    PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
      ->mme_ue_s1ap_id;

  OAILOG_INFO(
    LOG_NAS_ESM,
    "ESM-SAP   - Received Activate Dedicated EPS Bearer "
    "Context Accept message (ue_id=%d, pti=%d, ebi=%d)\n",
    ue_id,
    pti,
    ebi);

  /*
   * Procedure transaction identity checking
   */
  if (esm_pt_is_reserved(pti)) {
    /*
     * 3GPP TS 24.301, section 7.3.1, case f
     * * * * Reserved PTI value
     */
    OAILOG_WARNING(
      LOG_NAS_ESM, "ESM-SAP   - Invalid PTI value (pti=%d)\n", pti);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_PTI_VALUE);
  }
  /*
   * EPS bearer identity checking
   */
  else if (
    esm_ebr_is_reserved(ebi) || esm_ebr_is_not_in_use(emm_context, ebi)) {
    /*
     * 3GPP TS 24.301, section 7.3.2, case f
     * * * * Reserved or assigned value that does not match an existing EPS
     * * * * bearer context
     */
    OAILOG_WARNING(
      LOG_NAS_ESM, "ESM-SAP   - Invalid EPS bearer identity (ebi=%d)\n", ebi);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_EPS_BEARER_IDENTITY);
  }

  /*
   * Message processing
   */
  /*
   * Execute the dedicated EPS bearer context activation procedure accepted
   * * * * by the UE
   */
  int rc =
    esm_proc_dedicated_eps_bearer_context_accept(emm_context, ebi, &esm_cause);

  if (rc != RETURNerror) {
    esm_cause = ESM_CAUSE_SUCCESS;
  }

  /*
   * Return the ESM cause value
   */
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, esm_cause);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_recv_activate_dedicated_eps_bearer_context_reject()   **
 **                                                                        **
 ** Description: Processes Activate Dedicated EPS Bearer Context Reject    **
 **      message                                                   **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **          pti:       Procedure transaction identity             **
 **      ebi:       EPS bearer identity                        **
 **      msg:       The received ESM message                   **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    ESM cause code whenever the processing of  **
 **             the ESM message fail                       **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
esm_cause_t esm_recv_activate_dedicated_eps_bearer_context_reject(
  emm_context_t *emm_context,
  proc_tid_t pti,
  ebi_t ebi,
  const activate_dedicated_eps_bearer_context_reject_msg *msg)
{
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  esm_cause_t esm_cause = ESM_CAUSE_SUCCESS;
  mme_ue_s1ap_id_t ue_id =
    PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
      ->mme_ue_s1ap_id;

  OAILOG_INFO(
    LOG_NAS_ESM,
    "ESM-SAP   - Received Activate Dedicated EPS Bearer "
    "Context Reject message (ue_id=%d, pti=%d, ebi=%d)\n",
    ue_id,
    pti,
    ebi);

  /*
   * Procedure transaction identity checking
   */
  if (esm_pt_is_reserved(pti)) {
    /*
     * 3GPP TS 24.301, section 7.3.1, case f
     * * * * Reserved PTI value
     */
    OAILOG_WARNING(
      LOG_NAS_ESM, "ESM-SAP   - Invalid PTI value (pti=%d)\n", pti);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_PTI_VALUE);
  }
  /*
   * EPS bearer identity checking
   */
  else if (
    esm_ebr_is_reserved(ebi) || esm_ebr_is_not_in_use(emm_context, ebi)) {
    /*
     * 3GPP TS 24.301, section 7.3.2, case f
     * * * * Reserved or assigned value that does not match an existing EPS
     * * * * bearer context
     */
    OAILOG_WARNING(
      LOG_NAS_ESM, "ESM-SAP   - Invalid EPS bearer identity (ebi=%d)\n", ebi);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_EPS_BEARER_IDENTITY);
  }

  /*
   * Message processing
   */
  /*
   * Execute the dedicated EPS bearer context activation procedure not
   * * * *  accepted by the UE
   */
  int rc =
    esm_proc_dedicated_eps_bearer_context_reject(emm_context, ebi, &esm_cause);

  if (rc != RETURNerror) {
    esm_cause = ESM_CAUSE_SUCCESS;
  }

  /*
   * Return the ESM cause value
   */
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, esm_cause);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_recv_deactivate_eps_bearer_context_accept()           **
 **                                                                        **
 ** Description: Processes Deactivate EPS Bearer Context Accept message    **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **          pti:       Procedure transaction identity             **
 **      ebi:       EPS bearer identity                        **
 **      msg:       The received ESM message                   **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    ESM cause code whenever the processing of  **
 **             the ESM message fails                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
esm_cause_t esm_recv_deactivate_eps_bearer_context_accept(
  emm_context_t *emm_context,
  proc_tid_t pti,
  ebi_t ebi,
  const deactivate_eps_bearer_context_accept_msg *msg)
{
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  esm_cause_t esm_cause = ESM_CAUSE_SUCCESS;
  mme_ue_s1ap_id_t ue_id =
    PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
      ->mme_ue_s1ap_id;

  OAILOG_INFO(
    LOG_NAS_ESM,
    "ESM-SAP   - Received Deactivate EPS Bearer Context "
    "Accept message (ue_id=%d, pti=%d, ebi=%d)\n",
    ue_id,
    pti,
    ebi);

  /*
   * Procedure transaction identity checking
   */
  if (esm_pt_is_reserved(pti)) {
    /*
     * 3GPP TS 24.301, section 7.3.1, case f
     * * * * Reserved PTI value
     */
    OAILOG_WARNING(
      LOG_NAS_ESM, "ESM-SAP   - Invalid PTI value (pti=%d)\n", pti);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_PTI_VALUE);
  }
  /*
   * EPS bearer identity checking
   */
  else if (
    esm_ebr_is_reserved(ebi) || esm_ebr_is_not_in_use(emm_context, ebi)) {
    /*
     * 3GPP TS 24.301, section 7.3.2, case f
     * * * * Reserved or assigned value that does not match an existing EPS
     * * * * bearer context
     */
    OAILOG_WARNING(
      LOG_NAS_ESM, "ESM-SAP   - Invalid EPS bearer identity (ebi=%d)\n", ebi);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_EPS_BEARER_IDENTITY);
  }

  /*
   * Message processing
   */
  /*
   * Execute the dedicated EPS bearer context deactivation procedure accepted
   * * * * by the UE
   */
  esm_proc_eps_bearer_context_deactivate_accept(emm_context, ebi, &esm_cause);

  /*
   * Return the ESM cause value
   */
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, esm_cause);
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/
