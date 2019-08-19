/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
 * ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are those
 * of the authors and should not be interpreted as representing official policies,
 * either expressed or implied, of the FreeBSD Project.
 */
//WARNING: Do not include this header directly. Use intertask_interface.h instead.

/*! \file nas_messages_def.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
MESSAGE_DEF(
  NAS_CONNECTION_RELEASE_IND,
  MESSAGE_PRIORITY_MED,
  itti_nas_conn_rel_ind_t,
  nas_conn_rel_ind)
MESSAGE_DEF(
  NAS_UPLINK_DATA_IND,
  MESSAGE_PRIORITY_MED,
  itti_nas_ul_data_ind_t,
  nas_ul_data_ind)
MESSAGE_DEF(
  NAS_DOWNLINK_DATA_CNF,
  MESSAGE_PRIORITY_MED,
  itti_nas_dl_data_cnf_t,
  nas_dl_data_cnf)
MESSAGE_DEF(
  NAS_DOWNLINK_DATA_REJ,
  MESSAGE_PRIORITY_MED,
  itti_nas_dl_data_rej_t,
  nas_dl_data_rej)
MESSAGE_DEF(
  NAS_ERAB_SETUP_REQ,
  MESSAGE_PRIORITY_MED,
  itti_erab_setup_req_t,
  itti_erab_setup_req)
//MESSAGE_DEF(NAS_RAB_RELEASE_REQ,                MESSAGE_PRIORITY_MED,   itti_nas_rab_rel_req_t,          nas_rab_rel_req)

/* NAS layer -> MME app messages */
MESSAGE_DEF(
  NAS_AUTHENTICATION_PARAM_REQ,
  MESSAGE_PRIORITY_MED,
  itti_nas_auth_param_req_t,
  nas_auth_param_req)
MESSAGE_DEF(
  NAS_SGS_DETACH_REQ,
  MESSAGE_PRIORITY_MED,
  itti_nas_sgs_detach_req_t,
  nas_sgs_detach_req)
MESSAGE_DEF(
  NAS_IMPLICIT_DETACH_UE_IND,
  MESSAGE_PRIORITY_MED,
  itti_nas_implicit_detach_ue_ind_t,
  nas_implicit_detach_ue_ind)
MESSAGE_DEF(
  NAS_NW_INITIATED_DETACH_UE_REQ,
  MESSAGE_PRIORITY_MED,
  itti_nas_nw_initiated_detach_ue_req_t,
  nas_nw_initiated_detach_ue_req)
MESSAGE_DEF(
  NAS_EXTENDED_SERVICE_REQ,
  MESSAGE_PRIORITY_MED,
  itti_nas_extended_service_req_t,
  nas_extended_service_req)
MESSAGE_DEF(
  NAS_CS_SERVICE_NOTIFICATION,
  MESSAGE_PRIORITY_MED,
  itti_nas_cs_service_notification_t,
  nas_cs_service_notification)
MESSAGE_DEF(
  NAS_CS_DOMAIN_LOCATION_UPDATE_REQ,
  MESSAGE_PRIORITY_MED,
  itti_nas_cs_domain_location_update_req_t,
  nas_cs_domain_location_update_req)
MESSAGE_DEF(
  NAS_CS_DOMAIN_LOCATION_UPDATE_ACC,
  MESSAGE_PRIORITY_MED,
  itti_nas_cs_domain_location_update_acc_t,
  nas_cs_domain_location_update_acc)
MESSAGE_DEF(
  NAS_CS_DOMAIN_LOCATION_UPDATE_FAIL,
  MESSAGE_PRIORITY_MED,
  itti_nas_cs_domain_location_update_fail_t,
  nas_cs_domain_location_update_fail)
MESSAGE_DEF(
  NAS_TAU_COMPLETE,
  MESSAGE_PRIORITY_MED,
  itti_nas_tau_complete_t,
  nas_tau_complete)
MESSAGE_DEF(
  NAS_NOTIFY_SERVICE_REJECT,
  MESSAGE_PRIORITY_MED,
  itti_nas_notify_service_reject_t,
  nas_notify_service_reject)
