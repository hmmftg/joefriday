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

// Package json handles JSON based processing of network interface
// information. Instead of returning a Go struct, it returns JSON serialized
// bytes.  A function to deserialize the JSON serialized bytes into a
// structs.Info struct is provided.
package json

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/mohae/joefriday/net/info"
	"github.com/mohae/joefriday/net/structs"
)

// Profiler is used to process the /proc/net/dev file using JSON.
type Profiler struct {
	Prof *info.Profiler
}

// Initializes and returns a network information profiler.
func New() (prof *Profiler, err error) {
	p, err := info.New()
	if err != nil {
		return nil, err
	}
	return &Profiler{Prof: p}, nil
}

// Get returns the current network information as JSON serialized bytes.
func (prof *Profiler) Get() (p []byte, err error) {
	inf, err := prof.Prof.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(inf)
}

var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current network interface information as JSON serialized
// bytes using the package's globla Profiler.
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

// Ticker processes network interface information on a ticker.  The generated
// data is sent to the out channel.  Any errors encountered are sent to the
// errs channel.  Processing ends when a done signal is received.
//
// It is the callers responsibility to close the done and errs channels.
func (prof *Profiler) Ticker(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	outCh := make(chan *structs.Info)
	defer close(outCh)
	go prof.Prof.Ticker(interval, outCh, done, errs)
	for {
		select {
		case inf, ok := <-outCh:
			if !ok {
				return
			}
			b, err := prof.Serialize(inf)
			if err != nil {
				errs <- err
				continue
			}
			out <- b
		}
	}
}

// Ticker gathers network interface information on a ticker using the
// specified interval.  This uses a local Profiler as using the global
// doesn't make sense for an ongoing ticker.
func Ticker(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	p, err := New()
	if err != nil {
		errs <- err
		return
	}
	p.Ticker(interval, out, done, errs)
}

// Serialize network interface information using JSON
func (prof *Profiler) Serialize(inf *structs.Info) ([]byte, error) {
	return json.Marshal(inf)
}

// Serialize network interface information using JSON with the package global
// Profiler.
func Serialize(inf *structs.Info) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = New()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(inf)
}

// Marshal is an alias for Serialize
func (prof *Profiler) Marshal(inf *structs.Info) ([]byte, error) {
	return prof.Serialize(inf)
}

// Marshal is an alias for Serialize; uses the package global Profiler.
func Marshal(inf *structs.Info) ([]byte, error) {
	return Serialize(inf)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// structs.Info
func Deserialize(p []byte) (*structs.Info, error) {
	info := &structs.Info{}
	err := json.Unmarshal(p, info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// Unmarshal is an alias for Deserialize
func Unmarshal(p []byte) (*structs.Info, error) {
	return Deserialize(p)
}
