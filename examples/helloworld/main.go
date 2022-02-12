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
	"log"
	"math/rand"
	"time"

	events2prom "github.com/rakyll/events2prom"
	"github.com/rakyll/events2prom/event"
)

var pods = []string{"pod-1e0", "pod-1ff", "pod-def"}

func main() {
	for {
		n := rand.Intn(400)
		events2prom.Publish(
			event.Event{
				Name:  "request_latency_ms",
				Value: float64(n),
				Labels: map[string]string{
					"pod": pods[rand.Intn(len(pods))],
				},
			},
			event.Event{
				Name:  "event_not_collected",
				Value: float64(n),
			},
		)
		log.Println("Published two events.")
		time.Sleep(10 * time.Millisecond)
	}
}
