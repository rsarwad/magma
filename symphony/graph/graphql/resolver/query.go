// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/pkg/actions"
	"github.com/facebookincubator/symphony/pkg/actions/core"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type queryResolver struct{ resolver }

func (queryResolver) Me(ctx context.Context) (*viewer.Viewer, error) {
	return viewer.FromContext(ctx), nil
}

func (r queryResolver) Node(ctx context.Context, id string) (ent.Noder, error) {
	n, err := r.ClientFrom(ctx).Noder(ctx, id)
	if err == nil {
		return n, nil
	}
	r.log.For(ctx).
		Debug("cannot query node",
			zap.String("id", id),
			zap.Error(err),
		)
	var e *ent.ErrNotFound
	if xerrors.As(err, &e) {
		err = nil
	}
	return nil, err
}

func (r queryResolver) Location(ctx context.Context, id string) (*ent.Location, error) {
	l, err := r.ClientFrom(ctx).Location.Get(ctx, id)
	if err != nil && !ent.IsNotFound(err) {
		return nil, errors.Wrapf(err, "querying location: id=%q", id)
	}
	return l, nil
}

func (r queryResolver) LocationType(ctx context.Context, id string) (*ent.LocationType, error) {
	lt, err := r.ClientFrom(ctx).LocationType.Get(ctx, id)
	if err != nil && !ent.IsNotFound(err) {
		return nil, errors.Wrapf(err, "querying location type: id=%q", id)
	}
	return lt, nil
}

func (r queryResolver) LocationTypes(
	ctx context.Context,
	after *ent.Cursor, first *int,
	before *ent.Cursor, last *int,
) (*ent.LocationTypeConnection, error) {
	return r.ClientFrom(ctx).LocationType.Query().
		Paginate(ctx, after, first, before, last)
}

func (r queryResolver) Locations(
	ctx context.Context, onlyTopLevel *bool,
	types []string, name *string, needsSiteSurvey *bool,
	after *ent.Cursor, first *int,
	before *ent.Cursor, last *int,
) (*ent.LocationConnection, error) {
	query := r.ClientFrom(ctx).Location.Query()
	if onlyTopLevel != nil && *onlyTopLevel {
		query = query.Where(location.Not(location.HasParent()))
	}
	if name != nil {
		query = query.Where(location.NameContainsFold(*name))
	}
	if len(types) > 0 {
		query = query.Where(location.HasTypeWith(locationtype.IDIn(types...)))
	}
	if needsSiteSurvey != nil {
		query = query.Where(location.SiteSurveyNeeded(*needsSiteSurvey))
	}
	return query.Paginate(ctx, after, first, before, last)
}

func (r queryResolver) NearestSites(ctx context.Context, latitude, longitude float64, first int) ([]*ent.Location, error) {
	sites := r.ClientFrom(ctx).Location.Query().Where(location.HasTypeWith(locationtype.Site(true))).AllX(ctx)
	var lr locationResolver
	sort.Slice(sites, func(i, j int) bool {
		d1, _ := lr.DistanceKm(ctx, sites[i], latitude, longitude)
		d2, _ := lr.DistanceKm(ctx, sites[j], latitude, longitude)
		return d1 < d2
	})
	if len(sites) < first {
		return sites, nil
	}
	return sites[:first], nil
}

func (r queryResolver) Equipment(ctx context.Context, id string) (*ent.Equipment, error) {
	e, err := r.ClientFrom(ctx).Equipment.Get(ctx, id)
	if err != nil && !ent.IsNotFound(err) {
		return nil, errors.Wrapf(err, "querying equipment: id=%q", id)
	}
	return e, nil
}

func (r queryResolver) EquipmentType(ctx context.Context, id string) (*ent.EquipmentType, error) {
	et, err := r.ClientFrom(ctx).EquipmentType.Get(ctx, id)
	if err != nil && !ent.IsNotFound(err) {
		return nil, errors.Wrapf(err, "querying equipment type: id=%q", id)
	}
	return et, nil
}

func (r queryResolver) EquipmentTypes(
	ctx context.Context,
	after *ent.Cursor, first *int,
	before *ent.Cursor, last *int,
) (*ent.EquipmentTypeConnection, error) {
	return r.ClientFrom(ctx).EquipmentType.Query().
		Paginate(ctx, after, first, before, last)
}

func (r queryResolver) EquipmentPortType(ctx context.Context, id string) (*ent.EquipmentPortType, error) {
	e, err := r.ClientFrom(ctx).EquipmentPortType.Get(ctx, id)
	if err != nil && !ent.IsNotFound(err) {
		return nil, errors.Wrapf(err, "querying equipment port type: id=%q", id)
	}
	return e, nil
}

func (r queryResolver) EquipmentPortTypes(
	ctx context.Context,
	after *ent.Cursor, first *int,
	before *ent.Cursor, last *int,
) (*ent.EquipmentPortTypeConnection, error) {
	return r.ClientFrom(ctx).EquipmentPortType.Query().
		Paginate(ctx, after, first, before, last)
}

func (r queryResolver) EquipmentPortDefinitions(
	ctx context.Context,
	after *ent.Cursor, first *int,
	before *ent.Cursor, last *int,
) (*ent.EquipmentPortDefinitionConnection, error) {
	return r.ClientFrom(ctx).EquipmentPortDefinition.Query().
		Paginate(ctx, after, first, before, last)
}

func (r queryResolver) WorkOrder(ctx context.Context, id string) (*ent.WorkOrder, error) {
	wo, err := r.ClientFrom(ctx).WorkOrder.Get(ctx, id)
	if err != nil && !ent.IsNotFound(err) {
		return nil, errors.Wrapf(err, "querying equipment position definition: id=%q", id)
	}
	return wo, nil
}

func (r queryResolver) WorkOrders(
	ctx context.Context,
	after *ent.Cursor, first *int,
	before *ent.Cursor, last *int,
	showCompleted *bool,
) (*ent.WorkOrderConnection, error) {
	query := r.ClientFrom(ctx).WorkOrder.Query()
	if showCompleted != nil && !*showCompleted {
		query = query.Where(workorder.StatusIn(
			models.WorkOrderStatusPending.String(),
			models.WorkOrderStatusPlanned.String(),
		))
	}
	return query.Paginate(ctx, after, first, before, last)
}

func (r queryResolver) WorkOrderType(ctx context.Context, id string) (*ent.WorkOrderType, error) {
	lt, err := r.ClientFrom(ctx).WorkOrderType.Get(ctx, id)
	if err != nil && !ent.IsNotFound(err) {
		return nil, errors.Wrapf(err, "querying work order type: id=%q", id)
	}
	return lt, nil
}

func (r queryResolver) WorkOrderTypes(
	ctx context.Context,
	after *ent.Cursor, first *int,
	before *ent.Cursor, last *int,
) (*ent.WorkOrderTypeConnection, error) {
	return r.ClientFrom(ctx).WorkOrderType.Query().
		Paginate(ctx, after, first, before, last)
}

func (r queryResolver) SearchForEntity(
	ctx context.Context, name string,
	_ *ent.Cursor, limit *int,
	_ *ent.Cursor, _ *int,
) (*models.SearchEntriesConnection, error) {
	if limit == nil {
		return nil, errors.New("first is a mandatory param")
	}
	client := r.ClientFrom(ctx)
	locations, err := client.Location.Query().
		Where(
			location.Or(
				location.NameContainsFold(name),
				location.ExternalIDContainsFold(name),
			),
		).
		Limit(*limit).
		All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error searching for locations")
	}

	edges := make([]*models.SearchEntryEdge, len(locations))
	for i, l := range locations {
		lt := l.QueryType().FirstX(ctx)
		edges[i] = &models.SearchEntryEdge{
			Node: &models.SearchEntry{
				EntityType: "location",
				EntityID:   l.ID,
				Name:       l.Name,
				Type:       lt.Name,
			},
		}
	}
	if len(locations) == *limit {
		return &models.SearchEntriesConnection{Edges: edges}, nil
	}

	equipments, err := client.Equipment.Query().
		Where(equipment.Or(
			equipment.NameContainsFold(name),
			equipment.ExternalIDContainsFold(name),
		)).
		Limit(*limit - len(locations)).
		All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error searching for equipments")
	}
	for _, e := range equipments {
		et := e.QueryType().FirstX(ctx)
		edges = append(edges, &models.SearchEntryEdge{
			Node: &models.SearchEntry{
				EntityType: "equipment",
				EntityID:   e.ID,
				Name:       e.Name,
				Type:       et.Name,
			},
		})
	}
	return &models.SearchEntriesConnection{Edges: edges}, nil
}

func (r queryResolver) PossibleProperties(ctx context.Context, entityType models.PropertyEntity) ([]*ent.PropertyType, error) {
	var pts []*ent.PropertyType
	var err error
	switch entityType {
	case models.PropertyEntityEquipment:
		pts, err = r.ClientFrom(ctx).EquipmentType.Query().QueryPropertyTypes().All(ctx)
	case models.PropertyEntityService:
		pts, err = r.ClientFrom(ctx).ServiceType.Query().QueryPropertyTypes().All(ctx)
	case models.PropertyEntityLink:
		pts, err = r.ClientFrom(ctx).EquipmentPortType.Query().QueryLinkPropertyTypes().All(ctx)
	case models.PropertyEntityPort:
		pts, err = r.ClientFrom(ctx).EquipmentPortType.Query().QueryPropertyTypes().All(ctx)
	case models.PropertyEntityLocation:
		pts, err = r.ClientFrom(ctx).LocationType.Query().QueryPropertyTypes().All(ctx)
	default:
		return nil, errors.Errorf("entity type is not supported: %s", entityType)
	}

	if err != nil {
		return nil, errors.Wrap(err, "querying property types")
	}

	type key struct{ name, typ string }
	var (
		groups = map[key]struct{}{}
		types  []*ent.PropertyType
	)
	for _, pt := range pts {
		k := key{pt.Name, pt.Type}
		if _, ok := groups[k]; !ok {
			groups[k] = struct{}{}
			types = append(types, pt)
		}
	}
	return types, nil
}

func (r queryResolver) Surveys(ctx context.Context) ([]*ent.Survey, error) {
	surveys, err := r.ClientFrom(ctx).Survey.Query().All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "querying all surveys")
	}
	return surveys, nil
}

func (r queryResolver) Service(ctx context.Context, id string) (*ent.Service, error) {
	s, err := r.ClientFrom(ctx).Service.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "querying service: id=%q", id)
	}
	return s, nil
}

func (r queryResolver) ServiceType(ctx context.Context, id string) (*ent.ServiceType, error) {
	st, err := r.ClientFrom(ctx).ServiceType.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "querying service type: id=%q", id)
	}
	return st, nil
}

func (r queryResolver) ServiceTypes(
	ctx context.Context,
	after *ent.Cursor, first *int,
	before *ent.Cursor, last *int,
) (*ent.ServiceTypeConnection, error) {
	return r.ClientFrom(ctx).ServiceType.Query().
		Paginate(ctx, after, first, before, last)
}

func (r queryResolver) Customer(ctx context.Context, id string) (*ent.Customer, error) {
	st, err := r.ClientFrom(ctx).Customer.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "querying customer: id=%q", id)
	}
	return st, nil
}

func (r queryResolver) Customers(
	ctx context.Context,
	after *ent.Cursor, first *int,
	before *ent.Cursor, last *int,
) (*ent.CustomerConnection, error) {
	return r.ClientFrom(ctx).Customer.Query().
		Paginate(ctx, after, first, before, last)
}

func (r queryResolver) ActionsRules(
	ctx context.Context,
) (*models.ActionsRulesSearchResult, error) {
	results, err := r.ClientFrom(ctx).ActionsRule.Query().All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query action rules")
	}

	return &models.ActionsRulesSearchResult{
		Results: results,
		Count:   len(results),
	}, nil
}

func (r queryResolver) ActionsTriggers(
	ctx context.Context,
) (*models.ActionsTriggersSearchResult, error) {
	ac := actions.FromContext(ctx)
	triggers := ac.Triggers()

	ret := make([]*models.ActionsTrigger, len(triggers))
	for i, trigger := range ac.Triggers() {
		ret[i] = &models.ActionsTrigger{
			TriggerID:   trigger.ID(),
			Description: trigger.Description(),
		}
	}

	return &models.ActionsTriggersSearchResult{
		Results: ret,
		Count:   len(ret),
	}, nil
}

func (r queryResolver) ActionsTrigger(
	ctx context.Context, triggerID core.TriggerID,
) (*models.ActionsTrigger, error) {
	ac := actions.FromContext(ctx)
	trigger, err := ac.TriggerForID(triggerID)
	if err != nil {
		return nil, errors.Wrap(err, "getting trigger")
	}
	return &models.ActionsTrigger{
		TriggerID:   triggerID,
		Description: trigger.Description(),
	}, nil
}

func (r queryResolver) FindLocationWithDuplicateProperties(ctx context.Context, locationTypeID string, propertyName string) ([]string, error) {
	query := r.ClientFrom(ctx).
		LocationType.
		Query().
		Where(locationtype.ID(locationTypeID)).
		QueryLocations().
		QueryProperties().
		Where(property.HasTypeWith(
			propertytype.Name(propertyName),
		))
	properties, err := query.All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "querying properties")
	}

	var values []string
	for _, p := range properties {
		count, err := query.Clone().
			Where(property.StringVal(p.StringVal)).
			Count(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "querying count properties of properties with same value")
		}
		if count > 1 {
			values = append(values, p.StringVal)
		}
	}
	return values, nil
}

func (queryResolver) LatestPythonPackage(context.Context) (*models.LatestPythonPackageResult, error) {
	var packages []models.PythonPackage
	if err := json.Unmarshal([]byte(PyinventoryConsts), &packages); err != nil {
		return nil, errors.Wrap(err, "decoding python packages")
	}
	if len(packages) == 0 {
		return nil, nil
	}
	lastBreakingChange := len(packages) - 1
	for i, pkg := range packages {
		if pkg.HasBreakingChange {
			lastBreakingChange = i
			break
		}
	}
	return &models.LatestPythonPackageResult{
		LastPythonPackage:         &packages[0],
		LastBreakingPythonPackage: &packages[lastBreakingChange],
	}, nil
}

func (r queryResolver) Vertex(ctx context.Context, id string) (*ent.Node, error) {
	return r.ClientFrom(ctx).Node(ctx, id)
}
