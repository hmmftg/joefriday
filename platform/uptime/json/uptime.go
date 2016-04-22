// Copyright 2016 Joel Scoble and The JoeFriday authors.
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

// Package json handles JSON based processing of uptime information.  Instead
// of returning a Go struct, it returns JSON serialized bytes.  A function to
// deserialize the JSON serialized bytes into an uptime.Uptime struct is
// provided.
package json

import (
	"encoding/json"
	"sync"

	"github.com/mohae/joefriday/platform/uptime"
)

// Profiler is used to process the uptime information, /proc/version, using
// JSON.
type Profiler struct {
	*uptime.Profiler
}

// Initializes and returns a json.Profiler for uptime information.
func New() (prof *Profiler, err error) {
	p, err := uptime.New()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p}, nil
}

// Get returns the current uptime information as JSON serialized bytes.
func (prof *Profiler) Get() (p []byte, err error) {
	k, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(k)
}

var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current uptime information as JSON serialized bytes using
// the package's global Profiler.
func Get() (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = New()
		if err != nil {
			return nil, err
		}
	}
	return std.Get()
}

// Serialize uptime.Uptime using JSON
func (prof *Profiler) Serialize(u uptime.Uptime) ([]byte, error) {
	return json.Marshal(u)
}

// Serialize uptime.Uptime using JSON with the package global Profiler.
func Serialize(u uptime.Uptime) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = New()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(u)
}

// Marshal is an alias for Serialize
func (prof *Profiler) Marshal(u uptime.Uptime) ([]byte, error) {
	return prof.Serialize(u)
}

// Marshal is an alias for Serialize that uses the package's global profiler.
func Marshal(u uptime.Uptime) ([]byte, error) {
	return Serialize(u)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// uptime.Uptime.
func Deserialize(p []byte) (uptime.Uptime, error) {
	var u uptime.Uptime
	err := json.Unmarshal(p, &u)
	if err != nil {
		return u, err
	}
	return u, nil
}

// Unmarshal is an alias for Deserialize
func Unmarshal(p []byte) (uptime.Uptime, error) {
	return Deserialize(p)
}
