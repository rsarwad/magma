#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import Dict, List, Optional, Tuple

from pysymphony import SymphonyClient

from .._utils import PropertyValue, format_properties, get_graphql_property_inputs
from ..common.cache import SERVICE_TYPES
from ..common.data_class import (
    Customer,
    EquipmentPort,
    EquipmentPortDefinition,
    Link,
    Service,
    ServiceEndpoint,
    ServiceType,
)
from ..common.data_enum import Entity
from ..exceptions import EntityNotFoundError
from ..graphql.enum.service_status import ServiceStatus
from ..graphql.input.add_service_endpoint import AddServiceEndpointInput
from ..graphql.input.service_create_data import ServiceCreateData
from ..graphql.input.service_type_create_data import ServiceTypeCreateData
from ..graphql.mutation.add_service import AddServiceMutation
from ..graphql.mutation.add_service_endpoint import AddServiceEndpointMutation
from ..graphql.mutation.add_service_link import AddServiceLinkMutation
from ..graphql.mutation.add_service_type import AddServiceTypeMutation
from ..graphql.mutation.remove_service import RemoveServiceMutation
from ..graphql.mutation.remove_service_type import RemoveServiceTypeMutation
from ..graphql.query.service_details import ServiceDetailsQuery
from ..graphql.query.service_type_services import ServiceTypeServicesQuery
from ..graphql.query.service_types import ServiceTypesQuery


def _populate_service_types(client: SymphonyClient) -> None:
    service_types = ServiceTypesQuery.execute(client)
    if not service_types:
        return
    edges = service_types.edges
    for edge in edges:
        node = edge.node
        if node is not None:
            SERVICE_TYPES[node.name] = ServiceType(
                name=node.name,
                id=node.id,
                hasCustomer=node.hasCustomer,
                property_types=node.propertyTypes,
            )


def add_service_type(
    client: SymphonyClient,
    name: str,
    hasCustomer: bool,
    properties: List[Tuple[str, str, Optional[PropertyValue], Optional[bool]]],
) -> ServiceType:

    new_property_types = format_properties(properties)
    result = AddServiceTypeMutation.execute(
        client,
        data=ServiceTypeCreateData(
            name=name, hasCustomer=hasCustomer, properties=new_property_types
        ),
    )

    service_type = ServiceType(
        name=result.name,
        id=result.id,
        hasCustomer=result.hasCustomer,
        property_types=result.propertyTypes,
    )
    SERVICE_TYPES[name] = service_type
    return service_type


def add_service(
    client: SymphonyClient,
    name: str,
    external_id: Optional[str],
    service_type: str,
    customer: Optional[Customer],
    properties_dict: Dict[str, PropertyValue],
) -> Service:
    property_types = SERVICE_TYPES[service_type].property_types
    properties = get_graphql_property_inputs(property_types, properties_dict)
    service_create_data = ServiceCreateData(
        name=name,
        externalId=external_id,
        serviceTypeId=SERVICE_TYPES[service_type].id,
        status=ServiceStatus.PENDING,
        customerId=customer.id if customer is not None else None,
        properties=properties,
        upstreamServiceIds=[],
    )
    result = AddServiceMutation.execute(client, data=service_create_data)
    returned_customer = result.customer
    endpoints = []
    for e in result.endpoints:
        port = e.port
        link = port.link if port else None
        endpoints.append(
            ServiceEndpoint(
                id=e.id,
                port=EquipmentPort(
                    id=port.id,
                    properties=port.properties,
                    definition=EquipmentPortDefinition(
                        id=port.definition.id, name=port.definition.name
                    ),
                    link=Link(
                        link.id,
                        properties=link.properties,
                        service_ids=[s.id for s in link.services],
                    )
                    if link
                    else None,
                )
                if port
                else None,
                # TODO add service_endpoint_type api
                type="1",
            )
        )
    return Service(
        name=result.name,
        id=result.id,
        externalId=result.externalId,
        customer=Customer(
            name=returned_customer.name,
            id=returned_customer.id,
            externalId=returned_customer.externalId,
        )
        if returned_customer
        else None,
        endpoints=endpoints,
        links=[
            Link(
                id=l.id, properties=l.properties, service_ids=[s.id for s in l.services]
            )
            for l in result.links
        ],
    )


def add_service_endpoint(
    client: SymphonyClient, service: Service, port: EquipmentPort
) -> None:
    AddServiceEndpointMutation.execute(
        client,
        input=AddServiceEndpointInput(
            id=service.id, portId=port.id, definition="1", equipmentID="1"
        ),
    )


def add_service_link(client: SymphonyClient, service: Service, link: Link) -> None:
    AddServiceLinkMutation.execute(client, id=service.id, linkId=link.id)


def get_service(client: SymphonyClient, id: str) -> Service:
    result = ServiceDetailsQuery.execute(client, id=id)
    if result is None:
        raise EntityNotFoundError(entity=Entity.Service, entity_id=id)
    customer = result.customer
    endpoints = []
    for e in result.endpoints:
        port = e.port
        link = port.link if port else None
        endpoints.append(
            ServiceEndpoint(
                id=e.id,
                port=EquipmentPort(
                    id=port.id,
                    properties=port.properties,
                    definition=EquipmentPortDefinition(
                        id=port.definition.id, name=port.definition.name
                    ),
                    link=Link(
                        id=link.id,
                        properties=link.properties,
                        service_ids=[s.id for s in link.services],
                    )
                    if link
                    else None,
                )
                if port
                else None,
                # TODO add service_endpoint_type api
                type="1",
            )
        )
    return Service(
        name=result.name,
        id=result.id,
        externalId=result.externalId,
        customer=Customer(
            name=customer.name, id=customer.id, externalId=customer.externalId
        )
        if customer is not None
        else None,
        endpoints=endpoints,
        links=[
            Link(
                id=l.id, properties=l.properties, service_ids=[s.id for s in l.services]
            )
            for l in result.links
        ],
    )


def delete_service_type_with_services(
    client: SymphonyClient, service_type: ServiceType
) -> None:
    service_type_with_services = ServiceTypeServicesQuery.execute(
        client, id=service_type.id
    )
    if not service_type_with_services:
        raise EntityNotFoundError(entity=Entity.ServiceType, entity_id=service_type.id)
    services = service_type_with_services.services
    for service in services:
        RemoveServiceMutation.execute(client, id=service.id)
    RemoveServiceTypeMutation.execute(client, id=service_type.id)
