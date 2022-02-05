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
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/rakyll/events-to-prom/engine"
	"github.com/rakyll/events-to-prom/event"
)

type eventsServer struct {
	port   int
	events chan event.Event
}

func (s *eventsServer) listenAndServe() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: s.port})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Printf("Listening events at %v, let's 🧹...", conn.LocalAddr())
	message := make([]byte, 2048)
	for {
		len, _, err := conn.ReadFromUDP(message[:])
		if err != nil {
			log.Printf("Cannot read event: %v", err)
			continue
		}
		event, err := event.Parse(message[:len-1])
		if err != nil {
			log.Printf("Error parsing event: %s", message)
			continue
		}
		s.events <- event
	}
}

type adminServer struct {
	collections chan engine.Collection
	removals    chan string
}

func (s *adminServer) handlePost(w http.ResponseWriter, r *http.Request) {
	var col engine.Collection
	if err := json.NewDecoder(r.Body).Decode(&col); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.collections <- col
}

func (s *adminServer) handleDelete(w http.ResponseWriter, r *http.Request) {
	var col struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&col); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.removals <- col.Name
}
