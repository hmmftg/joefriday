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

// Package json handles JSON based processing of CPU facts.  Instead of
// returning a Go struct, it returns JSON serialized bytes.  A function to
// deserialize the JSON serialized bytes into a facts.Facts struct is
// provided.
package json

import (
	"encoding/json"
	"sync"

	"github.com/mohae/joefriday/cpu/facts"
)

// Profiler is used to process the /proc/cpuinfo file.
type Profiler struct {
	*facts.Profiler
}

// Initializes and returns a cpu Facts profiler.
func NewProfiler() (prof *Profiler, err error) {
	p, err := facts.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p}, nil
}

// Get returns the current cpuinfo (Facts) as JSON serialized bytes.
func (prof *Profiler) Get() (p []byte, err error) {
	fct, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(fct)
}

var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current cpuinfo (Facts) as JSON serialized bytes using
// the package's global Profiler.
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

// Serialize cpu Facts as JSON
func (prof *Profiler) Serialize(fct *facts.Facts) ([]byte, error) {
	return json.Marshal(fct)
}

// Serialize cpu Facts as JSON using package globals.
func Serialize(fct *facts.Facts) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(fct)
}

// Marshal is an alias for serialize.
func (prof *Profiler) Marshal(fct *facts.Facts) ([]byte, error) {
	return prof.Serialize(fct)
}

// Marshal is an alias for Serialize using package globals.
func Marshal(fct *facts.Facts) ([]byte, error) {
	return std.Serialize(fct)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// facts.Facts
func Deserialize(p []byte) (*facts.Facts, error) {
	fct := &facts.Facts{}
	err := json.Unmarshal(p, fct)
	if err != nil {
		return nil, err
	}
	return fct, nil
}

// Unmarshal is an alias for Deserialize
func Unmarshal(p []byte) (*facts.Facts, error) {
	return Deserialize(p)
}
