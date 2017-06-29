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

// Package version handles JSON based processing of kernel and version
// information: /proc/version. Instead of returning a Go struct, it returns
// JSON serialized bytes. A function to deserialize the JSON serialized bytes
// into a version.Info struct is provided.
//
// Note: the package name is version and not the final element of the import
// path (json). 
package version

import (
	"encoding/json"
	"sync"

	v "github.com/mohae/joefriday/system/version"
)

// Profiler is used to process the version information, /proc/version, using
// JSON.
type Profiler struct {
	*v.Profiler
}

// Initializes and returns a json.Profiler for version information.
func NewProfiler() (prof *Profiler, err error) {
	p, err := v.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p}, nil
}

// Get returns the current version information as JSON serialized bytes.
func (prof *Profiler) Get() (p []byte, err error) {
	inf, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(inf)
}

var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current version information as JSON serialized bytes using
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

// Serialize version.Info using JSON
func (prof *Profiler) Serialize(inf *v.Info) ([]byte, error) {
	return json.Marshal(inf)
}

// Serialize version.Info using JSON with the package global Profiler.
func Serialize(inf *v.Info) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(inf)
}

// Marshal is an alias for Serialize
func (prof *Profiler) Marshal(inf *v.Info) ([]byte, error) {
	return prof.Serialize(inf)
}

// Marshal is an alias for Serialize that uses the package's global profiler.
func Marshal(inf *v.Info) ([]byte, error) {
	return Serialize(inf)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// version.Info.
func Deserialize(p []byte) (*v.Info, error) {
	inf := &v.Info{}
	err := json.Unmarshal(p, inf)
	if err != nil {
		return nil, err
	}
	return inf, nil
}

// Unmarshal is an alias for Deserialize
func Unmarshal(p []byte) (*v.Info, error) {
	return Deserialize(p)
}
