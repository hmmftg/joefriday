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

// Package flat handles Flatbuffer based processing of network interface
//usage.  Interface usage is calculated by taking the difference in two
// /proc/net/dev snapshots and reflect bytes received and transmitted since
// the prior snapshot.  Instead of returning a Go struct, it returns
// Flatbuffer serialized bytes.  A function to deserialize the Flatbuffer
// serialized bytes into a structs.Usage struct.  After the first use, the
// flatbuffer builder is reused.
package flat

import (
	"sync"
	"time"

	"github.com/mohae/joefriday/net/info/flat"
	"github.com/mohae/joefriday/net/structs"
)

// Profiler is used to process the network interface usage using Flatbuffers.
type Profiler struct {
	Flat *flat.Profiler
}

// Initializes and returns a network interface usage profiler that utilizes
// FlatBuffers.
func New() (prof *Profiler, err error) {
	p, err := flat.New()
	if err != nil {
		return nil, err
	}
	return &Profiler{Flat: p}, nil
}

// Get returns the current network interface usage as Flatbuffer serialized
// bytes.
// TODO: should this be changed so that this calculates usage since the last
// time the network info was obtained.  If there aren't pre-existing info
// it would get current usage (which may be a separate method (or should be?))
func (prof *Profiler) Get() (p []byte, err error) {
	u, err := prof.Flat.Prof.Get()
	if err != nil {
		return nil, err
	}
	return prof.Flat.Serialize(u), nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns the current network interface usage as Flatbuffer serialized
// bytes using the package's global Profiler.
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

// Ticker processes network interface usage on a ticker.  The generated data
// is sent to the out channel.  Any errors encountered are sent to the errs
// channel.  Processing ends when a done signal is received.
//
// It is the callers responsibility to close the done and errs channels.
//
// TODO: better handle errors, e.g. restore cur from prior so that there
// isn't the possibility of temporarily having bad data, just a missed
// collection interval.
func (prof *Profiler) Ticker(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	outCh := make(chan *structs.Info)
	defer close(outCh)
	go prof.Flat.Prof.Ticker(interval, outCh, done, errs)
	for {
		select {
		case inf, ok := <-outCh:
			if !ok {
				return
			}
			out <- prof.Flat.Serialize(inf)
		}
	}
}

// Ticker gathers network interface usage on a ticker using the specified
// interval.  This uses a local Profiler as using the global doesn't make
// sense for an ongoing ticker.
func Ticker(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	prof, err := New()
	if err != nil {
		errs <- err
		close(out)
		return
	}
	prof.Ticker(interval, out, done, errs)
}
