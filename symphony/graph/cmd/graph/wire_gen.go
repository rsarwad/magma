// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"context"
	"fmt"
	"github.com/facebookincubator/symphony/graph/event"
	"github.com/facebookincubator/symphony/graph/graphevents"
	"github.com/facebookincubator/symphony/graph/graphgrpc"
	"github.com/facebookincubator/symphony/graph/graphhttp"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/mysql"
	"github.com/facebookincubator/symphony/pkg/server"
	"gocloud.dev/server/health"
	"google.golang.org/grpc"
)

import (
	_ "github.com/facebookincubator/symphony/graph/ent/runtime"
	_ "gocloud.dev/pubsub/mempubsub"
	_ "gocloud.dev/pubsub/natspubsub"
)

// Injectors from wire.go:

func newApplication(ctx context.Context, flags *cliFlags) (*application, func(), error) {
	config := flags.LogConfig
	logger, cleanup, err := log.ProvideLogger(config)
	if err != nil {
		return nil, nil, err
	}
	mysqlConfig := flags.MySQLConfig
	mySQLTenancy, err := newMySQLTenancy(mysqlConfig, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	eventConfig := flags.EventConfig
	topicEmitter, cleanup2, err := event.ProvideEmitter(ctx, eventConfig)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	tenancy, err := newTenancy(mySQLTenancy, logger, topicEmitter)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	url := flags.AuthURL
	urlSubscriber := event.ProvideSubscriber(eventConfig)
	telemetryConfig := &flags.TelemetryConfig
	v := newHealthChecks(mySQLTenancy)
	orc8rConfig := flags.Orc8rConfig
	graphhttpConfig := graphhttp.Config{
		Tenancy:      tenancy,
		AuthURL:      url,
		Subscriber:   urlSubscriber,
		Logger:       logger,
		Telemetry:    telemetryConfig,
		HealthChecks: v,
		Orc8r:        orc8rConfig,
	}
	server, cleanup3, err := graphhttp.NewServer(graphhttpConfig)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	config2 := &flags.MySQLConfig
	db, cleanup4 := mysql.Provider(config2)
	graphgrpcConfig := graphgrpc.Config{
		DB:      db,
		Logger:  logger,
		Orc8r:   orc8rConfig,
		Tenancy: tenancy,
	}
	grpcServer, cleanup5, err := graphgrpc.NewServer(graphgrpcConfig)
	if err != nil {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	grapheventsConfig := graphevents.Config{
		Tenancy:    tenancy,
		Subscriber: urlSubscriber,
		Logger:     logger,
	}
	grapheventsServer, cleanup6, err := graphevents.NewServer(grapheventsConfig)
	if err != nil {
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	mainApplication := newApp(logger, server, grpcServer, grapheventsServer, flags)
	return mainApplication, func() {
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

func newApp(logger log.Logger, httpServer *server.Server, grpcServer *grpc.Server, eventServer *graphevents.Server, flags *cliFlags) *application {
	var app application
	app.Logger = logger.Background()
	app.http.Server = httpServer
	app.http.addr = flags.HTTPAddress.String()
	app.grpc.Server = grpcServer
	app.grpc.addr = flags.GRPCAddress.String()
	app.event = eventServer
	return &app
}

func newTenancy(tenancy *viewer.MySQLTenancy, logger log.Logger, emitter event.Emitter) (viewer.Tenancy, error) {
	eventer := event.Eventer{Logger: logger, Emitter: emitter}
	return viewer.NewCacheTenancy(tenancy, eventer.HookTo), nil
}

func newHealthChecks(tenancy *viewer.MySQLTenancy) []health.Checker {
	return []health.Checker{tenancy}
}

func newMySQLTenancy(config mysql.Config, logger log.Logger) (*viewer.MySQLTenancy, error) {
	tenancy, err := viewer.NewMySQLTenancy(config.String())
	if err != nil {
		return nil, fmt.Errorf("creating mysql tenancy: %w", err)
	}
	tenancy.SetLogger(logger)
	mysql.SetLogger(logger)
	return tenancy, nil
}
