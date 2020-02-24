#!/usr/bin/env python3
# @generated AUTOGENERATED file. Do not Change!

from dataclasses import dataclass, field
from datetime import datetime
from functools import partial
from numbers import Number
from typing import Any, Callable, List, Mapping, Optional

from dataclasses_json import dataclass_json
from marshmallow import fields as marshmallow_fields

from .datetime_utils import fromisoformat

from .survey_create_data_input import SurveyCreateData


DATETIME_FIELD = field(
    metadata={
        "dataclasses_json": {
            "encoder": datetime.isoformat,
            "decoder": fromisoformat,
            "mm_field": marshmallow_fields.DateTime(format="iso"),
        }
    }
)


@dataclass_json
@dataclass
class CreateSurveyMutation:
    __QUERY__ = """
    mutation CreateSurveyMutation($data: SurveyCreateData!) {
  createSurvey(data: $data)
}

    """

    @dataclass_json
    @dataclass
    class CreateSurveyMutationData:
        createSurvey: Optional[str] = None

    data: Optional[CreateSurveyMutationData] = None
    errors: Optional[Any] = None

    @classmethod
    # fmt: off
    def execute(cls, client, data: SurveyCreateData):
        # fmt: off
        variables = {"data": data}
        response_text = client.call(cls.__QUERY__, variables=variables)
        return cls.from_json(response_text).data
