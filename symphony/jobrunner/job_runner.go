// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jobrunner

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/facebookincubator/symphony/graph/graphgrpc/schema"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

const (
	grpcAddr         = "%s:443"
	jobsURL          = "http://%s/jobs/%s"
	graphHostEnv     = "GRAPH_HOST"
	defaultGraphHost = "graph"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getTenantList(client schema.TenantServiceClient) []string {
	var tenantNames []string
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Printf("getting tenant list")
	tenants, err := client.List(ctx, &empty.Empty{})
	if err != nil {
		log.Fatalf("failed to fetch list of tenants: %v", err)
	}
	for _, tenant := range tenants.Tenants {
		tenantNames = append(tenantNames, tenant.Name)
	}
	return tenantNames
}

func runJob(url, tenant string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("x-auth-organization", tenant)
	req.Header.Add("x-auth-automation-name", "job_runner")
	req.Header.Add("x-auth-user-role", "OWNER")
	log.Printf("running job on url %s, tenant %s", url, tenant)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get response: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}
	log.Printf(
		"job with url %s, tenant %s succeeded with status %s output \"%s\"",
		url,
		tenant,
		resp.Status,
		body)
	return nil
}

func RunJobOnAllTenants(jobs ...string) {
	graphHost := getEnv(graphHostEnv, defaultGraphHost)
	address := fmt.Sprintf(grpcAddr, graphHost)
	log.Printf("connecting to grpc server %s", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := schema.NewTenantServiceClient(conn)
	tenants := getTenantList(client)
	for _, tenant := range tenants {
		for _, job := range jobs {
			url := fmt.Sprintf(jobsURL, graphHost, job)
			err := runJob(url, tenant)
			if err != nil {
				log.Printf("failed connecting url %s: %v", url, err)
			}
		}
	}
}
