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

// Package json handles JSON based processing of OS release information,
// /etc/os-release.  Instead of returning a Go struct, it returns JSON
// serialized bytes.  A function to deserialize the JSON serialized bytes
// into a release.Release struct is provided.
package json

import (
	"encoding/json"
	"sync"

	"github.com/mohae/joefriday/platform/release"
)

// Profiler is used to process the OS release information file using JSON.
type Profiler struct {
	*release.Profiler
}

// Initializes and returns a json.Profiler for OS release information.
func NewProfiler() (prof *Profiler, err error) {
	p, err := release.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p}, nil
}

// Get returns the current OS release information as JSON serialized bytes.
func (prof *Profiler) Get() (p []byte, err error) {
	k, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(k)
}

var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current OS release information as JSON serialized bytes
// using the package's global Profiler.
func Get() (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Get()
}

// Serialize release.Release using JSON
func (prof *Profiler) Serialize(r *release.Release) ([]byte, error) {
	return json.Marshal(r)
}

// Serialize release.Release using JSON with the package global Profiler.
func Serialize(r *release.Release) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(r)
}

// Marshal is an alias for Serialize
func (prof *Profiler) Marshal(r *release.Release) ([]byte, error) {
	return prof.Serialize(r)
}

// Marshal is an alias for Serialize using the package's global profiler.
func Marshal(r *release.Release) ([]byte, error) {
	return Serialize(r)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// release.Release.
func Deserialize(p []byte) (*release.Release, error) {
	k := &release.Release{}
	err := json.Unmarshal(p, k)
	if err != nil {
		return nil, err
	}
	return k, nil
}

// Unmarshal is an alias for Deserialize
func Unmarshal(p []byte) (*release.Release, error) {
	return Deserialize(p)
}
