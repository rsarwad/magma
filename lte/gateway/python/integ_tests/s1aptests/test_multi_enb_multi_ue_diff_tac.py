"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import time
import unittest
import s1ap_types
import s1ap_wrapper


class TestMultiEnbWithDifferentTac(unittest.TestCase):

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_detach_multienb_different_tac(self):
        """ Multi Enb attach with different TAC values """
        num_of_enbs = 5
        # column is a enb parameter,  row is a number of enbs
        """            Cell Id,   Tac, EnbType, PLMN Id """
        enb_list = list([[1,       1,     1,    "001010"],
                         [2,       2,     1,    "001010"],
                         [3,       3,     1,    "001010"],
                         [4,       4,     1,    "001010"],
                         [5,       5,     1,    "001010"]])

        assert (num_of_enbs == len(enb_list)), "Number of enbs configured"
        "not equal to enbs in the list!!!"

        self._s1ap_wrapper.multiEnbConfig(num_of_enbs, enb_list)

        time.sleep(2)

        ue_ids = []
        # UEs will attach to the ENBs in a round-robin fashion
        # each ENBs will be connected with 32UEs
        num_ues = 5
        self._s1ap_wrapper.configUEDevice(num_ues)
        for _ in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            print("******************** Calling attach for UE id ",
                  req.ue_id)
            self._s1ap_wrapper.s1_util.attach(
                req.ue_id,
                s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t)
            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()
            ue_ids.append(req.ue_id)

        for ue in ue_ids:
            print("************************* Calling detach for UE id ", ue)
            self._s1ap_wrapper.s1_util.detach(
                ue,
                s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value)


if __name__ == "__main__":
    unittest.main()
