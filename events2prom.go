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

package events2prom

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/rakyll/events2prom/event"
)

var defaultAddr = "0.0.0.0:6678"

var defaultClient *Client

func init() {
	var err error
	if host := os.Getenv("EVENTS2PROM_HOST"); host != "" {
		defaultAddr = host + ":6678"
	}
	// TODO(jbd): Only dial if Publish is called.
	defaultClient, err = NewClient(defaultAddr)
	if err != nil {
		log.Println(err)
	}
}

func Publish(e ...event.Event) {
	defaultClient.Publish(e...)
}

type Client struct {
	conn net.Conn
}

func NewClient(addr string) (*Client, error) {
	if addr == "" {
		addr = defaultAddr
	}
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

func (c *Client) Publish(e ...event.Event) {
	for _, ee := range e {
		fmt.Fprintln(c.conn, ee.Text())
	}
}

func (c *Client) Close() error {
	return c.conn.Close()
}
