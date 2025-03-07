/*
 * Copyright 2021 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"log"
	"time"

	"github.com/cloudwego/kitex-examples/hello/kitex_gen/api"
	"github.com/cloudwego/kitex-examples/hello/kitex_gen/api/hello"
	"github.com/cloudwego/kitex/client"
	polaris "github.com/kitex-contrib/registry-polaris"
	"github.com/polarismesh/polaris-go/pkg/config"
)

const (
	confPath  = "polaris.yaml"
	Namespace = "Polaris"
	// At present,polaris server tag is v1.4.0，can't support auto create namespace,
	// if you want to use a namespace other than default,Polaris ,before you register an instance,
	// you should create the namespace at polaris console first.
)

func main() {
	Conf, err := config.LoadConfigurationByFile(confPath)
	if err != nil {
		log.Fatal(err)
	}
	polarisAddresses := Conf.Global.ServerConnector.Addresses

	r, err := polaris.NewPolarisResolver(polarisAddresses)
	if err != nil {
		log.Fatal(err)
	}

	// client.WithTag sets the namespace tag for service discovery
	newClient := hello.MustNewClient("echo", client.WithTag("namespace", Namespace),
		client.WithResolver(r), client.WithRPCTimeout(time.Second*60))
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		resp, err := newClient.Echo(ctx, &api.Request{Message: "Hi,polaris!"})
		cancel()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(resp)
		time.Sleep(time.Second)
	}
}
