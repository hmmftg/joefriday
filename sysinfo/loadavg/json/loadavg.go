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

// Package json handles JSON based processing of /proc/loadavgl.  Instead of
// returning a Go struct, it returns JSON serialized bytes.  A function to
// deserialize the JSON serialized bytes into a loadavg.LoadAvg struct is
// provided.
package json

import (
	"encoding/json"
	"time"

	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/sysinfo/loadavg"
)

// Get returns the current loadavgl as JSON serialized bytes using the
// package's global Profiler.
func Get() (p []byte, err error) {
	var l loadavg.LoadAvg
	err = l.Get()
	if err != nil {
		return nil, err
	}
	return json.Marshal(&l)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// loadavg.LoadAvg.
func Deserialize(p []byte) (*loadavg.LoadAvg, error) {
	var l loadavg.LoadAvg
	err := json.Unmarshal(p, &l)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

// Unmarshal is an alias for Deserialize
func Unmarshal(p []byte) (*loadavg.LoadAvg, error) {
	return Deserialize(p)
}

type Ticker struct {
	*joe.Ticker
	Data chan []byte
}

// NewTicker returns a new Ticker containing a ticker channel, T,
func NewTicker(d time.Duration) (joe.Tocker, error) {
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan []byte)}
	go t.Run()
	return &t, nil
}

// Stop sends a signal to the done channel, stopping the goroutine.  This
// also closes all channels that the Ticker holds.
func (t *Ticker) Stop() {
	t.Done <- struct{}{}
}

func (t *Ticker) Run() {
	defer t.Close()
	defer close(t.Data)
	// read until done signal is received
	for {
		select {
		case <-t.Done:
			return
		case <-t.Ticker.C:
			p, err := Get()
			if err != nil {
				t.Errs <- err
				continue
			}
			t.Data <- p
		}
	}
}
