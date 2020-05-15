// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewer_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/log/logtest"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestViewerHandler(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(*http.Request)
		expect  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "TestTenant",
			prepare: func(req *http.Request) {
				req.Header.Set(viewer.TenantHeader, "test")
				req.Header.Set(viewer.UserHeader, "user")
				req.Header.Set(viewer.RoleHeader, string(user.RoleOWNER))
			},
			expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "test", rec.Body.String())
			},
		},
		{
			name: "NoTenant",
			expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.NotZero(t, rec.Body.Len())
			},
		},
		{
			name: "NoRole",
			prepare: func(req *http.Request) {
				req.Header.Set(viewer.TenantHeader, "test")
				req.Header.Set(viewer.UserHeader, "user")
			},
			expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.NotZero(t, rec.Body.Len())
			},
		},
		{
			name: "WithUserEntNotExist",
			prepare: func(req *http.Request) {
				req.Header.Set(viewer.TenantHeader, "test")
				req.Header.Set(viewer.UserHeader, "new_user")
				req.Header.Set(viewer.RoleHeader, string(user.RoleOWNER))
			},
			expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
			},
		},
		{
			name: "WithNoUserInViewer",
			prepare: func(req *http.Request) {
				req.Header.Set(viewer.TenantHeader, "test")
				req.Header.Set(viewer.RoleHeader, string(user.RoleOWNER))
			},
			expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
	}

	client := viewertest.NewTestClient(t)
	ctx := ent.NewContext(context.Background(), client)
	_, err := client.User.Create().
		SetAuthID("user").
		Save(privacy.DecisionContext(ctx, privacy.Allow))
	require.NoError(t, err)

	h := viewer.TenancyHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			v := viewer.FromContext(ctx)
			require.NotNil(t, v)
			require.Equal(t, r.Header.Get(viewer.UserHeader), v.Name())
			assert.NotNil(t, log.FieldsFromContext(ctx))
			_, _ = io.WriteString(w, v.Tenant())
		}),
		viewer.NewFixedTenancy(client),
		logtest.NewTestLogger(t),
	)
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.prepare != nil {
				tc.prepare(req)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			tc.expect(t, rec)
		})
	}
}

func TestWebSocketUpgradeHandler(t *testing.T) {
	var upgrader websocket.Upgrader
	srv := httptest.NewServer(
		viewer.WebSocketUpgradeHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if websocket.IsWebSocketUpgrade(r) {
					conn, err := upgrader.Upgrade(w, r, nil)
					require.NoError(t, err)
					defer conn.Close()
					for {
						if _, _, err := conn.ReadMessage(); err != nil {
							return
						}
					}
				}
				w.WriteHeader(http.StatusOK)
			}),
			"",
		),
	)
	defer srv.Close()

	t.Run("NoUpgrade", func(t *testing.T) {
		rsp, err := srv.Client().Get(srv.URL)
		require.NoError(t, err)
		defer rsp.Body.Close()
		assert.Equal(t, http.StatusOK, rsp.StatusCode)
	})
	t.Run("AuthenticatedUpgrade", func(t *testing.T) {
		header := http.Header{}
		header.Set(viewer.TenantHeader, "test")
		u, _ := url.Parse(srv.URL)
		u.Scheme = "ws"
		conn, rsp, err := websocket.DefaultDialer.DialContext(
			context.Background(), u.String(), header,
		)
		require.NoError(t, err)
		rsp.Body.Close()
		conn.Close()
	})

	const host = "test.example.com"
	authenticator := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, host, r.Header.Get("X-Forwarded-Host"))
			if username, _, ok := r.BasicAuth(); ok {
				err := json.NewEncoder(w).Encode(map[string]string{
					"organization": "test",
					"email":        username,
					"role":         "user",
				})
				require.NoError(t, err)
			} else if cookie, err := r.Cookie("Viewer"); err == nil {
				data, err := base64.StdEncoding.DecodeString(cookie.Value)
				require.NoError(t, err)
				_, err = w.Write(data)
				require.NoError(t, err)
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		}),
	)
	defer authenticator.Close()

	t.Run("Auth", func(t *testing.T) {
		authenticate := func(t *testing.T, authReq func(*http.Request)) {
			var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set(viewer.TenantHeader, r.Header.Get(viewer.TenantHeader))
				w.Header().Set(viewer.UserHeader, r.Header.Get(viewer.UserHeader))
				w.Header().Set(viewer.RoleHeader, r.Header.Get(viewer.RoleHeader))
				w.WriteHeader(http.StatusOK)
			})
			handler = viewer.WebSocketUpgradeHandler(handler, authenticator.URL)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Host = host
			authReq(req)
			req.Header.Set("Connection", "upgrade")
			req.Header.Set("Upgrade", "websocket")
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			require.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "test", rec.Header().Get(viewer.TenantHeader))
			assert.Equal(t, "tester", rec.Header().Get(viewer.UserHeader))
			assert.Equal(t, "user", rec.Header().Get(viewer.RoleHeader))
		}

		t.Run("Basic", func(t *testing.T) {
			authenticate(t, func(req *http.Request) {
				req.SetBasicAuth("tester", "tester")
			})
		})
		t.Run("Session", func(t *testing.T) {
			data, err := json.Marshal(map[string]string{
				"organization": "test",
				"email":        "tester",
				"role":         "user",
			})
			require.NoError(t, err)
			authenticate(t, func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  "Viewer",
					Value: base64.StdEncoding.EncodeToString(data),
				})
			})
		})
	})
}

type mockHandler struct{ mock.Mock }

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func TestDeactivatedUser(t *testing.T) {
	client := viewertest.NewTestClient(t)
	ctx := ent.NewContext(context.Background(), client)
	ctx = privacy.DecisionContext(ctx, privacy.Allow)
	deactivatedUser, err := client.User.Create().
		SetAuthID("deactivated_user").
		SetStatus(user.StatusDEACTIVATED).
		Save(ctx)
	require.NoError(t, err)
	v := viewer.NewUser("test", deactivatedUser)
	ctx = viewer.NewContext(ctx, v)
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	var m mockHandler
	defer m.AssertExpectations(t)
	viewer.UserHandler(&m, logtest.NewTestLogger(t)).ServeHTTP(rec, req)
	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Equal(t, "user is deactivated\n", rec.Body.String())
}

func TestAutomationViewerIsNotDeactivated(t *testing.T) {
	client := viewertest.NewTestClient(t)
	ctx := ent.NewContext(context.Background(), client)
	v := viewer.NewAutomation("test", "automation", user.RoleADMIN)
	ctx = viewer.NewContext(ctx, v)
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	var m mockHandler
	m.On("ServeHTTP", mock.Anything, mock.Anything).Once()
	defer m.AssertExpectations(t)
	viewer.UserHandler(&m, logtest.NewTestLogger(t)).ServeHTTP(rec, req)
}

func TestViewerMarshalLog(t *testing.T) {
	core, o := observer.New(zap.InfoLevel)
	logger := zap.New(core)
	c := viewertest.NewTestClient(t)
	ctx := privacy.DecisionContext(context.Background(), privacy.Allow)
	u, err := c.User.Create().SetAuthID("tester").Save(ctx)
	require.NoError(t, err)
	v := viewer.NewUser("test", u)
	logger.Info("viewer log test", zap.Object("viewer", v))

	logs := o.TakeAll()
	require.Len(t, logs, 1)
	field, ok := logs[0].ContextMap()["viewer"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, v.Tenant(), field["tenant"])
	assert.Equal(t, u.AuthID, field["user"])
}

type testExporter struct {
	mock.Mock
}

func (te *testExporter) ExportSpan(s *trace.SpanData) {
	te.Called(s)
}

func TestViewerSpanAttributes(t *testing.T) {
	client := viewertest.NewTestClient(t)
	h := viewer.TenancyHandler(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusAccepted)
		}),
		viewer.NewFixedTenancy(client),
		logtest.NewTestLogger(t),
	)
	t.Run("WithSpan", func(t *testing.T) {
		var te testExporter
		trace.RegisterExporter(&te)
		defer trace.UnregisterExporter(&te)

		te.On("ExportSpan", mock.AnythingOfType("*trace.SpanData")).
			Run(func(args mock.Arguments) {
				s := args.Get(0).(*trace.SpanData)
				assert.Equal(t, viewertest.DefaultTenant, s.Attributes["viewer.tenant"])
				assert.Equal(t, viewertest.DefaultUser, s.Attributes["viewer.user"])
			}).
			Once()
		defer te.AssertExpectations(t)

		ctx, span := trace.StartSpan(context.Background(), "test",
			trace.WithSampler(trace.AlwaysSample()),
		)
		defer span.End()

		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
		viewertest.SetDefaultViewerHeaders(req)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusAccepted, rec.Code)
	})
	t.Run("WithoutSpan", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		viewertest.SetDefaultViewerHeaders(req)
		rec := httptest.NewRecorder()
		assert.NotPanics(t, func() { h.ServeHTTP(rec, req) })
		assert.Equal(t, http.StatusAccepted, rec.Code)
	})
}

func TestViewerTags(t *testing.T) {
	measure := stats.Int64("", "", stats.UnitDimensionless)
	v := &view.View{
		Name: "viewer/test_tags",
		TagKeys: []tag.Key{
			viewer.KeyTenant,
			viewer.KeyUser,
			viewer.KeyRole,
		},
		Measure:     measure,
		Aggregation: view.Count(),
	}
	err := view.Register(v)
	require.NoError(t, err)
	defer view.Unregister(v)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(viewer.TenantHeader, "test-tenant")
	req.Header.Set(viewer.UserHeader, "test-user")
	req.Header.Set(viewer.RoleHeader, string(user.RoleUSER))
	rec := httptest.NewRecorder()
	client := viewertest.NewTestClient(t)
	viewer.TenancyHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := stats.RecordWithTags(r.Context(), nil, measure.M(1))
			require.NoError(t, err)
			w.WriteHeader(http.StatusNoContent)
		}),
		viewer.NewFixedTenancy(client),
		logtest.NewTestLogger(t),
	).ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	rows, err := view.RetrieveData(v.Name)
	require.NoError(t, err)
	require.Len(t, rows, 1)

	hasTag := func(key tag.Key, value string) assert.Comparison {
		return func() bool {
			for _, t := range rows[0].Tags {
				if t.Key.Name() == key.Name() {
					return t.Value == value
				}
			}
			return false
		}
	}
	assert.Condition(t, hasTag(viewer.KeyTenant, "test-tenant"))
	assert.Condition(t, hasTag(viewer.KeyUser, "test-user"))
	assert.Condition(t, hasTag(viewer.KeyRole, "USER"))
}

func TestViewerTenancy(t *testing.T) {
	t.Run("Regular", func(t *testing.T) {
		client := viewertest.NewTestClient(t)
		h := viewer.TenancyHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.True(t, client == ent.FromContext(r.Context()))
				v, ok := viewer.FromContext(r.Context()).(*viewer.UserViewer)
				assert.True(t, ok)
				assert.Equal(t, viewertest.DefaultTenant, v.Tenant())
				assert.Equal(t, viewertest.DefaultRole, v.Role())
				assert.Equal(t, viewertest.DefaultUser, v.Name())
				w.WriteHeader(http.StatusAccepted)
			}),
			viewer.NewFixedTenancy(client),
			logtest.NewTestLogger(t),
		)
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		viewertest.SetDefaultViewerHeaders(req)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusAccepted, rec.Code)
	})
	t.Run("WithFeatures", func(t *testing.T) {
		client := viewertest.NewTestClient(t)
		h := viewer.TenancyHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				v := viewer.FromContext(r.Context())
				assert.True(t, v.Features().Enabled("feature1"))
				assert.True(t, v.Features().Enabled("feature2"))
				assert.False(t, v.Features().Enabled("feature3"))
				assert.Equal(t, "feature1,feature2", v.Features().String())
				w.WriteHeader(http.StatusAccepted)
			}),
			viewer.NewFixedTenancy(client),
			logtest.NewTestLogger(t),
		)
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		viewertest.SetDefaultViewerHeaders(req)
		req.Header.Set(viewer.FeaturesHeader, "feature1,feature2")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusAccepted, rec.Code)
	})
	t.Run("WithAutomation", func(t *testing.T) {
		client := viewertest.NewTestClient(t)
		name := "Scheduler"
		h := viewer.TenancyHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.True(t, client == ent.FromContext(r.Context()))
				v, ok := viewer.FromContext(r.Context()).(*viewer.AutomationViewer)
				assert.True(t, ok)
				assert.Equal(t, viewertest.DefaultTenant, v.Tenant())
				assert.Equal(t, viewertest.DefaultRole, v.Role())
				assert.Equal(t, name, v.Name())
				w.WriteHeader(http.StatusAccepted)
			}),
			viewer.NewFixedTenancy(client),
			logtest.NewTestLogger(t),
		)
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(viewertest.TenantHeader, viewertest.DefaultTenant)
		req.Header.Set(viewer.AutomationHeader, name)
		req.Header.Set(viewertest.RoleHeader, string(viewertest.DefaultRole))
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusAccepted, rec.Code)
	})
	t.Run("WithNoViewer", func(t *testing.T) {
		client := viewertest.NewTestClient(t)
		h := viewer.TenancyHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			viewer.NewFixedTenancy(client),
			logtest.NewTestLogger(t),
		)
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(viewertest.TenantHeader, viewertest.DefaultTenant)
		req.Header.Set(viewertest.RoleHeader, string(viewertest.DefaultRole))
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
