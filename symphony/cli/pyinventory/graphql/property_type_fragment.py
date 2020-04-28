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

from gql.gql.enum_utils import enum_field
from .property_kind_enum import PropertyKind

QUERY: List[str] = ["""
fragment PropertyTypeFragment on PropertyType {
  id
  externalId
  name
  type
  index
  category
  stringValue
  intValue
  booleanValue
  floatValue
  latitudeValue
  longitudeValue
  rangeFromValue
  rangeToValue
  isEditable
  isInstanceProperty
  isMandatory
  isDeleted
}

"""]

@dataclass
class PropertyTypeFragment(DataClassJsonMixin):
    id: str
    name: str
    type: PropertyKind = enum_field(PropertyKind)
    externalId: Optional[str]
    index: Optional[int]
    category: Optional[str]
    stringValue: Optional[str]
    intValue: Optional[int]
    booleanValue: Optional[bool]
    floatValue: Optional[Number]
    latitudeValue: Optional[Number]
    longitudeValue: Optional[Number]
    rangeFromValue: Optional[Number]
    rangeToValue: Optional[Number]
    isEditable: Optional[bool]
    isInstanceProperty: Optional[bool]
    isMandatory: Optional[bool]
    isDeleted: Optional[bool]
