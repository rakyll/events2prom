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

package event

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/valyala/fastjson"
)

var fastParser fastjson.Parser

type Event struct {
	Name      string            `json:"event,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"` // label keys should match the regex [a-zA-Z0-9_]*
	Value     float64           `json:"value,omitempty"`
	Timestamp time.Time         `json:"ts,omitempty"`
}

func (e Event) Text() string {
	var buf bytes.Buffer
	buf.WriteString(e.Name)
	buf.WriteByte('|')
	buf.WriteString(strconv.FormatFloat(e.Value, 'f', -1, 64))
	buf.WriteByte('|')
	buf.WriteString(strconv.FormatInt(e.Timestamp.UnixNano(), 10))
	for k, v := range e.Labels {
		buf.WriteByte('|')
		buf.WriteString(k)
		buf.WriteByte(':')
		buf.WriteString(v)
	}
	return buf.String()
}

func ParseJSON(buf []byte) (Event, error) {
	v, err := fastParser.ParseBytes(buf)
	if err != nil {
		return Event{}, err
	}
	name := string(v.GetStringBytes("event"))
	value := v.GetFloat64("value")
	o, err := v.Object()
	if err != nil {
		return Event{}, err
	}

	labels := make(map[string]string)
	o.Visit(func(k []byte, v *fastjson.Value) {
		labels[string(k)] = v.String()
	})
	// TODO(jbd): Handle timestamp.
	return Event{
		Name:   name,
		Value:  value,
		Labels: labels,
	}, nil
}

func Parse(buf []byte) (Event, error) {
	// TODO(jbd): Remove bytes.Split.
	const minSections = 3
	sections := bytes.Split(buf, []byte("|"))
	if len(sections) < minSections {
		return Event{}, errors.New("invalid event")
	}

	name := string(sections[0])
	value, err := strconv.ParseFloat(string(sections[1]), 64)
	if err != nil {
		return Event{}, err
	}

	labels := make(map[string]string, len(sections)-minSections)
	for i := minSections; i < len(sections); i++ {
		keyValue := sections[i]
		idx := bytes.IndexByte(keyValue, byte(':'))
		if idx <= 0 {
			return Event{}, fmt.Errorf("invalid label: %s", keyValue)
		}
		labels[string(keyValue[:idx])] = string(keyValue[idx+1:])
	}
	// TODO(jbd): Handle timestamp.
	return Event{
		Name:   name,
		Value:  value,
		Labels: labels,
	}, nil
}
