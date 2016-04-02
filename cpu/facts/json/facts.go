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

// Package mem gets and processes /proc/meminfo, returning the data in the
// appropriate format.
package json

import (
	"encoding/json"
	"sync"

	"github.com/mohae/joefriday/cpu/facts"
)

type Profiler struct {
	Prof *facts.Profiler
}

func New() (prof *Profiler, err error) {
	p, err := facts.New()
	if err != nil {
		return nil, err
	}
	return &Profiler{Prof: p}, nil
}

// Get returns some of the results of /proc/meminfo.
func (prof *Profiler) Get() (p []byte, err error) {
	prof.Prof.Reset()
	fct, err := prof.Prof.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(fct)
}

// TODO: is it even worth it to have this as a global?  Should GetInfo()
// just instantiate a local version and use that?  InfoTicker does...
var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get get's the current meminfo.
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

// Serialize mem.Info as JSON
func (prof *Profiler) Serialize(fct *facts.Facts) ([]byte, error) {
	return json.Marshal(fct)
}

// Unmarshal unmarshals JSON into *Info.
func Unmarshal(p []byte) (*facts.Facts, error) {
	fct := &facts.Facts{}
	err := json.Unmarshal(p, fct)
	if err != nil {
		return nil, err
	}
	return fct, nil
}
