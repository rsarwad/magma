// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build wireinject

package graphhttp

import (
	"net/http"
	"net/url"

	"github.com/facebookincubator/symphony/graph/event"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/pkg/actions/action/magmarebootnode"
	"github.com/facebookincubator/symphony/pkg/actions/executor"
	"github.com/facebookincubator/symphony/pkg/actions/trigger/magmaalert"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/mysql"
	"github.com/facebookincubator/symphony/pkg/oc"
	"github.com/facebookincubator/symphony/pkg/orc8r"
	"github.com/facebookincubator/symphony/pkg/server"
	"github.com/facebookincubator/symphony/pkg/server/xserver"
	"go.opencensus.io/stats/view"

	"github.com/google/wire"
	"github.com/gorilla/mux"
	"gocloud.dev/server/health"
)

// Config defines the http server config.
type Config struct {
	Tenancy    *viewer.MySQLTenancy
	AuthURL    *url.URL
	Emitter    event.Emitter
	Subscriber event.Subscriber
	Logger     log.Logger
	Census     oc.Options
	Orc8r      orc8r.Config
}

// NewServer creates a server from config.
func NewServer(cfg Config) (*server.Server, func(), error) {
	wire.Build(
		xserver.ServiceSet,
		provideViews,
		newHealthChecker,
		wire.FieldsOf(new(Config), "Tenancy", "Logger", "Census"),
		newRouterConfig,
		newRouter,
		wire.Bind(new(http.Handler), new(*mux.Router)),
	)
	return nil, nil, nil
}

func newHealthChecker(tenancy *viewer.MySQLTenancy) []health.Checker {
	return []health.Checker{tenancy}
}

func newRouterConfig(config Config) (cfg routerConfig, err error) {
	client, _ := orc8r.NewClient(config.Orc8r)
	registry := executor.NewRegistry()
	if err = registry.RegisterTrigger(magmaalert.New()); err != nil {
		return
	}
	if err = registry.RegisterAction(magmarebootnode.New(client)); err != nil {
		return
	}
	cfg = routerConfig{logger: config.Logger}
	cfg.viewer.tenancy = config.Tenancy
	cfg.viewer.authurl = config.AuthURL.String()
	cfg.events.emitter = config.Emitter
	cfg.events.subscriber = config.Subscriber
	cfg.orc8r.client = client
	cfg.actions.registry = registry
	return cfg, nil
}

func provideViews() []*view.View {
	views := xserver.DefaultViews()
	views = append(views, mysql.DefaultViews...)
	views = append(views, event.DefaultViews...)
	return views
}
