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

// Package json handles JSON based processing of network usage.  Instead of
// returning a Go struct, it returns JSON serialized bytes.  A function to
// deserialize the JSON serialized bytes into a structs.Info struct is
// provided.
package json

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/mohae/joefriday/net/structs"
	"github.com/mohae/joefriday/net/usage"
)

// Profiler is used to process the network usage JSON.
type Profiler struct {
	Prof *usage.Profiler
}

// Initializes and returns a network usage profiler.
func New() (prof *Profiler, err error) {
	p, err := usage.New()
	if err != nil {
		return nil, err
	}
	return &Profiler{Prof: p}, nil
}

// Get returns the current network usage as JSON serialized bytes.
func (prof *Profiler) Get() (p []byte, err error) {
	inf, err := prof.Prof.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(inf)
}

// TODO: is it even worth it to have this as a global?  Should GetInfo()
// just instantiate a local version and use that?  InfoTicker does...
var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current network usage as JSON serialized bytes using the
// package's global Profiler.
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

// Ticker processes network usage on a ticker.  The generated data is sent to
// the out channel.  Any errors encountered are sent to the errs channel.
// Processing ends when a done signal is received.
//
// It is the callers responsibility to close the done and errs channels.
//
// TODO: better handle errors, e.g. restore cur from prior so that there
// isn't the possibility of temporarily having bad data, just a missed
// collection interval.
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

// Ticker gathers information on a ticker using the specified interval.
// This uses a local Profiler as using the global doesn't make sense for
// an ongoing ticker.
func Ticker(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	p, err := New()
	if err != nil {
		errs <- err
		return
	}
	p.Ticker(interval, out, done, errs)
}

// Serialize network usage using JSON
func (prof *Profiler) Serialize(inf *structs.Info) ([]byte, error) {
	return json.Marshal(inf)
}

// Marshal is an alias for Serialize
func (prof *Profiler) Marshal(inf *structs.Info) ([]byte, error) {
	return prof.Serialize(inf)
}

// Deserialize deserializes JSON serialized bytes.
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
