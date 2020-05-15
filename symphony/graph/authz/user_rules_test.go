// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"errors"
	"testing"

	"github.com/facebookincubator/symphony/graph/authz"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/stretchr/testify/require"
)

func TestUserCannotBeDeleted(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	u, err := c.User.Create().SetAuthID("someone").Save(ctx)
	require.NoError(t, err)
	err = c.User.DeleteOne(u).Exec(ctx)
	require.True(t, errors.Is(err, privacy.Deny))
}

func TestAdminUserCanEditUsers(t *testing.T) {
	client := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(
		context.Background(),
		client,
		viewertest.WithRole(user.RoleADMIN),
		viewertest.WithPermissions(authz.AdminPermissions()))
	_, err := client.UsersGroup.Create().
		SetName("NewGroup").
		Save(ctx)
	require.NoError(t, err)
	_, err = client.User.Create().
		SetAuthID("new_user").
		Save(ctx)
	require.NoError(t, err)
}
