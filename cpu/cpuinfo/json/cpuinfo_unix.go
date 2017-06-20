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

// Package cpuinfo (json) handles JSON based processing of CPU info. Instead
// of returning a Go struct, it returns JSON serialized bytes. A function to
// deserialize the JSON serialized bytes into a cpuinfo.CPUs struct is
// provided.
//
// Note: the package name is cpuinfo and not the final element of the import
// path (json). 
package cpuinfo

import (
	"encoding/json"
	"sync"

	inf "github.com/mohae/joefriday/cpu/cpuinfo"
)

// Profiler is used to process the cpuinfo (cpus) as JSON serialized bytes.
type Profiler struct {
	*inf.Profiler
}

// Initializes and returns a cpuinfo profiler.
func NewProfiler() (prof *Profiler, err error) {
	p, err := inf.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p}, nil
}

// Get returns the current cpuinfo (CPUs) as JSON serialized bytes.
func (prof *Profiler) Get() (p []byte, err error) {
	cpus, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(cpus)
}

var std *Profiler
var stdMu sync.Mutex //protects standard to prevent data race on checking/instantiation

// Get returns the current cpuinfo (CPUs) as JSON serialized bytes using the
// package's global Profiler.
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

// Serialize cpuinfo (CPUs) as JSON.
func (prof *Profiler) Serialize(cpus *inf.CPUs) ([]byte, error) {
	return json.Marshal(cpus)
}

// Serialize cpuinfo (CPUs) as JSON using package globals.
func Serialize(cpus *inf.CPUs) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(cpus)
}

// Marshal is an alias for serialize.
func (prof *Profiler) Marshal(cpus *inf.CPUs) ([]byte, error) {
	return prof.Serialize(cpus)
}

// Marshal is an alias for Serialize using package globals.
func Marshal(cpus *inf.CPUs) ([]byte, error) {
	return std.Serialize(cpus)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// inf.CPUs
func Deserialize(p []byte) (*inf.CPUs, error) {
	cpus := &inf.CPUs{}
	err := json.Unmarshal(p, cpus)
	if err != nil {
		return nil, err
	}
	return cpus, nil
}

// Unmarshal is an alias for Deserialize using package globals.
func Unmarshal(p []byte) (*inf.CPUs, error) {
	return Deserialize(p)
}
