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

from ..fragment.equipment_port_definition import EquipmentPortDefinitionFragment, QUERY as EquipmentPortDefinitionFragmentQuery
from ..fragment.link import LinkFragment, QUERY as LinkFragmentQuery
from ..fragment.property import PropertyFragment, QUERY as PropertyFragmentQuery
QUERY: List[str] = EquipmentPortDefinitionFragmentQuery + LinkFragmentQuery + PropertyFragmentQuery + ["""
fragment EquipmentPortFragment on EquipmentPort {
  id
  properties {
    ...PropertyFragment
  }
  definition {
    ...EquipmentPortDefinitionFragment
  }
  link {
    ...LinkFragment
  }
}

"""]

@dataclass
class EquipmentPortFragment(DataClassJsonMixin):
    @dataclass
    class Property(PropertyFragment):
        pass

    @dataclass
    class EquipmentPortDefinition(EquipmentPortDefinitionFragment):
        pass

    @dataclass
    class Link(LinkFragment):
        pass

    id: str
    properties: List[Property]
    definition: EquipmentPortDefinition
    link: Optional[Link]
