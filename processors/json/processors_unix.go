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

// Package processors handles JSON based processing of Processor info. Instead
// of returning a Go struct, it returns JSON serialized bytes. A function to
// deserialize the JSON serialized bytes into a processors.Processors struct
// is provided.
//
// Note: the package name is processors and not the final element of the import
// path (json). 
package processors

import (
	"encoding/json"
	"sync"

	procs "github.com/mohae/joefriday/processors"
)

// Profiler is used to process the processor info as JSON serialized bytes.
type Profiler struct {
	*procs.Profiler
}

// Initializes and returns a processor profiler.
func NewProfiler() (p *Profiler, err error) {
	prof, err := procs.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: prof}, nil
}

// Get returns the current processor info as JSON serialized bytes.
func (p *Profiler) Get() (b []byte, err error) {
	proc, err := p.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return p.Serialize(proc)
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

// Serialize processor info as JSON.
func (p *Profiler) Serialize(proc *procs.Processors) ([]byte, error) {
	return json.Marshal(proc)
}

// Serialize processor info as JSON using package globals.
func Serialize(proc *procs.Processors) (p []byte, err error) {
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
func (p *Profiler) Marshal(proc *procs.Processors) ([]byte, error) {
	return p.Serialize(proc)
}

// Marshal is an alias for Serialize using package globals.
func Marshal(proc *procs.Processors) ([]byte, error) {
	return std.Serialize(proc)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// Processors
func Deserialize(p []byte) (*procs.Processors, error) {
	proc := &procs.Processors{}
	err := json.Unmarshal(p, proc)
	if err != nil {
		return nil, err
	}
	return proc, nil
}

// Unmarshal is an alias for Deserialize
func Unmarshal(p []byte) (*procs.Processors, error) {
	return Deserialize(p)
}
