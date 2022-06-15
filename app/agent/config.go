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

package main

import (
	"io/ioutil"

	"github.com/opsmx/oes-birger/internal/tunnel"
	"gopkg.in/yaml.v3"
)

const (
	defaultCertPath       = "/app/secrets/agent/tls.crt"
	defaultKeyPath        = "/app/secrets/agent/tls.key"
	defaultUserconfigPath = "/app/config/services.yaml"
)

// agentConfig holds all the configuration for the agent.  The
// configuration file is loaded from disk first, and then any
// environment variables are applied.
type agentConfig struct {
	ControllerHostname        string  `yaml:"controllerHostname,omitempty"`
	CACert64                  *string `yaml:"caCert64,omitempty"`
	CertFile                  string  `yaml:"certFile,omitempty"`
	KeyFile                   string  `yaml:"keyFile,omitempty"`
	ServicesConfigPath        string  `yaml:"servicesConfigPath,omitempty"`
	InsecureControllerAllowed bool    `yaml:"insecureControllerAllowed,omitempty"`
}

func (c *agentConfig) applyDefaults() {
	if len(c.ControllerHostname) == 0 {
		c.ControllerHostname = "forwarder-controller:9001"
	}

	if len(c.CertFile) == 0 {
		c.CertFile = defaultCertPath
	}

	if len(c.KeyFile) == 0 {
		c.KeyFile = defaultKeyPath
	}

	if len(c.ServicesConfigPath) == 0 {
		c.ServicesConfigPath = defaultUserconfigPath
	}
}

// loadConfig will load YAML configuration from the provided filename, and then apply
// environment variables to override some subset of available options.
func loadConfig(filename string) (*agentConfig, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &agentConfig{}
	err = yaml.Unmarshal(buf, config)
	if err != nil {
		return nil, err
	}

	config.applyDefaults()

	return config, nil
}

type AgentInfoContainer struct {
	AgentInfo tunnel.AgentInfo `yaml:"agentInfo,omitempty"`
}

func loadAgentInfo(filename string) (*tunnel.AgentInfo, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	info := AgentInfoContainer{}
	err = yaml.Unmarshal(buf, &info)
	if err != nil {
		return nil, err
	}

	return &info.AgentInfo, nil
}
