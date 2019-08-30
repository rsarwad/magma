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

/*! \file nas_messages_types.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_NAS_MESSAGES_TYPES_SEEN
#define FILE_NAS_MESSAGES_TYPES_SEEN

#include <stdint.h>

#include "3gpp_23.003.h"
#include "3gpp_29.274.h"
#include "nas/as_message.h"
#include "common_ies.h"
#include "nas/networkDef.h"

#define NAS_UL_DATA_IND(mSGpTR) (mSGpTR)->ittiMsg.nas_ul_data_ind
#define NAS_DL_DATA_REQ(mSGpTR) (mSGpTR)->ittiMsg.nas_dl_data_req
#define NAS_BEARER_PARAM(mSGpTR) (mSGpTR)->ittiMsg.nas_bearer_param
#define NAS_AUTHENTICATION_REQ(mSGpTR) (mSGpTR)->ittiMsg.nas_auth_req
#define NAS_SGS_DETACH_REQ(mSGpTR) (mSGpTR)->ittiMsg.nas_sgs_detach_req
#define NAS_ERAB_SETUP_REQ(mSGpTR) (mSGpTR)->ittiMsg.itti_erab_setup_req
#define NAS_IMPLICIT_DETACH_UE_IND(mSGpTR)                                     \
  (mSGpTR)->ittiMsg.nas_implicit_detach_ue_ind
#define NAS_CS_DOMAIN_LOCATION_UPDATE_FAIL(mSGpTR)                             \
  (mSGpTR)->ittiMsg.nas_cs_domain_location_update_fail
#define NAS_CS_SERVICE_NOTIFICATION(mSGpTR)                                    \
  (mSGpTR)->ittiMsg.nas_cs_service_notification
#define NAS_DATA_LENGHT_MAX 256
#define NAS_EXTENDED_SERVICE_REQ(mSGpTR)                                       \
  (mSGpTR)->ittiMsg.nas_extended_service_req
#define NAS_NOTIFY_SERVICE_REJECT(mSGpTR)                                      \
  (mSGpTR)->ittiMsg.nas_notify_service_reject
#define NAS_ERAB_REL_CMD(mSGpTR) (mSGpTR)->ittiMsg.itti_erab_rel_cmd

typedef struct itti_nas_cs_service_notification_s {
  mme_ue_s1ap_id_t ue_id; /* UE lower layer identifier        */
#define NAS_PAGING_ID_IMSI 0X00
#define NAS_PAGING_ID_TMSI 0X01
  uint8_t paging_id; /* Paging UE ID, to be sent in CS Service Notification */
  bstring
    cli; /* If CLI received in Sgsap-Paging_Req,shall sent in CS Service Notification */
} itti_nas_cs_service_notification_t;
typedef struct itti_nas_conn_est_rej_s {
  mme_ue_s1ap_id_t ue_id;    /* UE lower layer identifier   */
  s_tmsi_t s_tmsi;           /* UE identity                 */
  nas_error_code_t err_code; /* Transaction status          */
  bstring nas_msg;           /* NAS message to transfer     */
  uint32_t nas_ul_count;     /* UL NAS COUNT                */
  uint16_t selected_encryption_algorithm;
  uint16_t selected_integrity_algorithm;
} itti_nas_conn_est_rej_t;
typedef struct itti_nas_conn_rel_ind_s {
} itti_nas_conn_rel_ind_t;

typedef struct itti_nas_info_transfer_s {
  mme_ue_s1ap_id_t ue_id; /* UE lower layer identifier        */
  //nas_error_code_t err_code;     /* Transaction status               */
  bstring nas_msg; /* Uplink NAS message           */
} itti_nas_info_transfer_t;

typedef struct itti_nas_ul_data_ind_s {
  mme_ue_s1ap_id_t ue_id; /* UE lower layer identifier        */
  bstring nas_msg;        /* Uplink NAS message           */
  tai_t
    tai; /* Indicating the Tracking Area from which the UE has sent the NAS message.  */
  ecgi_t
    cgi; /* Indicating the cell from which the UE has sent the NAS message.   */
} itti_nas_ul_data_ind_t;


typedef struct itti_erab_setup_req_s {
  mme_ue_s1ap_id_t ue_id; /* UE lower layer identifier   */
  ebi_t ebi;              /* EPS bearer id        */
  bstring nas_msg; /* NAS erab bearer context activation message           */
  bitrate_t mbr_dl;
  bitrate_t mbr_ul;
  bitrate_t gbr_dl;
  bitrate_t gbr_ul;
} itti_erab_setup_req_t;

typedef struct itti_erab_rel_cmd_s {
  mme_ue_s1ap_id_t ue_id; /* UE lower layer identifier   */
  ebi_t ebi;              /* EPS bearer id        */
  bstring nas_msg; /* NAS erab bearer context activation message           */
} itti_erab_rel_cmd_t;

typedef struct itti_nas_implicit_detach_ue_ind_s {
  /* UE identifier */
  mme_ue_s1ap_id_t ue_id;
} itti_nas_implicit_detach_ue_ind_t;

typedef struct itti_nas_extended_service_req_s {
  /* UE identifier */
  mme_ue_s1ap_id_t ue_id;
  uint8_t servType; /* service type */
  /* csfb_response is valid only if service type Mobile Terminating CSFB */
  uint8_t csfb_response;
} itti_nas_extended_service_req_t;

typedef struct itti_nas_sgs_detach_req_s {
  /* UE identifier */
  mme_ue_s1ap_id_t ue_id;
  /* detach type */
  uint8_t detach_type;
} itti_nas_sgs_detach_req_t;

typedef struct itti_nas_cs_domain_location_update_fail_s {
/* UE identifier */
#define LAI (1 << 0)
  uint8_t presencemask;
  mme_ue_s1ap_id_t ue_id;
  int reject_cause;
  lai_t laicsfb;
} itti_nas_cs_domain_location_update_fail_t;

/* ITTI message used to intimate service reject for ongoing service request procedure
 * from mme_app to nas
 */
typedef struct itti_nas_notify_service_reject_s {
  mme_ue_s1ap_id_t ue_id;
  uint8_t emm_cause;
#define INTIAL_CONTEXT_SETUP_PROCEDURE_FAILED 0x00
#define UE_CONTEXT_MODIFICATION_PROCEDURE_FAILED 0x01
#define MT_CALL_CANCELLED_BY_NW_IN_IDLE_STATE 0x02
#define MT_CALL_CANCELLED_BY_NW_IN_CONNECTED_STATE 0x03
  uint8_t failed_procedure;
} itti_nas_notify_service_reject_t;

#endif /* FILE_NAS_MESSAGES_TYPES_SEEN */
