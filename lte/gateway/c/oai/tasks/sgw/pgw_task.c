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

/*! \file sgw_task.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#define PGW
#define PGW_TASK_C

#include <stdio.h>
#include <sys/types.h>

#include "log.h"
#include "intertask_interface.h"
#include "pgw_defs.h"
#include "pgw_handlers.h"
#include "sgw.h"
#include "common_defs.h"
#include "bstrlib.h"
#include "intertask_interface_types.h"
#include "spgw_config.h"

pgw_app_t pgw_app;

extern __pid_t g_pid;

static void pgw_exit(void);

//------------------------------------------------------------------------------
static void *pgw_intertask_interface(void *args_p)
{
  itti_mark_task_ready(TASK_PGW_APP);

  while (1) {
    MessageDef *received_message_p = NULL;

    itti_receive_msg(TASK_PGW_APP, &received_message_p);

    switch (ITTI_MSG_ID(received_message_p)) {
      case S5_CREATE_BEARER_REQUEST: {
        pgw_handle_create_bearer_request(
          &received_message_p->ittiMsg.s5_create_bearer_request);
      } break;

      case S5_NW_INITIATED_ACTIVATE_BEARER_RESP: {
        pgw_handle_nw_init_activate_bearer_rsp(
          &received_message_p->ittiMsg.s5_nw_init_actv_bearer_response);
      } break;

      case S5_NW_INITIATED_DEACTIVATE_BEARER_RESP: {
        pgw_handle_nw_init_deactivate_bearer_rsp(
          &received_message_p->ittiMsg.s5_nw_init_deactv_bearer_response);
      } break;

      case TERMINATE_MESSAGE: {
        pgw_exit();
        itti_exit_task();
      } break;

      default: {
        OAILOG_DEBUG(
          LOG_PGW_APP,
          "Unkwnon message ID %d:%s\n",
          ITTI_MSG_ID(received_message_p),
          ITTI_MSG_NAME(received_message_p));
      } break;
    }
    itti_free(ITTI_MSG_ORIGIN_ID(received_message_p), received_message_p);
    received_message_p = NULL;
  }

  return NULL;
}

int pgw_init(spgw_config_t *spgw_config_pP)
{
  if (itti_create_task(TASK_PGW_APP, &pgw_intertask_interface, NULL) < 0) {
    perror("pthread_create");
    OAILOG_ALERT(LOG_PGW_APP, "Initializing PGW-APP task interface: ERROR\n");
    return RETURNerror;
  }

  FILE *fp = NULL;
  bstring filename = bformat("/tmp/pgw_%d.status", g_pid);
  fp = fopen(bdata(filename), "w+");
  bdestroy(filename);
  fprintf(fp, "STARTED\n");
  fflush(fp);
  fclose(fp);

  OAILOG_DEBUG(LOG_PGW_APP, "Initializing PGW-APP task interface: DONE\n");
  return RETURNok;
}

static void pgw_exit(void)
{
  return;
}
