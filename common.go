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

package polaris

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/cloudwego/kitex/pkg/discovery"
	perrors "github.com/pkg/errors"
	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/config"
	"github.com/polarismesh/polaris-go/pkg/model"
)

// GetPolarisConfig get polaris config from endpoints.
func GetPolarisConfig(endpoints []string) (api.SDKContext, error) {
	if len(endpoints) == 0 {
		return nil, perrors.New("endpoints is empty!")
	}

	serverConfigs := make([]string, 0, len(endpoints))
	for _, addr := range endpoints {
		ip, portStr, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, perrors.WithMessagef(err, "split [%s] ", addr)
		}
		port, _ := strconv.Atoi(portStr)
		serverConfigs = append(serverConfigs, fmt.Sprintf("%s:%d", ip, uint64(port)))
	}

	polarisConf := config.NewDefaultConfiguration(serverConfigs)

	sdkCtx, err := api.InitContextByConfig(polarisConf)
	if err != nil {
		return nil, err
	}
	return sdkCtx, nil
}

// SplitDescription splits description to namespace and serviceName.
func SplitDescription(description string) (string, string) {
	str := strings.Split(description, ":")
	namespace, serviceName := str[0], str[1]
	return namespace, serviceName
}

// ChangePolarisInstanceToKitex transforms polaris instance to Kitex instance.
func ChangePolarisInstanceToKitex(PolarisInstance model.Instance) discovery.Instance {
	weight := PolarisInstance.GetWeight()
	if weight <= 0 {
		weight = defaultWeight
	}
	addr := PolarisInstance.GetHost() + ":" + strconv.Itoa(int(PolarisInstance.GetPort()))

	tags := map[string]string{
		"namespace": PolarisInstance.GetNamespace(),
	}

	KitexInstance := discovery.NewInstance(PolarisInstance.GetProtocol(), addr, weight, tags)
	// In KitexInstance , tags can be used as IDC、Cluster、Env 、namespace、and so on.
	return KitexInstance
}

// GetLocalIPv4Address gets local ipv4 address when info host is empty.
func GetLocalIPv4Address() (string, error) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addr {
		ipNet, isIpNet := addr.(*net.IPNet)
		if isIpNet && !ipNet.IP.IsLoopback() {
			ipv4 := ipNet.IP.To4()
			if ipv4 != nil {
				return ipv4.String(), nil
			}
		}
	}
	return "", fmt.Errorf("not found ipv4 address")
}

// GetInfoHostAndPort gets Host and port from info.Addr.
func GetInfoHostAndPort(Addr string) (string, int, error) {
	infoHost, port, err := net.SplitHostPort(Addr)
	if err != nil {
		return "", 0, err
	} else {
		if port == "" {
			return infoHost, 0, fmt.Errorf("registry info addr missing port")
		}
		if infoHost == "" {
			ipv4, err := GetLocalIPv4Address()
			if err != nil {
				return "", 0, fmt.Errorf("get local ipv4 error, cause %v", err)
			}
			infoHost = ipv4
		}
	}
	infoPort, err := strconv.Atoi(port)
	if err != nil {
		return "", 0, err
	}
	return infoHost, infoPort, nil
}

// GetInstanceKey generates instanceKey  for one instance.
func GetInstanceKey(namespace, serviceName, host, port string) string {
	var instanceKey strings.Builder
	instanceKey.WriteString(namespace)
	instanceKey.WriteString(":")
	instanceKey.WriteString(serviceName)
	instanceKey.WriteString(":")
	instanceKey.WriteString(host)
	instanceKey.WriteString(":")
	instanceKey.WriteString(port)
	return instanceKey.String()
}
