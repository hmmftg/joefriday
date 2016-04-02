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

package flat

import (
	"sync"
	"time"

	"github.com/mohae/joefriday/net/info/flat"
	"github.com/mohae/joefriday/net/structs"
)

type Profiler struct {
	Flat *flat.Profiler
}

func New() (prof *Profiler, err error) {
	p, err := flat.New()
	if err != nil {
		return nil, err
	}
	return &Profiler{Flat: p}, nil
}

// Get returns the network usage.  Usage calculations requires two pieces of
// data.  This func gets a snapshot of /proc/net/dev, sleeps for a/ second,
// and takes another snapshot and calcualtes the usage from the two snapshots.
// If ongoing usage information is desired, Ticker should be called; it's
// better suited for ongoing usage information: using less cpu cycles and
// generating less garbage.
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

// Ticker processes network usage information on a ticker.  Errors are sent
// on the errs channel.  Processing ends when a done signal is received.
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
