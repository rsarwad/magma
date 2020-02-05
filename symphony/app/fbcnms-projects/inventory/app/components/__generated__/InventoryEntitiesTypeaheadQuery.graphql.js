/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 4a939b24aae4e4cf85a1f10952b74378
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type InventoryEntitiesTypeaheadQueryVariables = {|
  name: string
|};
export type InventoryEntitiesTypeaheadQueryResponse = {|
  +searchForEntity: {|
    +edges: ?$ReadOnlyArray<{|
      +node: ?{|
        +entityId: string,
        +entityType: string,
        +name: string,
        +type: string,
      |}
    |}>
  |}
|};
export type InventoryEntitiesTypeaheadQuery = {|
  variables: InventoryEntitiesTypeaheadQueryVariables,
  response: InventoryEntitiesTypeaheadQueryResponse,
|};
*/


/*
query InventoryEntitiesTypeaheadQuery(
  $name: String!
) {
  searchForEntity(name: $name, first: 10) {
    edges {
      node {
        entityId
        entityType
        name
        type
      }
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "name",
    "type": "String!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "searchForEntity",
    "storageKey": null,
    "args": [
      {
        "kind": "Literal",
        "name": "first",
        "value": 10
      },
      {
        "kind": "Variable",
        "name": "name",
        "variableName": "name"
      }
    ],
    "concreteType": "SearchEntriesConnection",
    "plural": false,
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "edges",
        "storageKey": null,
        "args": null,
        "concreteType": "SearchEntryEdge",
        "plural": true,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "node",
            "storageKey": null,
            "args": null,
            "concreteType": "SearchEntry",
            "plural": false,
            "selections": [
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "entityId",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "entityType",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "name",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "type",
                "args": null,
                "storageKey": null
              }
            ]
          }
        ]
      }
    ]
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "InventoryEntitiesTypeaheadQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "InventoryEntitiesTypeaheadQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "InventoryEntitiesTypeaheadQuery",
    "id": null,
    "text": "query InventoryEntitiesTypeaheadQuery(\n  $name: String!\n) {\n  searchForEntity(name: $name, first: 10) {\n    edges {\n      node {\n        entityId\n        entityType\n        name\n        type\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '213cfdf6a9220dfc44c9c26e93a3acf0';
module.exports = node;
