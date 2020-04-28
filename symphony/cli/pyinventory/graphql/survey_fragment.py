#!/usr/bin/env python3
# @generated AUTOGENERATED file. Do not Change!

from dataclasses import dataclass
from datetime import datetime
from gql.gql.datetime_utils import DATETIME_FIELD
from gql.gql.graphql_client import GraphqlClient
from gql.gql.client import OperationException
from gql.gql.reporter import FailedOperationException
from functools import partial
from numbers import Number
from typing import Any, Callable, List, Mapping, Optional
from time import perf_counter
from dataclasses_json import DataClassJsonMixin

from .survey_question_fragment import SurveyQuestionFragment, QUERY as SurveyQuestionFragmentQuery
QUERY: List[str] = SurveyQuestionFragmentQuery + ["""
fragment SurveyFragment on Survey {
  id
  name
  completionTimestamp
  sourceFile {
    id
    fileName
    storeKey
  }
  surveyResponses {
    ...SurveyQuestionFragment
  }
}

"""]

@dataclass
class SurveyFragment(DataClassJsonMixin):
    @dataclass
    class File(DataClassJsonMixin):
        id: str
        fileName: str
        storeKey: Optional[str]

    @dataclass
    class SurveyQuestion(SurveyQuestionFragment):
        pass

    id: str
    name: str
    completionTimestamp: int
    surveyResponses: List[SurveyQuestion]
    sourceFile: Optional[File]
