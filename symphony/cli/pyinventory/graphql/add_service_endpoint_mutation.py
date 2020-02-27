#!/usr/bin/env python3
# @generated AUTOGENERATED file. Do not Change!

from dataclasses import dataclass, field
from datetime import datetime
from gql.gql.datetime_utils import fromisoformat
from gql.gql.graphql_client import GraphqlClient
from functools import partial
from numbers import Number
from typing import Any, Callable, List, Mapping, Optional

from dataclasses_json import dataclass_json
from marshmallow import fields as marshmallow_fields

from .service_endpoint_role_enum import ServiceEndpointRole

from .add_service_endpoint_input import AddServiceEndpointInput


DATETIME_FIELD = field(
    metadata={
        "dataclasses_json": {
            "encoder": datetime.isoformat,
            "decoder": fromisoformat,
            "mm_field": marshmallow_fields.DateTime(format="iso"),
        }
    }
)


def enum_field(enum_type):
    def encode_enum(value):
        return value.value

    def decode_enum(t, value):
        return t(value)

    return field(
        metadata={
            "dataclasses_json": {
                "encoder": encode_enum,
                "decoder": partial(decode_enum, enum_type),
            }
        }
    )



@dataclass_json
@dataclass
class AddServiceEndpointMutation:
    __QUERY__ = """
    mutation AddServiceEndpointMutation($input: AddServiceEndpointInput!) {
  addServiceEndpoint(input: $input) {
    id
    name
    externalId
    customer {
      id
      name
      externalId
    }
    endpoints {
      id
      port {
        id
      }
      role
    }
    links {
      id
    }
  }
}

    """

    @dataclass_json
    @dataclass
    class AddServiceEndpointMutationData:
        @dataclass_json
        @dataclass
        class Service:
            @dataclass_json
            @dataclass
            class Customer:
                id: str
                name: str
                externalId: Optional[str] = None

            @dataclass_json
            @dataclass
            class ServiceEndpoint:
                @dataclass_json
                @dataclass
                class EquipmentPort:
                    id: str

                id: str
                port: EquipmentPort
                role: ServiceEndpointRole = enum_field(ServiceEndpointRole)

            @dataclass_json
            @dataclass
            class Link:
                id: str

            id: str
            name: str
            endpoints: List[ServiceEndpoint]
            links: List[Link]
            externalId: Optional[str] = None
            customer: Optional[Customer] = None

        addServiceEndpoint: Optional[Service] = None

    data: Optional[AddServiceEndpointMutationData] = None
    errors: Optional[Any] = None

    @classmethod
    # fmt: off
    def execute(cls, client: GraphqlClient, input: AddServiceEndpointInput):
        # fmt: off
        variables = {"input": input}
        response_text = client.call(cls.__QUERY__, variables=variables)
        return cls.from_json(response_text).data
