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

// Package json handles JSON based processing of load using syscall.
// Instead of returning a Go struct, it returns JSON serialized bytes.  A
// function to deserialize the JSON serialized bytes into a load.LoadAvg
// struct is provided.
package json

import (
	"encoding/json"
	"time"

	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/sysinfo/load"
)

// Get returns the current load as JSON serialized bytes.
func Get() (p []byte, err error) {
	var l load.LoadAvg
	err = l.Get()
	if err != nil {
		return nil, err
	}
	return json.Marshal(&l)
}

// Deserialize takes some JSON serialized bytes and unmarshals them as
// load.LoadAvg.
func Deserialize(p []byte) (*load.LoadAvg, error) {
	var l load.LoadAvg
	err := json.Unmarshal(p, &l)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

// Unmarshal is an alias for Deserialize
func Unmarshal(p []byte) (*load.LoadAvg, error) {
	return Deserialize(p)
}

// Ticker delivers load.LoadAvg as JSON serialized bytes at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan []byte
}

// NewTicker returns a new Ticker continaing a Data channel that delivers
// the data at intervals and an error channel that delivers any errors
// encountered.  Stop the ticker to signal the ticker to stop running; it
// does not close the Data channel.  Close the ticker to close all ticker
// channels.
func NewTicker(d time.Duration) (joe.Tocker, error) {
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan []byte)}
	go t.Run()
	return &t, nil
}

// Run runs the ticker.
func (t *Ticker) Run() {
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

// Close closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Close()
	close(t.Data)
}
