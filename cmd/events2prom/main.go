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
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rakyll/events-to-prom/engine"
	"github.com/rakyll/events-to-prom/event"

	_ "net/http/pprof"
)

var (
	config string
)

func main() {
	flag.StringVar(&config, "config", "events2prom.yaml", "")
	flag.Parse()

	conf, err := readConfig(config)
	if err != nil {
		log.Fatalf("Can't read the config at %q: %v", config, err)
	}

	events := make(chan event.Event, 32*1024)
	collections := make(chan engine.Collection, 32)
	removals := make(chan string, 32)

	loop := engine.NewLoop(conf.BufferSize, events, collections, removals)
	loop.BufferFlushWindow = conf.Window

	server := &eventsServer{port: conf.Port, events: events}
	admin := &adminServer{collections: collections, removals: removals}
	http.HandleFunc("/collections", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			admin.handlePost(w, r)
		case "DELETE":
			admin.handleDelete(w, r)
		}
	})
	http.Handle("/metrics", promhttp.HandlerFor(loop.Registry(), promhttp.HandlerOpts{}))

	// Register collections if any.
	for _, col := range conf.Collections {
		collections <- col
	}

	go server.listenAndServe()
	go loop.Run()

	log.Printf("Listening to admin server at %q...", conf.Endpoint)
	log.Fatal(http.ListenAndServe(conf.Endpoint, nil))
}
