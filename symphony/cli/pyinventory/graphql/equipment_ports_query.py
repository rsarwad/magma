#!/usr/bin/env python3
# @generated AUTOGENERATED file. Do not Change!

from dataclasses import dataclass
from datetime import datetime
from gql.gql.datetime_utils import DATETIME_FIELD
from gql.gql.graphql_client import GraphqlClient
from functools import partial
from numbers import Number
from typing import Any, Callable, List, Mapping, Optional

from dataclasses_json import DataClassJsonMixin

from .property_fragment import PropertyFragment, QUERY as PropertyFragmentQuery

@dataclass
class EquipmentPortsQuery(DataClassJsonMixin):
    @dataclass
    class EquipmentPortsQueryData(DataClassJsonMixin):
        @dataclass
        class Node(DataClassJsonMixin):
            @dataclass
            class EquipmentPort(DataClassJsonMixin):
                @dataclass
                class Property(PropertyFragment):
                    pass

                @dataclass
                class EquipmentPortDefinition(DataClassJsonMixin):
                    @dataclass
                    class EquipmentPortType(DataClassJsonMixin):
                        id: str
                        name: str

                    id: str
                    name: str
                    portType: Optional[EquipmentPortType] = None

                @dataclass
                class Link(DataClassJsonMixin):
                    @dataclass
                    class Service(DataClassJsonMixin):
                        id: str

                    id: str
                    services: List[Service]

                id: str
                properties: List[Property]
                definition: EquipmentPortDefinition
                link: Optional[Link] = None

            ports: List[EquipmentPort]

        equipment: Optional[Node] = None

    data: EquipmentPortsQueryData

    __QUERY__: str = PropertyFragmentQuery + """
    query EquipmentPortsQuery($id: ID!) {
  equipment: node(id: $id) {
    ... on Equipment {
      ports {
        id
        properties {
          ...PropertyFragment
        }
        definition {
          id
          name
          portType {
            id
            name
          }
        }
        link {
          id
          services {
            id
          }
        }
      }
    }
  }
}

    """

    @classmethod
    # fmt: off
    def execute(cls, client: GraphqlClient, id: str) -> EquipmentPortsQueryData:
        # fmt: off
        variables = {"id": id}
        response_text = client.call(cls.__QUERY__, variables=variables)
        return cls.from_json(response_text).data
