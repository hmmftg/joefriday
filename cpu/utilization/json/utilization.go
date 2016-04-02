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

package json

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/mohae/joefriday/cpu/utilization"
)

type Profiler struct {
	Prof *utilization.Profiler
}

func New() (prof *Profiler, err error) {
	p, err := utilization.New()
	if err != nil {
		return nil, err
	}
	return &Profiler{Prof: p}, nil
}

// Get returns some of the results of
func (prof *Profiler) Get() (p []byte, err error) {
	prof.Prof.Reset()
	st, err := prof.Prof.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(st)
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

// Ticker processes CPU utilization information on a ticker as JSON
// serialized bytes.  Any errors encountered are put on the errs channel.
// Processing ends when a done signal is received.
//
// It is the callers responsibility to close the done and errs channels.
//
// TODO: better handle errors, e.g. restore cur from prior so that there
// isn't the possibility of temporarily having bad data, just a missed
// collection interval.
func (prof *Profiler) Ticker(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	outCh := make(chan *utilization.Utilization)
	defer close(outCh)
	go prof.Prof.Ticker(interval, outCh, done, errs)
	for {
		select {
		case u, ok := <-outCh:
			if !ok {
				return
			}
			b, err := prof.Serialize(u)
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
	prof, err := New()
	if err != nil {
		errs <- err
		close(out)
		return
	}
	prof.Ticker(interval, out, done, errs)
}

// Serialize mem.Info as JSON
func (prof *Profiler) Serialize(ut *utilization.Utilization) ([]byte, error) {
	return json.Marshal(ut)
}

// Unmarshal unmarshals JSON into *Info.
func Unmarshal(p []byte) (*utilization.Utilization, error) {
	ut := &utilization.Utilization{}
	err := json.Unmarshal(p, ut)
	if err != nil {
		return nil, err
	}
	return ut, nil
}
