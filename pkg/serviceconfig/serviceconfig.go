/*
 * Copyright 2021 OpsMx, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License")
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

//
// Package serviceconfig handles the top level agent configuration YAML handling.
//
// The service-level "Config" element is handled by the Make() methods on
// each service endpoint type.
//
package serviceconfig

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type incomingServiceConfig struct {
	Name               string `yaml:"name,omitempty"`
	Port               uint16 `yaml:"port,omitempty"`
	ServiceType        string `yaml:"serviceType,omitempty"`
	UseHTTP            bool   `yaml:"useHTTP,omitempty"`
	Destination        string `yaml:"destination,omitempty"`
	DestinationService string `yaml:"destinationService,omitempty"`
}

type outgoingServiceConfig struct {
	Enabled    bool                        `yaml:"enabled"`
	Name       string                      `yaml:"name"`
	Type       string                      `yaml:"type"`
	Config     map[interface{}]interface{} `yaml:"config,omitempty"`
	Namespaces []serviceNamespace          `yaml:"namespaces,omitempty"`
	AccountID  string                      `yaml:"accountId,omitempty"`
	AssumeRole string                      `yaml:"assumeRole,omitempty"`
}

type serviceNamespace struct {
	Name       string   `yaml:"name"`
	Namespaces []string `yaml:"namespaces"`
}

// ServiceConfig defines a service level configuration top-level list.
type ServiceConfig struct {
	OutgoingServices []outgoingServiceConfig `yaml:"outgoingServices,omitempty"`
	IncomingServices []incomingServiceConfig `yaml:"incomingServices,omitempty"`
}

// LoadServiceConfig loads a service configuration YAML file.
func LoadServiceConfig(filename string) (*ServiceConfig, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &ServiceConfig{}
	err = yaml.Unmarshal(buf, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
