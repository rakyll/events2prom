// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"
	"time"

	"github.com/rakyll/events-to-prom/engine"
	"gopkg.in/yaml.v2"
)

const (
	defaultPort        = 6678
	defaultEndpoint    = ":6677"
	defaultFlushWindow = 5 * time.Second
)

type serverConfig struct {
	// Port is the UDP port to listen to events.
	Port int `yaml:"port,omitempty"`

	// Endpoint is the endpoint to serve the control API.
	// Users can enable or disable new aggregation using the API.
	// Control API also serves the metrics in the Prometheus
	// exposition format at {Endpoint}/metrics.
	Endpoint string `yaml:"endpoint,omitempty"`

	// BufferSize is the max number of events to buffer in memory
	// before starting to aggregate.
	BufferSize int `yaml:"buffer_size,omitempty"`

	// Window is the max amount of duration to buffer events in
	// memory before starting to aggregate. Example values are:
	// 1s for one second, 2m for two minutes, 1m30s for one minute
	// thirty seconds.
	Window time.Duration `yaml:"window,omitempty"`

	// Collections are the collections to enable at start.
	Collections []engine.Collection `yaml:"collections,omitempty"`
}

func readConfig(filename string) (serverConfig, error) {
	var c serverConfig
	if filename != "" {
		f, err := os.Open(filename)
		if err != nil {
			return serverConfig{}, err
		}
		defer f.Close()

		decoder := yaml.NewDecoder(f)
		if err := decoder.Decode(&c); err != nil {
			return serverConfig{}, err
		}
	}
	if c.Port == 0 {
		c.Port = defaultPort
	}
	if c.Endpoint == "" {
		c.Endpoint = defaultEndpoint
	}
	if c.Window <= 0 {
		c.Window = defaultFlushWindow
	}
	return c, nil
}
