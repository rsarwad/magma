// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
)

// ServiceTypeWritePolicyRule grants write permission to service type based on policy.
func ServiceTypeWritePolicyRule() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		return cudBasedRule(FromContext(ctx).InventoryPolicy.ServiceType, m)
	})
}

// ServiceWritePolicyRule grants write permission to service based on policy.
func ServiceWritePolicyRule() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		return allowOrSkip(FromContext(ctx).InventoryPolicy.Equipment.Update)
	})
}

// ServiceEndpointWritePolicyRule grants write permission to service endpoint based on policy.
func ServiceEndpointWritePolicyRule() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		return allowOrSkip(FromContext(ctx).InventoryPolicy.Equipment.Update)
	})
}

// ServiceEndpointDefinitionWritePolicyRule grants write permission to service endpoint definition based on policy.
func ServiceEndpointDefinitionWritePolicyRule() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		return allowOrSkip(FromContext(ctx).InventoryPolicy.ServiceType.Update)
	})
}
