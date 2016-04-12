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

// Package json handles JSON based processing of loadavg information.  Instead
// of returning a Go struct, it returns JSON serialized bytes.  A function to
// deserialize the JSON serialized bytes into an loadavg.LoadAvg struct is
// provided.
package json

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/mohae/joefriday/platform/loadavg"
)

// Profiler is used to process the loadavg information, /proc/loadavg, using
// JSON.
type Profiler struct {
	*loadavg.Profiler
}

// Initializes and returns a json.Profiler for loadavg information.
func New() (prof *Profiler, err error) {
	p, err := loadavg.New()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p}, nil
}

// Get returns the current loadavg information as JSON serialized bytes.
func (prof *Profiler) Get() (p []byte, err error) {
	k, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(k)
}

// Ticker returns the current loadavg as Flatbuffer serialized bytes on a
// ticker.
func (prof *Profiler) Ticker(d time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	outCh := make(chan loadavg.LoadAvg)
	defer close(out)
	var (
		ok  bool
		p   []byte
		err error
		l   loadavg.LoadAvg
	)
	go prof.Profiler.Ticker(d, outCh, done, errs)
	for {
		select {
		case l, ok = <-outCh:
			if !ok {
				return
			}
			p, err = prof.Serialize(l)
			if err != nil {
				errs <- err
				continue
			}
			out <- p
		}
	}
}

var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current loadavg information as JSON serialized bytes using
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

// Ticker gathers information on a ticker using the specified interval.
// This uses a local Profiler as using the global doesn't make sense for
// an ongoing ticker.
func Ticker(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	prof, err := New()
	if err != nil {
		errs <- err
		close(out)
		return
	}
	prof.Ticker(interval, out, done, errs)
}

// Serialize loadavg.LoadAvg using JSON
func (prof *Profiler) Serialize(u loadavg.LoadAvg) ([]byte, error) {
	return json.Marshal(u)
}

// Serialize loadavg.LoadAvg using JSON with the package global Profiler.
func Serialize(u loadavg.LoadAvg) (p []byte, err error) {
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
func (prof *Profiler) Marshal(l loadavg.LoadAvg) ([]byte, error) {
	return prof.Serialize(l)
}

// Marshal is an alias for Serialize that uses the package's global profiler.
func Marshal(u loadavg.LoadAvg) ([]byte, error) {
	return Serialize(u)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// loadavg.LoadAvg.
func Deserialize(p []byte) (l loadavg.LoadAvg, err error) {
	err = json.Unmarshal(p, &l)
	if err != nil {
		return l, err
	}
	return l, nil
}

// Unmarshal is an alias for Deserialize
func Unmarshal(p []byte) (loadavg.LoadAvg, error) {
	return Deserialize(p)
}
