// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"

	"github.com/facebookincubator/symphony/pkg/authz/models"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/privacy"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/viewer"
)

func cudBasedCheck(cud *models.Cud, m ent.Mutation) bool {
	var permission *models.BasicPermissionRule
	switch {
	case m.Op().Is(ent.OpCreate):
		permission = cud.Create
	case m.Op().Is(ent.OpUpdateOne | ent.OpUpdate):
		permission = cud.Update
	case m.Op().Is(ent.OpDeleteOne | ent.OpDelete):
		permission = cud.Delete
	default:
		return false
	}
	return permission.IsAllowed == models.PermissionValueYes
}

func allowOrSkip(r *models.BasicPermissionRule) error {
	if r.IsAllowed == models.PermissionValueYes {
		return privacy.Allow
	}
	return privacy.Skip
}

func allowOrSkipLocations(r *models.LocationPermissionRule, locationTypeID int) error {
	switch r.IsAllowed {
	case models.PermissionValueYes:
		return privacy.Allow
	case models.PermissionValueByCondition:
		for _, typeID := range r.LocationTypeIds {
			if typeID == locationTypeID {
				return privacy.Allow
			}
		}
	}
	return privacy.Skip
}

func privacyDecision(allowed bool) error {
	if allowed {
		return privacy.Allow
	}
	return privacy.Skip
}

func checkWorkforce(r *models.WorkforcePermissionRule, workOrderTypeID *int, projectTypeID *int) bool {
	switch r.IsAllowed {
	case models.PermissionValueYes:
		return true
	case models.PermissionValueByCondition:
		if workOrderTypeID != nil {
			for _, typeID := range r.WorkOrderTypeIds {
				if typeID == *workOrderTypeID {
					return true
				}
			}
		}
		if projectTypeID != nil {
			for _, typeID := range r.ProjectTypeIds {
				if typeID == *projectTypeID {
					return true
				}
			}
		}
	}
	return false
}

func cudBasedRule(cud *models.Cud, m ent.Mutation) error {
	return privacyDecision(cudBasedCheck(cud, m))
}

func allowWritePermissionsRule() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, _ ent.Mutation) error {
		v := viewer.FromContext(ctx)
		if v != nil && v.Features().Enabled(viewer.FeaturePermissionPolicies) {
			return privacyDecision(v.Role() == user.RoleOWNER)
		}
		return privacyDecision(FromContext(ctx).CanWrite)
	})
}

func allowReadPermissionsRule() privacy.QueryRule {
	return privacy.QueryRuleFunc(func(ctx context.Context, _ ent.Query) error {
		switch v := viewer.FromContext(ctx); {
		case v == nil:
			return privacy.Skip
		case !v.Features().Enabled(viewer.FeaturePermissionPolicies),
			v.Role() == user.RoleOWNER:
			return privacy.Allow
		default:
			return privacy.Skip
		}
	})
}

func denyIfNoPermissionSettingsRule() privacy.QueryMutationRule {
	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
		if FromContext(ctx) == nil {
			return privacy.Deny
		}
		return privacy.Skip
	})
}

func allowOrSkipWorkOrder(ctx context.Context, p *models.PermissionSettings, wo *ent.WorkOrder) error {
	switch allowed, err := isViewerWorkOrderOwnerOrAssignee(ctx, wo); {
	case err != nil:
		return privacy.Denyf("cannot check work order viewer relation: %w", err)
	case allowed:
		return privacy.Allow
	}
	workOrderTypeID, err := wo.QueryType().OnlyID(ctx)
	if err != nil {
		return privacy.Denyf("cannot fetch work order type id: %w", err)
	}
	return privacyDecision(
		checkWorkforce(
			p.WorkforcePolicy.Data.Update, &workOrderTypeID, nil,
		),
	)
}

func allowOrSkipProject(ctx context.Context, p *models.PermissionSettings, proj *ent.Project) error {
	if userViewer, ok := viewer.FromContext(ctx).(*viewer.UserViewer); ok {
		switch isCreator, err := isUserProjectCreator(ctx, userViewer.User().ID, proj); {
		case err != nil:
			return privacy.Denyf("cannot check project viewer relation: %w", err)
		case isCreator:
			return privacy.Allow
		}
	}

	projectTypeID, err := proj.QueryType().OnlyID(ctx)
	if err != nil {
		return privacy.Denyf("cannot fetch project type id: %w", err)
	}
	return privacyDecision(
		checkWorkforce(
			p.WorkforcePolicy.Data.Update, nil, &projectTypeID,
		),
	)
}

func denyBulkEditOrDeleteRule() privacy.MutationRule {
	rule := privacy.DenyMutationOperationRule(ent.OpUpdate | ent.OpDelete)
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		if v := viewer.FromContext(ctx); v == nil ||
			!v.Features().Enabled(viewer.FeaturePermissionPolicies) {
			return privacy.Skip
		}
		return rule.EvalMutation(ctx, m)
	})
}
