/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 0e7cd7b39df56f179e555197cae29a33
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type WorkOrderDetails_workOrder$ref = any;
export type WorkOrderCardQueryVariables = {|
  workOrderId: string
|};
export type WorkOrderCardQueryResponse = {|
  +workOrder: ?{|
    +id?: string,
    +name?: string,
    +$fragmentRefs: WorkOrderDetails_workOrder$ref,
  |}
|};
export type WorkOrderCardQuery = {|
  variables: WorkOrderCardQueryVariables,
  response: WorkOrderCardQueryResponse,
|};
*/


/*
query WorkOrderCardQuery(
  $workOrderId: ID!
) {
  workOrder: node(id: $workOrderId) {
    __typename
    ... on WorkOrder {
      id
      name
      ...WorkOrderDetails_workOrder
    }
    id
  }
}

fragment CommentsBox_comments on Comment {
  ...CommentsLog_comments
}

fragment CommentsLog_comments on Comment {
  id
  ...TextCommentPost_comment
}

fragment DocumentTable_files on File {
  id
  fileName
  category
  ...FileAttachment_file
}

fragment DocumentTable_hyperlinks on Hyperlink {
  id
  category
  url
  displayName
  ...HyperlinkTableRow_hyperlink
}

fragment EntityDocumentsTable_files on File {
  ...DocumentTable_files
}

fragment EntityDocumentsTable_hyperlinks on Hyperlink {
  ...DocumentTable_hyperlinks
}

fragment EquipmentBreadcrumbs_equipment on Equipment {
  id
  name
  equipmentType {
    id
    name
  }
  locationHierarchy {
    id
    name
    locationType {
      name
      id
    }
  }
  positionHierarchy {
    id
    definition {
      id
      name
      visibleLabel
    }
    parentEquipment {
      id
      name
      equipmentType {
        id
        name
      }
    }
  }
}

fragment FileAttachment_file on File {
  id
  fileName
  sizeInBytes
  uploaded
  fileType
  storeKey
  category
  ...ImageDialog_img
}

fragment HyperlinkTableRow_hyperlink on Hyperlink {
  id
  category
  url
  displayName
  createTime
}

fragment ImageDialog_img on File {
  storeKey
  fileName
}

fragment LocationBreadcrumbsTitle_locationDetails on Location {
  id
  name
  locationType {
    name
    id
  }
  locationHierarchy {
    id
    name
    locationType {
      name
      id
    }
  }
}

fragment TextCommentPost_comment on Comment {
  id
  authorName
  text
  createTime
}

fragment WorkOrderDetailsPaneEquipmentItem_equipment on Equipment {
  id
  name
  equipmentType {
    id
    name
  }
  parentLocation {
    id
    name
    locationType {
      id
      name
    }
  }
  parentPosition {
    id
    definition {
      name
      visibleLabel
      id
    }
    parentEquipment {
      id
      name
    }
  }
}

fragment WorkOrderDetailsPaneLinkItem_link on Link {
  id
  futureState
  ports {
    id
    definition {
      id
      name
      visibleLabel
      portType {
        linkPropertyTypes {
          id
          name
          type
          index
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
        }
        id
      }
    }
    parentEquipment {
      id
      name
      futureState
      equipmentType {
        id
        name
        portDefinitions {
          id
          name
          visibleLabel
          bandwidth
          portType {
            id
            name
          }
        }
      }
      ...EquipmentBreadcrumbs_equipment
    }
    serviceEndpoints {
      role
      service {
        name
        id
      }
      id
    }
  }
  workOrder {
    id
    status
  }
  properties {
    id
    propertyType {
      id
      name
      type
      isEditable
      isMandatory
      isInstanceProperty
      stringValue
    }
    stringValue
    intValue
    floatValue
    booleanValue
    latitudeValue
    longitudeValue
    rangeFromValue
    rangeToValue
    equipmentValue {
      id
      name
    }
    locationValue {
      id
      name
    }
    serviceValue {
      id
      name
    }
  }
  services {
    id
    name
  }
}

fragment WorkOrderDetailsPane_workOrder on WorkOrder {
  id
  name
  equipmentToAdd {
    id
    ...WorkOrderDetailsPaneEquipmentItem_equipment
  }
  equipmentToRemove {
    id
    ...WorkOrderDetailsPaneEquipmentItem_equipment
  }
  linksToAdd {
    id
    ...WorkOrderDetailsPaneLinkItem_link
  }
  linksToRemove {
    id
    ...WorkOrderDetailsPaneLinkItem_link
  }
}

fragment WorkOrderDetails_workOrder on WorkOrder {
  id
  name
  description
  workOrderType {
    name
    id
  }
  location {
    name
    id
    latitude
    longitude
    locationType {
      mapType
      mapZoomLevel
      id
    }
    ...LocationBreadcrumbsTitle_locationDetails
  }
  ownerName
  assignee
  creationDate
  installDate
  status
  priority
  ...WorkOrderDetailsPane_workOrder
  properties {
    id
    propertyType {
      id
      name
      type
      isEditable
      isMandatory
      isInstanceProperty
      stringValue
    }
    stringValue
    intValue
    floatValue
    booleanValue
    latitudeValue
    longitudeValue
    rangeFromValue
    rangeToValue
    equipmentValue {
      id
      name
    }
    locationValue {
      id
      name
    }
    serviceValue {
      id
      name
    }
  }
  images {
    ...EntityDocumentsTable_files
    id
  }
  files {
    ...EntityDocumentsTable_files
    id
  }
  hyperlinks {
    ...EntityDocumentsTable_hyperlinks
    id
  }
  comments {
    ...CommentsBox_comments
    id
  }
  project {
    name
    id
    type {
      id
      name
    }
  }
  checkListCategories {
    id
    title
    description
    checkList {
      id
      index
      type
      title
      helpText
      checked
      enumValues
      stringValue
      enumSelectionMode
      selectedEnumValues
      yesNoResponse
      files {
        id
        fileName
        sizeInBytes
        modified
        uploaded
        fileType
        storeKey
        category
      }
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "workOrderId",
    "type": "ID!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "id",
    "variableName": "workOrderId"
  }
],
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
},
v4 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "description",
  "args": null,
  "storageKey": null
},
v5 = [
  (v3/*: any*/),
  (v2/*: any*/)
],
v6 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "locationHierarchy",
  "storageKey": null,
  "args": null,
  "concreteType": "Location",
  "plural": true,
  "selections": [
    (v2/*: any*/),
    (v3/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "locationType",
      "storageKey": null,
      "args": null,
      "concreteType": "LocationType",
      "plural": false,
      "selections": (v5/*: any*/)
    }
  ]
},
v7 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "status",
  "args": null,
  "storageKey": null
},
v8 = [
  (v2/*: any*/),
  (v3/*: any*/)
],
v9 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "equipmentType",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentType",
  "plural": false,
  "selections": (v8/*: any*/)
},
v10 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "visibleLabel",
  "args": null,
  "storageKey": null
},
v11 = [
  (v2/*: any*/),
  (v3/*: any*/),
  (v9/*: any*/),
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "parentLocation",
    "storageKey": null,
    "args": null,
    "concreteType": "Location",
    "plural": false,
    "selections": [
      (v2/*: any*/),
      (v3/*: any*/),
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "locationType",
        "storageKey": null,
        "args": null,
        "concreteType": "LocationType",
        "plural": false,
        "selections": (v8/*: any*/)
      }
    ]
  },
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "parentPosition",
    "storageKey": null,
    "args": null,
    "concreteType": "EquipmentPosition",
    "plural": false,
    "selections": [
      (v2/*: any*/),
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "definition",
        "storageKey": null,
        "args": null,
        "concreteType": "EquipmentPositionDefinition",
        "plural": false,
        "selections": [
          (v3/*: any*/),
          (v10/*: any*/),
          (v2/*: any*/)
        ]
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "parentEquipment",
        "storageKey": null,
        "args": null,
        "concreteType": "Equipment",
        "plural": false,
        "selections": (v8/*: any*/)
      }
    ]
  }
],
v12 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "futureState",
  "args": null,
  "storageKey": null
},
v13 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "type",
  "args": null,
  "storageKey": null
},
v14 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "index",
  "args": null,
  "storageKey": null
},
v15 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "stringValue",
  "args": null,
  "storageKey": null
},
v16 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "intValue",
  "args": null,
  "storageKey": null
},
v17 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "booleanValue",
  "args": null,
  "storageKey": null
},
v18 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "floatValue",
  "args": null,
  "storageKey": null
},
v19 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "latitudeValue",
  "args": null,
  "storageKey": null
},
v20 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "longitudeValue",
  "args": null,
  "storageKey": null
},
v21 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeFromValue",
  "args": null,
  "storageKey": null
},
v22 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeToValue",
  "args": null,
  "storageKey": null
},
v23 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isEditable",
  "args": null,
  "storageKey": null
},
v24 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isInstanceProperty",
  "args": null,
  "storageKey": null
},
v25 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isMandatory",
  "args": null,
  "storageKey": null
},
v26 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "properties",
  "storageKey": null,
  "args": null,
  "concreteType": "Property",
  "plural": true,
  "selections": [
    (v2/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "propertyType",
      "storageKey": null,
      "args": null,
      "concreteType": "PropertyType",
      "plural": false,
      "selections": [
        (v2/*: any*/),
        (v3/*: any*/),
        (v13/*: any*/),
        (v23/*: any*/),
        (v25/*: any*/),
        (v24/*: any*/),
        (v15/*: any*/)
      ]
    },
    (v15/*: any*/),
    (v16/*: any*/),
    (v18/*: any*/),
    (v17/*: any*/),
    (v19/*: any*/),
    (v20/*: any*/),
    (v21/*: any*/),
    (v22/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "equipmentValue",
      "storageKey": null,
      "args": null,
      "concreteType": "Equipment",
      "plural": false,
      "selections": (v8/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "locationValue",
      "storageKey": null,
      "args": null,
      "concreteType": "Location",
      "plural": false,
      "selections": (v8/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "serviceValue",
      "storageKey": null,
      "args": null,
      "concreteType": "Service",
      "plural": false,
      "selections": (v8/*: any*/)
    }
  ]
},
v27 = [
  (v2/*: any*/),
  (v12/*: any*/),
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "ports",
    "storageKey": null,
    "args": null,
    "concreteType": "EquipmentPort",
    "plural": true,
    "selections": [
      (v2/*: any*/),
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "definition",
        "storageKey": null,
        "args": null,
        "concreteType": "EquipmentPortDefinition",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          (v10/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "portType",
            "storageKey": null,
            "args": null,
            "concreteType": "EquipmentPortType",
            "plural": false,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "linkPropertyTypes",
                "storageKey": null,
                "args": null,
                "concreteType": "PropertyType",
                "plural": true,
                "selections": [
                  (v2/*: any*/),
                  (v3/*: any*/),
                  (v13/*: any*/),
                  (v14/*: any*/),
                  (v15/*: any*/),
                  (v16/*: any*/),
                  (v17/*: any*/),
                  (v18/*: any*/),
                  (v19/*: any*/),
                  (v20/*: any*/),
                  (v21/*: any*/),
                  (v22/*: any*/),
                  (v23/*: any*/),
                  (v24/*: any*/),
                  (v25/*: any*/)
                ]
              },
              (v2/*: any*/)
            ]
          }
        ]
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "parentEquipment",
        "storageKey": null,
        "args": null,
        "concreteType": "Equipment",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          (v12/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "equipmentType",
            "storageKey": null,
            "args": null,
            "concreteType": "EquipmentType",
            "plural": false,
            "selections": [
              (v2/*: any*/),
              (v3/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "portDefinitions",
                "storageKey": null,
                "args": null,
                "concreteType": "EquipmentPortDefinition",
                "plural": true,
                "selections": [
                  (v2/*: any*/),
                  (v3/*: any*/),
                  (v10/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "bandwidth",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "portType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "EquipmentPortType",
                    "plural": false,
                    "selections": (v8/*: any*/)
                  }
                ]
              }
            ]
          },
          (v6/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "positionHierarchy",
            "storageKey": null,
            "args": null,
            "concreteType": "EquipmentPosition",
            "plural": true,
            "selections": [
              (v2/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "definition",
                "storageKey": null,
                "args": null,
                "concreteType": "EquipmentPositionDefinition",
                "plural": false,
                "selections": [
                  (v2/*: any*/),
                  (v3/*: any*/),
                  (v10/*: any*/)
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "parentEquipment",
                "storageKey": null,
                "args": null,
                "concreteType": "Equipment",
                "plural": false,
                "selections": [
                  (v2/*: any*/),
                  (v3/*: any*/),
                  (v9/*: any*/)
                ]
              }
            ]
          }
        ]
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "serviceEndpoints",
        "storageKey": null,
        "args": null,
        "concreteType": "ServiceEndpoint",
        "plural": true,
        "selections": [
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "role",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "service",
            "storageKey": null,
            "args": null,
            "concreteType": "Service",
            "plural": false,
            "selections": (v5/*: any*/)
          },
          (v2/*: any*/)
        ]
      }
    ]
  },
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "workOrder",
    "storageKey": null,
    "args": null,
    "concreteType": "WorkOrder",
    "plural": false,
    "selections": [
      (v2/*: any*/),
      (v7/*: any*/)
    ]
  },
  (v26/*: any*/),
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "services",
    "storageKey": null,
    "args": null,
    "concreteType": "Service",
    "plural": true,
    "selections": (v8/*: any*/)
  }
],
v28 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "fileName",
  "args": null,
  "storageKey": null
},
v29 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "category",
  "args": null,
  "storageKey": null
},
v30 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "sizeInBytes",
  "args": null,
  "storageKey": null
},
v31 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "uploaded",
  "args": null,
  "storageKey": null
},
v32 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "fileType",
  "args": null,
  "storageKey": null
},
v33 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "storeKey",
  "args": null,
  "storageKey": null
},
v34 = [
  (v2/*: any*/),
  (v28/*: any*/),
  (v29/*: any*/),
  (v30/*: any*/),
  (v31/*: any*/),
  (v32/*: any*/),
  (v33/*: any*/)
],
v35 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "createTime",
  "args": null,
  "storageKey": null
},
v36 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "title",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "WorkOrderCardQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "workOrder",
        "name": "node",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": null,
        "plural": false,
        "selections": [
          {
            "kind": "InlineFragment",
            "type": "WorkOrder",
            "selections": [
              (v2/*: any*/),
              (v3/*: any*/),
              {
                "kind": "FragmentSpread",
                "name": "WorkOrderDetails_workOrder",
                "args": null
              }
            ]
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "WorkOrderCardQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "workOrder",
        "name": "node",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": null,
        "plural": false,
        "selections": [
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "__typename",
            "args": null,
            "storageKey": null
          },
          (v2/*: any*/),
          {
            "kind": "InlineFragment",
            "type": "WorkOrder",
            "selections": [
              (v3/*: any*/),
              (v4/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "workOrderType",
                "storageKey": null,
                "args": null,
                "concreteType": "WorkOrderType",
                "plural": false,
                "selections": (v5/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "location",
                "storageKey": null,
                "args": null,
                "concreteType": "Location",
                "plural": false,
                "selections": [
                  (v3/*: any*/),
                  (v2/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "latitude",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "longitude",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "locationType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "LocationType",
                    "plural": false,
                    "selections": [
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "mapType",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "mapZoomLevel",
                        "args": null,
                        "storageKey": null
                      },
                      (v2/*: any*/),
                      (v3/*: any*/)
                    ]
                  },
                  (v6/*: any*/)
                ]
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "ownerName",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "assignee",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "creationDate",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "installDate",
                "args": null,
                "storageKey": null
              },
              (v7/*: any*/),
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "priority",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "equipmentToAdd",
                "storageKey": null,
                "args": null,
                "concreteType": "Equipment",
                "plural": true,
                "selections": (v11/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "equipmentToRemove",
                "storageKey": null,
                "args": null,
                "concreteType": "Equipment",
                "plural": true,
                "selections": (v11/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "linksToAdd",
                "storageKey": null,
                "args": null,
                "concreteType": "Link",
                "plural": true,
                "selections": (v27/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "linksToRemove",
                "storageKey": null,
                "args": null,
                "concreteType": "Link",
                "plural": true,
                "selections": (v27/*: any*/)
              },
              (v26/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "images",
                "storageKey": null,
                "args": null,
                "concreteType": "File",
                "plural": true,
                "selections": (v34/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "files",
                "storageKey": null,
                "args": null,
                "concreteType": "File",
                "plural": true,
                "selections": (v34/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "hyperlinks",
                "storageKey": null,
                "args": null,
                "concreteType": "Hyperlink",
                "plural": true,
                "selections": [
                  (v2/*: any*/),
                  (v29/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "url",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "displayName",
                    "args": null,
                    "storageKey": null
                  },
                  (v35/*: any*/)
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "comments",
                "storageKey": null,
                "args": null,
                "concreteType": "Comment",
                "plural": true,
                "selections": [
                  (v2/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "authorName",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "text",
                    "args": null,
                    "storageKey": null
                  },
                  (v35/*: any*/)
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "project",
                "storageKey": null,
                "args": null,
                "concreteType": "Project",
                "plural": false,
                "selections": [
                  (v3/*: any*/),
                  (v2/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "type",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "ProjectType",
                    "plural": false,
                    "selections": (v8/*: any*/)
                  }
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "checkListCategories",
                "storageKey": null,
                "args": null,
                "concreteType": "CheckListCategory",
                "plural": true,
                "selections": [
                  (v2/*: any*/),
                  (v36/*: any*/),
                  (v4/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "checkList",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "CheckListItem",
                    "plural": true,
                    "selections": [
                      (v2/*: any*/),
                      (v14/*: any*/),
                      (v13/*: any*/),
                      (v36/*: any*/),
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "helpText",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "checked",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "enumValues",
                        "args": null,
                        "storageKey": null
                      },
                      (v15/*: any*/),
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "enumSelectionMode",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "selectedEnumValues",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "yesNoResponse",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "files",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "File",
                        "plural": true,
                        "selections": [
                          (v2/*: any*/),
                          (v28/*: any*/),
                          (v30/*: any*/),
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "modified",
                            "args": null,
                            "storageKey": null
                          },
                          (v31/*: any*/),
                          (v32/*: any*/),
                          (v33/*: any*/),
                          (v29/*: any*/)
                        ]
                      }
                    ]
                  }
                ]
              }
            ]
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "WorkOrderCardQuery",
    "id": null,
    "text": "query WorkOrderCardQuery(\n  $workOrderId: ID!\n) {\n  workOrder: node(id: $workOrderId) {\n    __typename\n    ... on WorkOrder {\n      id\n      name\n      ...WorkOrderDetails_workOrder\n    }\n    id\n  }\n}\n\nfragment CommentsBox_comments on Comment {\n  ...CommentsLog_comments\n}\n\nfragment CommentsLog_comments on Comment {\n  id\n  ...TextCommentPost_comment\n}\n\nfragment DocumentTable_files on File {\n  id\n  fileName\n  category\n  ...FileAttachment_file\n}\n\nfragment DocumentTable_hyperlinks on Hyperlink {\n  id\n  category\n  url\n  displayName\n  ...HyperlinkTableRow_hyperlink\n}\n\nfragment EntityDocumentsTable_files on File {\n  ...DocumentTable_files\n}\n\nfragment EntityDocumentsTable_hyperlinks on Hyperlink {\n  ...DocumentTable_hyperlinks\n}\n\nfragment EquipmentBreadcrumbs_equipment on Equipment {\n  id\n  name\n  equipmentType {\n    id\n    name\n  }\n  locationHierarchy {\n    id\n    name\n    locationType {\n      name\n      id\n    }\n  }\n  positionHierarchy {\n    id\n    definition {\n      id\n      name\n      visibleLabel\n    }\n    parentEquipment {\n      id\n      name\n      equipmentType {\n        id\n        name\n      }\n    }\n  }\n}\n\nfragment FileAttachment_file on File {\n  id\n  fileName\n  sizeInBytes\n  uploaded\n  fileType\n  storeKey\n  category\n  ...ImageDialog_img\n}\n\nfragment HyperlinkTableRow_hyperlink on Hyperlink {\n  id\n  category\n  url\n  displayName\n  createTime\n}\n\nfragment ImageDialog_img on File {\n  storeKey\n  fileName\n}\n\nfragment LocationBreadcrumbsTitle_locationDetails on Location {\n  id\n  name\n  locationType {\n    name\n    id\n  }\n  locationHierarchy {\n    id\n    name\n    locationType {\n      name\n      id\n    }\n  }\n}\n\nfragment TextCommentPost_comment on Comment {\n  id\n  authorName\n  text\n  createTime\n}\n\nfragment WorkOrderDetailsPaneEquipmentItem_equipment on Equipment {\n  id\n  name\n  equipmentType {\n    id\n    name\n  }\n  parentLocation {\n    id\n    name\n    locationType {\n      id\n      name\n    }\n  }\n  parentPosition {\n    id\n    definition {\n      name\n      visibleLabel\n      id\n    }\n    parentEquipment {\n      id\n      name\n    }\n  }\n}\n\nfragment WorkOrderDetailsPaneLinkItem_link on Link {\n  id\n  futureState\n  ports {\n    id\n    definition {\n      id\n      name\n      visibleLabel\n      portType {\n        linkPropertyTypes {\n          id\n          name\n          type\n          index\n          stringValue\n          intValue\n          booleanValue\n          floatValue\n          latitudeValue\n          longitudeValue\n          rangeFromValue\n          rangeToValue\n          isEditable\n          isInstanceProperty\n          isMandatory\n        }\n        id\n      }\n    }\n    parentEquipment {\n      id\n      name\n      futureState\n      equipmentType {\n        id\n        name\n        portDefinitions {\n          id\n          name\n          visibleLabel\n          bandwidth\n          portType {\n            id\n            name\n          }\n        }\n      }\n      ...EquipmentBreadcrumbs_equipment\n    }\n    serviceEndpoints {\n      role\n      service {\n        name\n        id\n      }\n      id\n    }\n  }\n  workOrder {\n    id\n    status\n  }\n  properties {\n    id\n    propertyType {\n      id\n      name\n      type\n      isEditable\n      isMandatory\n      isInstanceProperty\n      stringValue\n    }\n    stringValue\n    intValue\n    floatValue\n    booleanValue\n    latitudeValue\n    longitudeValue\n    rangeFromValue\n    rangeToValue\n    equipmentValue {\n      id\n      name\n    }\n    locationValue {\n      id\n      name\n    }\n    serviceValue {\n      id\n      name\n    }\n  }\n  services {\n    id\n    name\n  }\n}\n\nfragment WorkOrderDetailsPane_workOrder on WorkOrder {\n  id\n  name\n  equipmentToAdd {\n    id\n    ...WorkOrderDetailsPaneEquipmentItem_equipment\n  }\n  equipmentToRemove {\n    id\n    ...WorkOrderDetailsPaneEquipmentItem_equipment\n  }\n  linksToAdd {\n    id\n    ...WorkOrderDetailsPaneLinkItem_link\n  }\n  linksToRemove {\n    id\n    ...WorkOrderDetailsPaneLinkItem_link\n  }\n}\n\nfragment WorkOrderDetails_workOrder on WorkOrder {\n  id\n  name\n  description\n  workOrderType {\n    name\n    id\n  }\n  location {\n    name\n    id\n    latitude\n    longitude\n    locationType {\n      mapType\n      mapZoomLevel\n      id\n    }\n    ...LocationBreadcrumbsTitle_locationDetails\n  }\n  ownerName\n  assignee\n  creationDate\n  installDate\n  status\n  priority\n  ...WorkOrderDetailsPane_workOrder\n  properties {\n    id\n    propertyType {\n      id\n      name\n      type\n      isEditable\n      isMandatory\n      isInstanceProperty\n      stringValue\n    }\n    stringValue\n    intValue\n    floatValue\n    booleanValue\n    latitudeValue\n    longitudeValue\n    rangeFromValue\n    rangeToValue\n    equipmentValue {\n      id\n      name\n    }\n    locationValue {\n      id\n      name\n    }\n    serviceValue {\n      id\n      name\n    }\n  }\n  images {\n    ...EntityDocumentsTable_files\n    id\n  }\n  files {\n    ...EntityDocumentsTable_files\n    id\n  }\n  hyperlinks {\n    ...EntityDocumentsTable_hyperlinks\n    id\n  }\n  comments {\n    ...CommentsBox_comments\n    id\n  }\n  project {\n    name\n    id\n    type {\n      id\n      name\n    }\n  }\n  checkListCategories {\n    id\n    title\n    description\n    checkList {\n      id\n      index\n      type\n      title\n      helpText\n      checked\n      enumValues\n      stringValue\n      enumSelectionMode\n      selectedEnumValues\n      yesNoResponse\n      files {\n        id\n        fileName\n        sizeInBytes\n        modified\n        uploaded\n        fileType\n        storeKey\n        category\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'ba8f1936a3ba4fd0f7aa7b7b90970432';
module.exports = node;
