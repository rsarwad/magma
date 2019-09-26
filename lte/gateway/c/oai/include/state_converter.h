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
 *------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

#include "3gpp_23.003.h"
#include "hashtable.h"
#include "log.h"

#ifdef __cplusplus
}
#endif

#include <google/protobuf/map.h>

#include "lte/gateway/c/oai/protos/common_types.pb.h"

namespace magma {
namespace lte {

#define PLMN_BYTES 6

/**
 * StateConverter is a base class for state conversion between tasks state
 * structs and protobuf objects. This class is used to support specific state
 * conversion for each task that extends from it. The class doesn't hold memory,
 * all memory is owned by caller.
 */
class StateConverter {
 protected:
  StateConverter();
  ~StateConverter();

  static void guti_to_proto(const guti_t &guti_state, Guti *guti_proto);

  static void ecgi_to_proto(const ecgi_t &state_ecgi, Ecgi *ecgi_proto);

  /**
   * Function that converts hashtable struct to protobuf Map instance, using
   * a conversion function to convert each node of the hashtable, memory
   * is owned by the caller.
   * @tparam NodeType struct type of hashmap node entry
   * @tparam ProtoMessage protobuf type for proto map value entry
   * @param state_ht hashtable_ts_t struct to convert from
   * @param proto_map protobuf Map instance to convert to
   * @param conversion_callable conversion function for each entry of hashtable
   * @param log_task_level log level for task (LOG_MME_APP, LOG_SPGW_APP)
   */
  template<typename NodeType, typename ProtoMessage>
  static void hashtable_ts_to_proto(
    hash_table_ts_t *state_ht,
    google::protobuf::Map<unsigned int, ProtoMessage> *proto_map,
    std::function<void(NodeType *, ProtoMessage *)> conversion_callable,
    log_proto_t log_task_level)
  {
    hashtable_key_array_t *ht_keys = hashtable_ts_get_keys(state_ht);
    hashtable_rc_t ht_rc;
    if (ht_keys == nullptr) {
      return;
    }

    for (uint32_t i = 0; i < ht_keys->num_keys; i++) {
      NodeType *node;
      ht_rc = hashtable_ts_get(
        state_ht, (hash_key_t) ht_keys->keys[i], (void **) &node);
      if (ht_rc == HASH_TABLE_OK) {
        ProtoMessage proto;
        conversion_callable((NodeType *) node, &proto);
        (*proto_map)[ht_keys->keys[i]] = proto;
      } else {
        OAILOG_ERROR(
          log_task_level,
          "Key %u not found on %s hashtable",
          ht_keys->keys[i],
          state_ht->name->data);
      }
    }
  }

 private:
  static void plmn_to_chars(const plmn_t &state_plmn, char *plmn_array);
};

} // namespace lte
} // namespace magma
