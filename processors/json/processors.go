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

// Package json handles JSON based processing of Processor info.  Instead of
// returning a Go struct, it returns JSON serialized bytes.  A function to
// deserialize the JSON serialized bytes into a processors.Processors struct
// is provided.
package json

import (
	"encoding/json"
	"sync"

	"github.com/mohae/joefriday/processors"
)

// Profiler is used to process the processor info as JSON serialized bytes.
type Profiler struct {
	*processors.Profiler
}

// Initializes and returns a cpu Facts piler.
func NewProfiler() (p *Profiler, err error) {
	prof, err := processors.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: prof}, nil
}

// Get returns the current processor info as JSON serialized bytes.
func (p *Profiler) Get() (b []byte, err error) {
	procs, err := p.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return p.Serialize(procs)
}

var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current processors info as JSON serialized bytes using
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

// Serialize processors info as JSON
func (p *Profiler) Serialize(proc *processors.Processors) ([]byte, error) {
	return json.Marshal(proc)
}

// Serialize processors info as JSON using package globals.
func Serialize(proc *processors.Processors) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(proc)
}

// Marshal is an alias for serialize.
func (p *Profiler) Marshal(proc *processors.Processors) ([]byte, error) {
	return p.Serialize(proc)
}

// Marshal is an alias for Serialize using package globals.
func Marshal(proc *processors.Processors) ([]byte, error) {
	return std.Serialize(proc)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// Processors
func Deserialize(p []byte) (*processors.Processors, error) {
	proc := &processors.Processors{}
	err := json.Unmarshal(p, proc)
	if err != nil {
		return nil, err
	}
	return proc, nil
}

// Unmarshal is an alias for Deserialize
func Unmarshal(p []byte) (*processors.Processors, error) {
	return Deserialize(p)
}
